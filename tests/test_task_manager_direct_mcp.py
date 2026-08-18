import json
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "task-manager"
PRODUCTION_MCP_URL = "https://task-manager.xxsrez-work.chatgpt.site/api/mcp"


def read_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


class TaskManagerDirectMcpPackagingTest(unittest.TestCase):
    def test_marketplace_defers_authentication_until_first_use(self) -> None:
        marketplace = read_json(REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json")
        task_manager = next(
            plugin for plugin in marketplace["plugins"] if plugin["name"] == "task-manager"
        )

        self.assertEqual(task_manager["policy"]["authentication"], "ON_USE")

    def test_manifest_distributes_direct_mcp_without_registered_app(self) -> None:
        manifest = read_json(PLUGIN_ROOT / ".codex-plugin" / "plugin.json")

        self.assertNotIn("apps", manifest)
        self.assertEqual(manifest["mcpServers"], "./.mcp.json")
        self.assertFalse((PLUGIN_ROOT / ".app.json").exists())

    def test_direct_mcp_targets_production_oauth_resource(self) -> None:
        mcp_config = read_json(PLUGIN_ROOT / ".mcp.json")
        task_manager = mcp_config["mcpServers"]["task-manager"]

        self.assertEqual(
            task_manager,
            {
                "type": "http",
                "url": PRODUCTION_MCP_URL,
                "oauth_resource": PRODUCTION_MCP_URL,
            },
        )

    def test_both_task_manager_skills_are_bundled(self) -> None:
        for skill_name in ("task-manager", "ship-tasks"):
            self.assertTrue(
                (PLUGIN_ROOT / "skills" / skill_name / "SKILL.md").is_file(),
                f"missing bundled skill: {skill_name}",
            )


if __name__ == "__main__":
    unittest.main()
