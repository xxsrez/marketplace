import hashlib
import json
import os
import platform
import re
import shutil
import subprocess
import tarfile
import tempfile
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
MIND_DIARY_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "mind-diary"
UAT_OAUTH_RESOURCE = "https://mind-diary.xxsrez-work.chatgpt.site/api/mcp"
UAT_CODEX_MCP_URL = f"{UAT_OAUTH_RESOURCE}/2025-11-25"


def read_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


class MindDiaryDirectMcpPackagingTest(unittest.TestCase):
    def test_marketplace_installs_before_authentication(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        mind_diary = next(
            plugin for plugin in marketplace["plugins"]
            if plugin["name"] == "mind-diary"
        )

        self.assertEqual(mind_diary["policy"]["installation"], "AVAILABLE")
        self.assertEqual(mind_diary["policy"]["authentication"], "ON_USE")

    def test_manifest_distributes_direct_mcp_without_registered_app(self) -> None:
        manifest = read_json(
            MIND_DIARY_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )

        self.assertNotIn("apps", manifest)
        self.assertEqual(manifest["mcpServers"], "./.mcp.json")
        self.assertFalse((MIND_DIARY_PLUGIN_ROOT / ".app.json").exists())
        self.assertEqual(manifest["interface"]["displayName"], "Mind Diary UAT")
        self.assertIn(
            "restricted Mind Diary UAT pilot",
            manifest["interface"]["longDescription"],
        )
        self.assertRegex(
            manifest["version"],
            re.compile(r"^0\.1\.0\+codex\.\d{14}$"),
        )

    def test_direct_mcp_uses_codex_compatibility_with_canonical_oauth_resource(self) -> None:
        mcp_config = read_json(MIND_DIARY_PLUGIN_ROOT / ".mcp.json")
        mind_diary = mcp_config["mcpServers"]["mind-diary"]

        self.assertEqual(
            mind_diary,
            {
                "type": "http",
                "url": UAT_CODEX_MCP_URL,
                "oauth_resource": UAT_OAUTH_RESOURCE,
            },
        )

    def test_plugin_bundles_path_safe_local_companion_beside_hosted_mcp(self) -> None:
        mcp_config = read_json(MIND_DIARY_PLUGIN_ROOT / ".mcp.json")
        local = mcp_config["mcpServers"]["mind-diary-local"]
        self.assertEqual(
            local,
            {
                "command": "./bin/mind-diary-local-launcher",
                "cwd": ".",
                "env_vars": ["MIND_DIARY_WORKSPACE_ROOTS"],
                "startup_timeout_sec": 10,
                "tool_timeout_sec": 900,
            },
        )
        self.assertNotIn("env", local)
        self.assertNotIn("args", local)
        serialized = json.dumps(mcp_config).lower()
        self.assertNotIn("client_secret", serialized)
        self.assertNotIn("refresh_token", serialized)

        for name in (
            "mind-diary-local-launcher",
            "mind-diary-local-darwin-arm64",
            "mind-diary-local-darwin-amd64",
        ):
            path = MIND_DIARY_PLUGIN_ROOT / "bin" / name
            self.assertTrue(path.is_file(), name)
            self.assertTrue(os.access(path, os.X_OK), name)
            binary = path.read_bytes()
            for forbidden in (
                b"/usr/bin/security",
                b"/usr/bin/open",
                b"client_secret",
                b"refresh_token",
                b"transport-origin",
            ):
                self.assertNotIn(forbidden, binary, name)

        source = MIND_DIARY_PLUGIN_ROOT / "local-companion"
        for name in (
            "go.mod",
            "main.go",
            "config.go",
            "contract.go",
            "protocol.go",
            "local_file.go",
            "upload.go",
        ):
            self.assertTrue((source / name).is_file(), name)
        production_source = "\n".join(
            path.read_text(encoding="utf-8")
            for path in source.glob("*.go")
            if not path.name.endswith("_test.go")
        )
        self.assertNotIn("io.ReadAll", production_source)
        self.assertNotIn("os.ReadFile", production_source)
        self.assertIn("io.LimitReader", production_source)
        self.assertIn("application/octet-stream", production_source)

    def test_packaged_local_companion_lists_stable_two_tool_protocol_offline(self) -> None:
        if platform.system() != "Darwin":
            self.skipTest("packaged local companion currently supports macOS only")
        launcher = MIND_DIARY_PLUGIN_ROOT / "bin" / "mind-diary-local-launcher"
        requests = "\n".join(
            json.dumps(request)
            for request in (
                {
                    "jsonrpc": "2.0",
                    "id": 1,
                    "method": "initialize",
                    "params": {"protocolVersion": "2025-11-25"},
                },
                {"jsonrpc": "2.0", "id": 2, "method": "tools/list"},
            )
        ) + "\n"
        completed = subprocess.run(
            [str(launcher)],
            input=requests,
            text=True,
            capture_output=True,
            check=False,
            timeout=10,
        )
        self.assertEqual(completed.returncode, 0, completed.stderr)
        responses = [json.loads(line) for line in completed.stdout.splitlines()]
        self.assertEqual(responses[0]["result"]["protocolVersion"], "2025-11-25")
        manifest = read_json(
            MIND_DIARY_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )
        self.assertEqual(
            responses[0]["result"]["serverInfo"]["version"],
            manifest["version"],
        )
        tools = responses[1]["result"]["tools"]
        self.assertEqual(
            [tool["name"] for tool in tools],
            ["prepare_local_file", "upload_prepared_file"],
        )
        prepare, upload = tools
        self.assertEqual(prepare["inputSchema"]["required"], ["path"])
        self.assertEqual(
            upload["inputSchema"]["required"],
            ["local_file_ref", "upload_url"],
        )
        self.assertNotIn("path", upload["inputSchema"]["properties"])
        self.assertNotIn("upload_url", upload["outputSchema"]["properties"])
        self.assertNotIn("local_file_ref", upload["outputSchema"]["properties"])

    def test_packaged_binaries_reproduce_from_exact_head_without_vcs_metadata(self) -> None:
        go = shutil.which("go")
        self.assertIsNotNone(go, "Go is required to verify packaged binary provenance")
        binary_names = (
            "mind-diary-local-darwin-arm64",
            "mind-diary-local-darwin-amd64",
        )
        for name in binary_names:
            packaged = MIND_DIARY_PLUGIN_ROOT / "bin" / name
            build_info = subprocess.run(
                [go, "version", "-m", str(packaged)],
                text=True,
                capture_output=True,
                check=False,
                timeout=10,
            )
            self.assertEqual(build_info.returncode, 0, build_info.stderr)
            self.assertNotIn("\tbuild\tvcs", build_info.stdout, name)

        with tempfile.TemporaryDirectory() as directory:
            clean_tree = Path(directory) / "tree"
            clean_tree.mkdir()
            archive = Path(directory) / "source.tar"
            with archive.open("wb") as output:
                archived = subprocess.run(
                    [
                        "git",
                        "archive",
                        "--format=tar",
                        "HEAD",
                        "plugins/mind-diary/local-companion",
                        "plugins/mind-diary/.codex-plugin/plugin.json",
                    ],
                    cwd=REPOSITORY_ROOT,
                    stdout=output,
                    stderr=subprocess.PIPE,
                    check=False,
                    timeout=10,
                )
            self.assertEqual(archived.returncode, 0, archived.stderr.decode())
            with tarfile.open(archive) as source:
                source.extractall(clean_tree, filter="data")

            clean_plugin = clean_tree / "plugins" / "mind-diary"
            clean_manifest = read_json(
                clean_plugin / ".codex-plugin" / "plugin.json"
            )
            rebuilt = Path(directory) / "rebuilt"
            build = subprocess.run(
                [
                    str(clean_plugin / "local-companion" / "build-release.sh"),
                    clean_manifest["version"],
                    str(rebuilt),
                ],
                cwd=clean_plugin / "local-companion",
                text=True,
                capture_output=True,
                check=False,
                timeout=120,
            )
            self.assertEqual(build.returncode, 0, build.stderr)
            for name in binary_names:
                self.assertEqual(
                    sha256_file(rebuilt / name),
                    sha256_file(MIND_DIARY_PLUGIN_ROOT / "bin" / name),
                    f"{name} does not reproduce from the exact HEAD source",
                )

    def test_packaged_launcher_enforces_trusted_workspace_process_config(self) -> None:
        if platform.system() != "Darwin":
            self.skipTest("packaged local companion currently supports macOS only")
        launcher = MIND_DIARY_PLUGIN_ROOT / "bin" / "mind-diary-local-launcher"
        with tempfile.TemporaryDirectory() as workspace, tempfile.TemporaryDirectory() as outside:
            workspace_path = Path(workspace)
            outside_path = Path(outside)
            inside_file = workspace_path / "inside.bin"
            outside_file = outside_path / "outside.bin"
            inside_file.write_bytes(b"inside")
            outside_file.write_bytes(b"outside")
            requests = "\n".join(
                json.dumps(request)
                for request in (
                    {
                        "jsonrpc": "2.0",
                        "id": 1,
                        "method": "initialize",
                        "params": {"protocolVersion": "2025-11-25"},
                    },
                    {
                        "jsonrpc": "2.0",
                        "id": 2,
                        "method": "tools/call",
                        "params": {
                            "name": "prepare_local_file",
                            "arguments": {
                                "path": str(inside_file),
                                "source_kind": "workspace/generated_artifact",
                            },
                        },
                    },
                    {
                        "jsonrpc": "2.0",
                        "id": 3,
                        "method": "tools/call",
                        "params": {
                            "name": "prepare_local_file",
                            "arguments": {
                                "path": str(outside_file),
                                "source_kind": "workspace/generated_artifact",
                            },
                        },
                    },
                    {
                        "jsonrpc": "2.0",
                        "id": 4,
                        "method": "tools/call",
                        "params": {
                            "name": "prepare_local_file",
                            "arguments": {
                                "path": str(outside_file),
                                "source_kind": "local_path",
                            },
                        },
                    },
                )
            ) + "\n"
            environment = os.environ.copy()
            environment["MIND_DIARY_WORKSPACE_ROOTS"] = workspace
            completed = subprocess.run(
                [str(launcher)],
                input=requests,
                text=True,
                capture_output=True,
                check=False,
                timeout=10,
                env=environment,
            )

        self.assertEqual(completed.returncode, 0, completed.stderr)
        responses = {response["id"]: response for response in map(
            json.loads, completed.stdout.splitlines()
        )}
        inside = responses[2]["result"]
        self.assertFalse(inside["isError"])
        self.assertEqual(
            inside["structuredContent"]["source_kind"],
            "workspace/generated_artifact",
        )
        rejected = responses[3]["result"]
        self.assertTrue(rejected["isError"])
        self.assertEqual(
            rejected["structuredContent"]["error"]["code"],
            "file_ingress_source_unsupported",
        )
        local = responses[4]["result"]
        self.assertFalse(local["isError"])
        self.assertEqual(local["structuredContent"]["source_kind"], "local_path")
        serialized = json.dumps(responses)
        self.assertNotIn(str(inside_file), serialized)
        self.assertNotIn(str(outside_file), serialized)

    def test_packaged_launcher_consumes_ref_after_invalid_upload_url(self) -> None:
        if platform.system() != "Darwin":
            self.skipTest("packaged local companion currently supports macOS only")
        launcher = MIND_DIARY_PLUGIN_ROOT / "bin" / "mind-diary-local-launcher"
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory) / "source.bin"
            source.write_bytes(b"definitive invalid URL fixture")
            process = subprocess.Popen(
                [str(launcher)],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
            )
            self.assertIsNotNone(process.stdin)
            self.assertIsNotNone(process.stdout)
            self.assertIsNotNone(process.stderr)

            def request(payload: dict) -> dict:
                process.stdin.write(json.dumps(payload) + "\n")
                process.stdin.flush()
                line = process.stdout.readline()
                self.assertNotEqual(
                    line,
                    "",
                    "packaged launcher closed stdout before returning a response",
                )
                return json.loads(line)

            try:
                request({
                    "jsonrpc": "2.0",
                    "id": 1,
                    "method": "initialize",
                    "params": {"protocolVersion": "2025-11-25"},
                })
                prepared = request({
                    "jsonrpc": "2.0",
                    "id": 2,
                    "method": "tools/call",
                    "params": {
                        "name": "prepare_local_file",
                        "arguments": {"path": str(source)},
                    },
                })
                local_file_ref = prepared["result"]["structuredContent"]["local_file_ref"]
                invalid_upload = {
                    "jsonrpc": "2.0",
                    "id": 3,
                    "method": "tools/call",
                    "params": {
                        "name": "upload_prepared_file",
                        "arguments": {
                            "local_file_ref": local_file_ref,
                            "upload_url": "https://example.com/not-a-mind-diary-intent",
                        },
                    },
                }
                first = request(invalid_upload)
                invalid_upload["id"] = 4
                second = request(invalid_upload)
            finally:
                process.stdin.close()
                process.wait(timeout=10)
                stderr_output = process.stderr.read()
                process.stdout.close()
                process.stderr.close()

        self.assertEqual(process.returncode, 0, stderr_output)
        self.assertEqual(
            first["result"]["structuredContent"]["error"]["code"],
            "invalid_upload_url",
        )
        self.assertEqual(
            second["result"]["structuredContent"]["error"]["code"],
            "local_companion_ref_not_found",
        )
        serialized = json.dumps((first, second))
        self.assertNotIn(str(source), serialized)

    def test_skill_requires_authoritative_bindings_and_explicit_rebind(self) -> None:
        skill = (
            MIND_DIARY_PLUGIN_ROOT / "skills" / "mind-diary" / "SKILL.md"
        ).read_text(encoding="utf-8")

        binding_read = skill.index("Call `get_mind_bindings` before every content")
        content_read = skill.index("## Read workflow")
        self.assertLess(binding_read, content_read)
        self.assertIn("`list_minds` is discovery only", skill)
        self.assertIn("never rebind\nautomatically", skill)
        self.assertIn("Other attached Minds remain read-only", skill)
        self.assertRegex(
            skill,
            re.compile(
                r"Call `commit_changeset`[\s\S]+unchanged active "
                r"`write_binding_id`[\s\S]+`expected_revision`"
            ),
        )
        self.assertIn("principal, token, grant, email\nor internal Mind IDs", skill)

    def test_skill_keeps_automatic_capture_narrow_and_policy_authoritative(self) -> None:
        skill = (
            MIND_DIARY_PLUGIN_ROOT / "skills" / "mind-diary" / "SKILL.md"
        ).read_text(encoding="utf-8")

        capture = skill.index("## Automatic capture workflow")
        boundaries = skill.index("## Product boundaries")
        self.assertLess(capture, boundaries)
        self.assertIn("`automatic_capture.mode` is `routine_non_sensitive`", skill)
        self.assertIn("The Sites control plane is the only place", skill)
        self.assertIn("Never call `capture_knowledge` for credentials", skill)
        self.assertIn("cross-Mind sources", skill)
        self.assertIn("same writable Mind", skill)
        self.assertIn("Treat `captured` as one new immutable revision", skill)
        self.assertIn("`no_op` as successful\n   deduplication", skill)
        self.assertIn("never move, replace, merge or retry the payload", skill)

    def test_skill_routes_one_exact_local_file_through_hosted_intent(self) -> None:
        skill = (
            MIND_DIARY_PLUGIN_ROOT / "skills" / "mind-diary" / "SKILL.md"
        ).read_text(encoding="utf-8")

        local = skill.index("## Local regular-file workflow")
        capture = skill.index("## Automatic capture workflow")
        self.assertLess(local, capture)
        self.assertIn("`prepare_local_file`", skill)
        self.assertIn("`create_file_upload_intent`", skill)
        self.assertIn("`upload_prepared_file`", skill)
        self.assertIn("single `create_bundle_file` or `replace_bundle_file`", skill)
        self.assertIn("Never copy the path into a hosted tool", skill)
        self.assertIn("The companion never deletes or modifies the source file", skill)
        self.assertIn("not directory/bulk/archive", skill)

    def test_skill_routes_incremental_typed_okf_transfer_without_bulk_state(self) -> None:
        skill = (
            MIND_DIARY_PLUGIN_ROOT / "skills" / "mind-diary" / "SKILL.md"
        ).read_text(encoding="utf-8")
        start = skill.index("## Incremental typed OKF transfer")
        end = skill.index("## Local regular-file workflow", start)
        transfer = skill[start:end]

        self.assertLess(
            transfer.index("ordinary local workspace read capability"),
            transfer.index("`commit_changeset`"),
        )
        self.assertIn("explicitly enumerates every entry", transfer)
        self.assertIn(
            "Never infer,\ndiscover or add a related Markdown entry",
            transfer,
        )
        self.assertIn("every transferred entry must be selected by the user", transfer)
        self.assertNotIn("or a small explicitly related set", transfer)
        self.assertIn("`recorded_by`, `applies_to` and `sources`", transfer)
        self.assertIn("Markdown stays Markdown", transfer)
        self.assertRegex(transfer, re.compile(r"`/raw/`.*`/wiki/`.*`/output/`", re.S))
        self.assertIn("`before → after`", transfer)

        confirmation = transfer.index("Apply the existing confirmation rules")
        fresh_bindings = transfer.index("`get_mind_bindings`", confirmation)
        fresh_head = transfer.index("exact current HEAD", fresh_bindings)
        commit = transfer.index("`commit_changeset`", fresh_head)
        self.assertLess(confirmation, fresh_bindings)
        self.assertLess(fresh_bindings, fresh_head)
        self.assertLess(fresh_head, commit)

        attachment_tools = [
            "`prepare_local_file`",
            "`create_file_upload_intent`",
            "`upload_prepared_file`",
        ]
        positions = [transfer.index(tool) for tool in attachment_tools]
        self.assertEqual(positions, sorted(positions))
        for tool in (
            "`replace_index`",
            "`add_log_entry`",
            "`reconcile_file_stage`",
            "`reconcile_changeset`",
            "`get_bundle_file_download`",
            "`validate_mind`",
        ):
            self.assertIn(tool, transfer)

        for excluded in (
            "`AGENTS.md`",
            "`.agents/`",
            "`operations/events`",
            "`operations/revisions`",
            "`operations/config*`",
            "multi-writer or Drive protocol state",
        ):
            self.assertIn(excluded, transfer)
        self.assertIn("must not rewrite the first committed entry", transfer)
        self.assertIn("Never create a migration database", transfer)


if __name__ == "__main__":
    unittest.main()
