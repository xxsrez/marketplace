import json
import re
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
MIND_DIARY_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "mind-diary"
UAT_MCP_URL = "https://mind-diary.xxsrez-work.chatgpt.site/api/mcp"


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

    def test_direct_mcp_targets_the_exact_uat_oauth_resource(self) -> None:
        mcp_config = read_json(MIND_DIARY_PLUGIN_ROOT / ".mcp.json")
        mind_diary = mcp_config["mcpServers"]["mind-diary"]

        self.assertEqual(
            mind_diary,
            {
                "type": "http",
                "url": UAT_MCP_URL,
                "oauth_resource": UAT_MCP_URL,
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


if __name__ == "__main__":
    unittest.main()
