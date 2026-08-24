import json
import re
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


if __name__ == "__main__":
    unittest.main()
