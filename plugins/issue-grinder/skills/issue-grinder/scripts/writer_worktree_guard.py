#!/usr/bin/env python3
"""Create and verify isolated Issue Grinder writer worktrees.

The helper deliberately owns only mechanical Git invariants.  Scope, packet
selection, agent ownership, integration, and cleanup decisions stay with the
Issue Grinder coordinator.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
from typing import Any


WRITER_SCHEMA = "issue-grinder/writer-worktree/v1"
GUARD_SCHEMA = "issue-grinder/integration-guard/v1"
ERROR_SCHEMA = "issue-grinder/writer-worktree-error/v1"
OWNER = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._:/-]{0,127}$")


class GuardError(RuntimeError):
    """A fail-closed writer isolation error."""


def run_git(
    worktree: Path,
    *arguments: str,
    check: bool = True,
) -> subprocess.CompletedProcess[bytes]:
    result = subprocess.run(
        ["git", "-C", str(worktree), *arguments],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if check and result.returncode != 0:
        message = result.stderr.decode("utf-8", errors="replace").strip()
        raise GuardError(
            f"git {' '.join(arguments)} failed"
            + (f": {message}" if message else "")
        )
    return result


def git_text(worktree: Path, *arguments: str) -> str:
    return run_git(worktree, *arguments).stdout.decode("utf-8").strip()


def resolved_git_path(worktree: Path, argument: str) -> Path:
    value = Path(git_text(worktree, "rev-parse", argument))
    if not value.is_absolute():
        value = worktree / value
    return value.resolve()


def branch_ref(worktree: Path) -> str:
    result = run_git(worktree, "symbolic-ref", "--quiet", "HEAD", check=False)
    if result.returncode != 0:
        raise GuardError(f"{worktree} has detached HEAD")
    return result.stdout.decode("utf-8").strip()


def parse_worktrees(repository: Path) -> list[dict[str, str | bool]]:
    raw = run_git(repository, "worktree", "list", "--porcelain", "-z").stdout
    records: list[dict[str, str | bool]] = []
    current: dict[str, str | bool] = {}
    for token in raw.split(b"\0"):
        if not token:
            if current:
                records.append(current)
                current = {}
            continue
        text = token.decode("utf-8", errors="strict")
        key, separator, value = text.partition(" ")
        current[key] = value if separator else True
    if current:
        records.append(current)
    return records


def status_bytes(worktree: Path) -> bytes:
    return run_git(
        worktree,
        "status",
        "--porcelain=v1",
        "-z",
        "--untracked-files=all",
    ).stdout


def inspect_worktree(worktree: Path) -> dict[str, Any]:
    requested = worktree.resolve()
    top = Path(git_text(requested, "rev-parse", "--show-toplevel")).resolve()
    if top != requested:
        raise GuardError(
            f"expected worktree root {requested}, but Git resolved {top}"
        )
    git_dir = resolved_git_path(top, "--git-dir")
    common_dir = resolved_git_path(top, "--git-common-dir")
    head = git_text(top, "rev-parse", "HEAD")
    ref = branch_ref(top)
    status = status_bytes(top)
    records = parse_worktrees(top)
    record = next(
        (
            item
            for item in records
            if isinstance(item.get("worktree"), str)
            and Path(str(item["worktree"])).resolve() == top
        ),
        None,
    )
    if record is None:
        raise GuardError(f"{top} is absent from git worktree list --porcelain")
    return {
        "worktree": str(top),
        "git_dir": str(git_dir),
        "common_dir": str(common_dir),
        "head": head,
        "branch": ref,
        "clean": not status,
        "status_sha256": hashlib.sha256(status).hexdigest(),
        "linked": git_dir != common_dir,
        "locked": "locked" in record,
    }


def write_json(payload: dict[str, Any], output: Path | None) -> None:
    rendered = json.dumps(payload, ensure_ascii=False, indent=2, sort_keys=True) + "\n"
    if output is None:
        sys.stdout.write(rendered)
        return
    target = output.resolve()
    if not target.parent.is_dir():
        raise GuardError(f"receipt parent does not exist: {target.parent}")
    with tempfile.NamedTemporaryFile(
        "w",
        encoding="utf-8",
        dir=target.parent,
        prefix=f".{target.name}.",
        delete=False,
    ) as handle:
        handle.write(rendered)
        temporary = Path(handle.name)
    os.replace(temporary, target)
    sys.stdout.write(rendered)


def assert_owner(value: str) -> None:
    if not OWNER.fullmatch(value):
        raise GuardError(
            "owner must be 1-128 bounded letters, digits, dot, underscore, colon, slash, or dash"
        )


def inside(path: Path, parent: Path) -> bool:
    try:
        path.relative_to(parent)
        return True
    except ValueError:
        return False


def ensure_external_output(output: Path | None, *worktrees: Path) -> None:
    if output is None:
        return
    target = output.resolve()
    for worktree in worktrees:
        root = worktree.resolve()
        if target == root or inside(target, root):
            raise GuardError(
                f"receipt must be outside every Git worktree: {target} is inside {root}"
            )


def prepare(arguments: argparse.Namespace) -> dict[str, Any]:
    repository = Path(arguments.repo).resolve()
    target = Path(arguments.worktree).resolve()
    assert_owner(arguments.owner)
    repository_state = inspect_worktree(repository)
    base = git_text(repository, "rev-parse", f"{arguments.base}^{{commit}}")
    existing_worktrees = [
        Path(str(record["worktree"]))
        for record in parse_worktrees(repository)
        if isinstance(record.get("worktree"), str)
    ]
    ensure_external_output(arguments.output, *existing_worktrees, target)

    if target.exists():
        raise GuardError(f"writer worktree path already exists: {target}")
    for existing_value in existing_worktrees:
        existing = existing_value.resolve()
        if inside(target, existing) or inside(existing, target):
            raise GuardError(
                f"writer worktree path overlaps existing worktree: {existing}"
            )

    ref_check = run_git(
        repository,
        "check-ref-format",
        "--branch",
        arguments.branch,
        check=False,
    )
    if ref_check.returncode != 0:
        raise GuardError(f"invalid writer branch: {arguments.branch}")
    existing_branch = run_git(
        repository,
        "show-ref",
        "--verify",
        "--quiet",
        f"refs/heads/{arguments.branch}",
        check=False,
    )
    if existing_branch.returncode == 0:
        raise GuardError(f"writer branch already exists: {arguments.branch}")

    run_git(
        repository,
        "worktree",
        "add",
        "--lock",
        "--reason",
        f"issue-grinder:{arguments.owner}",
        "-b",
        arguments.branch,
        str(target),
        base,
    )
    state = inspect_worktree(target)
    expected_branch = f"refs/heads/{arguments.branch}"
    failures = []
    if state["common_dir"] != repository_state["common_dir"]:
        failures.append("different common Git directory")
    if not state["linked"]:
        failures.append("not a linked worktree")
    if not state["locked"]:
        failures.append("worktree is not locked")
    if state["branch"] != expected_branch:
        failures.append("unexpected branch")
    if state["head"] != base:
        failures.append("unexpected base commit")
    if not state["clean"]:
        failures.append("new worktree is dirty")
    if failures:
        raise GuardError(
            "prepared worktree failed verification: " + ", ".join(failures)
        )
    return {
        "schema": WRITER_SCHEMA,
        "operation": "prepare",
        "owner": arguments.owner,
        "base": base,
        **state,
    }


def admit(arguments: argparse.Namespace) -> dict[str, Any]:
    expected = Path(arguments.expect_worktree).resolve()
    current = Path.cwd().resolve()
    assert_owner(arguments.owner)
    ensure_external_output(arguments.output, expected)
    if current != expected:
        raise GuardError(
            f"current working directory {current} is not admitted writer worktree {expected}"
        )
    state = inspect_worktree(current)
    failures = []
    if not state["linked"]:
        failures.append("current checkout is the main worktree")
    if not state["locked"]:
        failures.append("writer worktree is not locked")
    if state["branch"] != f"refs/heads/{arguments.expect_branch}":
        failures.append("branch mismatch")
    if arguments.expect_head and state["head"] != arguments.expect_head:
        failures.append("HEAD mismatch")
    if arguments.expect_common_dir:
        common = str(Path(arguments.expect_common_dir).resolve())
        if state["common_dir"] != common:
            failures.append("common Git directory mismatch")
    if not arguments.allow_dirty and not state["clean"]:
        failures.append("writer worktree is dirty before admission")
    if failures:
        raise GuardError("writer admission rejected: " + ", ".join(failures))
    return {
        "schema": WRITER_SCHEMA,
        "operation": "admit",
        "owner": arguments.owner,
        **state,
    }


def snapshot(arguments: argparse.Namespace) -> dict[str, Any]:
    worktree = Path(arguments.worktree).resolve()
    ensure_external_output(arguments.output, worktree)
    state = inspect_worktree(worktree)
    if not state["clean"]:
        raise GuardError(
            "integration worktree must be clean before a parallel writer wave"
        )
    return {
        "schema": GUARD_SCHEMA,
        "operation": "snapshot",
        **state,
    }


def assert_unchanged(arguments: argparse.Namespace) -> dict[str, Any]:
    receipt_path = Path(arguments.snapshot).resolve()
    try:
        expected = json.loads(receipt_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        raise GuardError(f"cannot read integration snapshot: {error}") from error
    if expected.get("schema") != GUARD_SCHEMA:
        raise GuardError("integration snapshot has an unsupported schema")
    worktree = Path(str(expected.get("worktree", ""))).resolve()
    ensure_external_output(arguments.output, worktree)
    current = inspect_worktree(worktree)
    compared = ("worktree", "git_dir", "common_dir", "head", "branch", "status_sha256")
    changed = [key for key in compared if current.get(key) != expected.get(key)]
    if changed:
        raise GuardError(
            "integration worktree changed during writer wave: " + ", ".join(changed)
        )
    return {
        "schema": GUARD_SCHEMA,
        "operation": "assert-unchanged",
        "unchanged": True,
        **current,
    }


def output_argument(parser: argparse.ArgumentParser) -> None:
    parser.add_argument(
        "--output",
        type=Path,
        help="Optionally write the JSON receipt atomically to this existing parent.",
    )


def parser() -> argparse.ArgumentParser:
    root = argparse.ArgumentParser(description=__doc__)
    commands = root.add_subparsers(dest="command", required=True)

    prepare_parser = commands.add_parser("prepare")
    prepare_parser.add_argument("--repo", required=True)
    prepare_parser.add_argument("--worktree", required=True)
    prepare_parser.add_argument("--branch", required=True)
    prepare_parser.add_argument("--base", required=True)
    prepare_parser.add_argument("--owner", required=True)
    output_argument(prepare_parser)
    prepare_parser.set_defaults(handler=prepare)

    admit_parser = commands.add_parser("admit")
    admit_parser.add_argument("--expect-worktree", required=True)
    admit_parser.add_argument("--expect-branch", required=True)
    admit_parser.add_argument("--expect-head", required=True)
    admit_parser.add_argument("--expect-common-dir", required=True)
    admit_parser.add_argument("--owner", required=True)
    admit_parser.add_argument(
        "--allow-dirty",
        action="store_true",
        help="Only for a proven quiescent unfinished task-owned checkpoint.",
    )
    output_argument(admit_parser)
    admit_parser.set_defaults(handler=admit)

    snapshot_parser = commands.add_parser("snapshot")
    snapshot_parser.add_argument("--worktree", required=True)
    output_argument(snapshot_parser)
    snapshot_parser.set_defaults(handler=snapshot)

    unchanged_parser = commands.add_parser("assert-unchanged")
    unchanged_parser.add_argument("--snapshot", required=True)
    output_argument(unchanged_parser)
    unchanged_parser.set_defaults(handler=assert_unchanged)
    return root


def main() -> int:
    arguments = parser().parse_args()
    try:
        payload = arguments.handler(arguments)
        write_json(payload, arguments.output)
        return 0
    except GuardError as error:
        sys.stderr.write(
            json.dumps(
                {"schema": ERROR_SCHEMA, "error": str(error)},
                ensure_ascii=False,
                sort_keys=True,
            )
            + "\n"
        )
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
