import json
import os
import platform
import re
import subprocess
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
MIND_DIARY_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "mind-diary"
UAT_OAUTH_RESOURCE = "https://mind-diary.xxsrez-work.chatgpt.site/api/mcp"
UAT_CODEX_MCP_URL = f"{UAT_OAUTH_RESOURCE}/2025-11-25"


def read_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


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


if __name__ == "__main__":
    unittest.main()
