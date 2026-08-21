import json
import os
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TASK_MANAGER_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "task-manager"
SHIP_TASKS_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "ship-tasks"
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
        manifest = read_json(
            TASK_MANAGER_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )

        self.assertNotIn("apps", manifest)
        self.assertEqual(manifest["mcpServers"], "./.mcp.json")
        self.assertFalse((TASK_MANAGER_PLUGIN_ROOT / ".app.json").exists())

    def test_direct_mcp_targets_production_oauth_resource(self) -> None:
        mcp_config = read_json(TASK_MANAGER_PLUGIN_ROOT / ".mcp.json")
        task_manager = mcp_config["mcpServers"]["task-manager"]

        self.assertEqual(
            task_manager,
            {
                "type": "http",
                "url": PRODUCTION_MCP_URL,
                "oauth_resource": PRODUCTION_MCP_URL,
            },
        )

    def test_local_companion_is_bundled_without_plaintext_credentials(self) -> None:
        mcp_config = read_json(TASK_MANAGER_PLUGIN_ROOT / ".mcp.json")
        local = mcp_config["mcpServers"]["task-manager-local"]
        self.assertEqual(
            local,
            {
                "command": "./bin/task-manager-local-launcher",
                "cwd": ".",
                "startup_timeout_sec": 10,
                "tool_timeout_sec": 900,
            },
        )
        serialized = json.dumps(mcp_config).lower()
        self.assertNotIn("refresh_token", serialized)
        self.assertNotIn("client_secret", serialized)
        for name in (
            "task-manager-local-launcher",
            "task-manager-local-darwin-arm64",
            "task-manager-local-darwin-amd64",
        ):
            path = TASK_MANAGER_PLUGIN_ROOT / "bin" / name
            self.assertTrue(path.is_file(), name)
            self.assertTrue(os.access(path, os.X_OK), name)
            binary = path.read_bytes()
            for forbidden in (
                b"/usr/bin/security",
                b"serve-private-uat-ingress",
                b"transport-origin",
                b"OAI-Sites-Authorization",
                b"task-manager-uat.xxsrez-work.chatgpt.site",
            ):
                self.assertNotIn(forbidden, binary, name)

        source = TASK_MANAGER_PLUGIN_ROOT / "local-companion"
        self.assertTrue((source / "go.mod").is_file())
        self.assertTrue((source / "main.go").is_file())

    def test_marketplace_has_no_obsolete_uat_ingress_profile(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        names = [plugin["name"] for plugin in marketplace["plugins"]]
        self.assertEqual(names.count("task-manager"), 1)
        self.assertNotIn("task-manager-uat", names)
        self.assertFalse((REPOSITORY_ROOT / "plugins" / "task-manager-uat").exists())

        production_config = read_json(TASK_MANAGER_PLUGIN_ROOT / ".mcp.json")
        self.assertEqual(
            production_config["mcpServers"]["task-manager"]["url"],
            PRODUCTION_MCP_URL,
        )
        self.assertNotIn("task-manager-uat", production_config["mcpServers"])

    def test_task_manager_plugin_is_adapter_only(self) -> None:
        self.assertTrue(
            (
                TASK_MANAGER_PLUGIN_ROOT
                / "skills"
                / "task-manager"
                / "SKILL.md"
            ).is_file()
        )
        self.assertFalse(
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "ship-tasks").exists()
        )
        self.assertFalse(
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "strategic-explainer").exists()
        )

    def test_ship_tasks_is_a_separate_skill_only_plugin(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        names = [plugin["name"] for plugin in marketplace["plugins"]]
        self.assertEqual(names.count("task-manager"), 1)
        self.assertEqual(names.count("ship-tasks"), 1)

        manifest = read_json(
            SHIP_TASKS_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )
        self.assertEqual(manifest["name"], "ship-tasks")
        self.assertEqual(manifest["skills"], "./skills/")
        self.assertNotIn("mcpServers", manifest)
        self.assertNotIn("apps", manifest)
        self.assertFalse((SHIP_TASKS_PLUGIN_ROOT / ".mcp.json").exists())
        self.assertTrue(
            (
                SHIP_TASKS_PLUGIN_ROOT
                / "skills"
                / "ship-tasks"
                / "SKILL.md"
            ).is_file()
        )
        self.assertTrue(
            (
                SHIP_TASKS_PLUGIN_ROOT
                / "skills"
                / "strategic-explainer"
                / "SKILL.md"
            ).is_file()
        )

        metadata = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "agents"
            / "openai.yaml"
        ).read_text(encoding="utf-8")
        self.assertIn('value: "task-manager"', metadata)

    def test_ship_tasks_preserves_the_implicit_routing_gate(self) -> None:
        manifest = read_json(
            SHIP_TASKS_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )
        self.assertEqual(manifest["version"], "0.1.6+codex.20260821111515")
        self.assertIn("selected Task Manager", manifest["description"])
        self.assertIn("A delivery verb alone", manifest["interface"]["longDescription"])
        self.assertTrue(
            any(
                prompt.startswith("Create exactly one Task in Task Manager")
                for prompt in manifest["interface"]["defaultPrompt"]
            )
        )

        skill = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "SKILL.md"
        ).read_text(encoding="utf-8")
        frontmatter = skill.split("---", 2)[1]
        self.assertIn("Одного delivery-глагола недостаточно", frontmatter)
        self.assertIn("ровно одну Task именно в Task Manager", frontmatter)
        self.assertNotIn(
            "по естественным просьбам выполнить, исправить или довести",
            frontmatter,
        )
        for prompt in (
            "Выполни TM-123",
            "Почини X сейчас",
            "Исправь баг в plugin",
            "Реализуй это изменение в коде",
        ):
            self.assertIn(prompt, skill)

    def test_strategic_explainer_is_a_generic_sibling_skill(self) -> None:
        skill = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "SKILL.md"
        ).read_text(encoding="utf-8")
        metadata = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "agents"
            / "openai.yaml"
        ).read_text(encoding="utf-8")

        self.assertIn("name: strategic-explainer", skill)
        self.assertIn("Strategic Handoff", skill)
        self.assertIn("Problem to solve", skill)
        self.assertIn("Current-State Brief", skill)
        self.assertIn("PROBLEM_CONTEXT_ERROR", skill)
        self.assertIn("bounded read-only tools", skill)
        self.assertIn("source note", skill)
        self.assertIn("свободное стратегическое объяснение", skill)
        self.assertIn("CONTEXT_INTEGRITY_ERROR", skill)
        self.assertIn('fork_turns="none"', skill)
        self.assertIn("сформулирует окончательный user-facing текст", skill)
        self.assertNotIn("Technical Brief", skill)
        self.assertNotIn("User Brief", skill)
        self.assertNotIn("PARENT NOTES", skill)
        self.assertNotIn("ShipTask", skill)
        self.assertNotIn("Task Manager", skill)
        self.assertIn('display_name: "Strategic Explainer"', metadata)

    def test_ship_tasks_repairs_contaminated_explainer_invocation(self) -> None:
        skill = (
            SHIP_TASKS_PLUGIN_ROOT / "skills" / "ship-tasks" / "SKILL.md"
        ).read_text(encoding="utf-8")
        handoff = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "references"
            / "strategic-explainer.md"
        ).read_text(encoding="utf-8")

        for text in (skill, handoff):
            self.assertIn('fork_turns="none"', text)
            self.assertIn("CONTEXT_INTEGRITY_ERROR", text)
            self.assertIn("PROBLEM_CONTEXT_ERROR", text)
            self.assertIn("своими словами", text)
            self.assertIn("$ship-tasks:strategic-explainer", text)
        self.assertIn("один раз перезапустить", skill)
        self.assertIn("bounded read-only strategic discovery", skill)
        self.assertIn("local wording fallback", handoff)
        self.assertIn("доступные read-only tools", handoff)
        self.assertNotIn("User Brief` — единственный источник", handoff)


if __name__ == "__main__":
    unittest.main()
