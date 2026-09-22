#!/usr/bin/env python3
"""Fail closed unless a Scope Reviewer optic uses a fresh Luna Max route."""

from __future__ import annotations

import argparse
import json
from dataclasses import asdict, dataclass


SCHEMA = "scope-reviewer/lens-routing/v1"
EXPECTED_AGENT_TYPE = "default"
EXPECTED_MODEL = "gpt-5.6-luna"
EXPECTED_EFFORT = "max"
EXPECTED_FORK_TURNS = "none"


@dataclass(frozen=True)
class LensRoutingReceipt:
    schema: str
    optic: str
    snapshot_id: str
    agent_type: str
    requested_model: str
    requested_effort: str
    fork_turns: str
    telemetry_status: str
    actual_model: str | None
    actual_effort: str | None
    allowed: bool
    defects: tuple[str, ...]


def validate_route(
    *,
    optic: str,
    snapshot_id: str,
    agent_type: str,
    model: str,
    effort: str,
    fork_turns: str,
    actual_model: str | None = None,
    actual_effort: str | None = None,
) -> LensRoutingReceipt:
    normalized_optic = optic.strip()
    normalized_snapshot = snapshot_id.strip()
    normalized_agent = agent_type.strip().casefold()
    normalized_model = model.strip().casefold()
    normalized_effort = effort.strip().casefold()
    normalized_fork = fork_turns.strip().casefold()
    normalized_actual_model = (
        actual_model.strip().casefold() if actual_model is not None else None
    )
    normalized_actual_effort = (
        actual_effort.strip().casefold() if actual_effort is not None else None
    )

    defects: list[str] = []
    if not normalized_optic:
        defects.append("blank_optic")
    if not normalized_snapshot:
        defects.append("blank_snapshot_id")
    if normalized_agent != EXPECTED_AGENT_TYPE:
        defects.append("default_agent_required")
    if normalized_model != EXPECTED_MODEL:
        defects.append("luna_model_required")
    if normalized_effort != EXPECTED_EFFORT:
        defects.append("luna_max_effort_required")
    if normalized_fork != EXPECTED_FORK_TURNS:
        defects.append("fresh_fork_required")
    if (actual_model is None) != (actual_effort is None):
        defects.append("incomplete_actual_profile")
    if normalized_actual_model is not None and normalized_actual_model != EXPECTED_MODEL:
        defects.append("actual_luna_model_mismatch")
    if normalized_actual_effort is not None and normalized_actual_effort != EXPECTED_EFFORT:
        defects.append("actual_luna_effort_mismatch")

    return LensRoutingReceipt(
        schema=SCHEMA,
        optic=normalized_optic,
        snapshot_id=normalized_snapshot,
        agent_type=normalized_agent,
        requested_model=normalized_model,
        requested_effort=normalized_effort,
        fork_turns=normalized_fork,
        telemetry_status=(
            "observed"
            if normalized_actual_model is not None
            and normalized_actual_effort is not None
            else "telemetry_pending"
        ),
        actual_model=normalized_actual_model,
        actual_effort=normalized_actual_effort,
        allowed=not defects,
        defects=tuple(defects),
    )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--optic", required=True)
    parser.add_argument("--snapshot-id", required=True)
    parser.add_argument("--agent-type", default=EXPECTED_AGENT_TYPE)
    parser.add_argument("--model", required=True)
    parser.add_argument("--effort", required=True)
    parser.add_argument("--fork-turns", required=True)
    parser.add_argument("--actual-model")
    parser.add_argument("--actual-effort")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    receipt = validate_route(
        optic=args.optic,
        snapshot_id=args.snapshot_id,
        agent_type=args.agent_type,
        model=args.model,
        effort=args.effort,
        fork_turns=args.fork_turns,
        actual_model=args.actual_model,
        actual_effort=args.actual_effort,
    )
    print(json.dumps(asdict(receipt), ensure_ascii=False, sort_keys=True))
    return 0 if receipt.allowed else 2


if __name__ == "__main__":
    raise SystemExit(main())
