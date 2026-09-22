import json
import os
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TASK_MANAGER_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "task-manager"
ISSUE_GRINDER_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "issue-grinder"
STRATEGIC_EXPLAINER_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "strategic-explainer"
PRODUCTION_MCP_URL = "https://task-manager.xxsrez-work.chatgpt.site/api/mcp"


def read_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


class TaskManagerDirectMcpPackagingTest(unittest.TestCase):
    def test_marketplace_authenticates_direct_mcp_during_install(self) -> None:
        marketplace = read_json(REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json")
        task_manager = next(
            plugin for plugin in marketplace["plugins"] if plugin["name"] == "task-manager"
        )

        self.assertEqual(task_manager["policy"]["authentication"], "ON_INSTALL")

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
        self.assertEqual(
            names,
            [
                "task-manager",
                "mind-diary",
                "strategic-explainer",
                "issue-grinder",
                "interpreter",
            ],
        )
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
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "task-composer").exists()
        )
        self.assertFalse(
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "scope-reviewer").exists()
        )
        self.assertFalse(
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "strategic-explainer").exists()
        )
        self.assertFalse(
            (TASK_MANAGER_PLUGIN_ROOT / "skills" / "issue-grinder").exists()
        )

    def test_issue_grinder_is_a_separate_skill_only_plugin(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        entry = next(
            plugin
            for plugin in marketplace["plugins"]
            if plugin["name"] == "issue-grinder"
        )
        manifest = read_json(
            ISSUE_GRINDER_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )

        self.assertEqual(
            entry["source"],
            {"source": "local", "path": "./plugins/issue-grinder"},
        )
        self.assertEqual(entry["policy"]["installation"], "AVAILABLE")
        self.assertEqual(manifest["name"], "issue-grinder")
        self.assertRegex(manifest["version"], r"^0\.1\.0\+codex\.\d{14}$")
        self.assertEqual(manifest["skills"], "./skills/")
        self.assertNotIn("mcpServers", manifest)
        self.assertNotIn("apps", manifest)
        self.assertFalse((ISSUE_GRINDER_PLUGIN_ROOT / ".mcp.json").exists())

        issue_root = ISSUE_GRINDER_PLUGIN_ROOT / "skills" / "issue-grinder"
        composer_root = ISSUE_GRINDER_PLUGIN_ROOT / "skills" / "task-composer"
        reviewer_root = ISSUE_GRINDER_PLUGIN_ROOT / "skills" / "scope-reviewer"
        self.assertTrue((issue_root / "SKILL.md").is_file())
        self.assertTrue(
            (issue_root / "references" / "thread-title.md").is_file()
        )
        self.assertTrue(
            (issue_root / "references" / "execution-modes.md").is_file()
        )
        mode_files = {
            "Соло": "solo.md",
            "Классический": "classic.md",
            "Баланс": "balance.md",
            "Экономичный": "economical.md",
        }
        for filename in mode_files.values():
            self.assertTrue(
                (issue_root / "references" / "modes" / filename).is_file()
            )
        self.assertTrue((issue_root / "references" / "mode-help.md").is_file())
        self.assertTrue((issue_root / "references" / "run-and-goal.md").is_file())
        self.assertTrue((issue_root / "scripts" / "model_routing_guard.py").is_file())
        self.assertTrue((composer_root / "SKILL.md").is_file())
        self.assertTrue((reviewer_root / "SKILL.md").is_file())
        self.assertTrue(
            (reviewer_root / "scripts" / "lens_routing_guard.py").is_file()
        )
        self.assertFalse(
            (ISSUE_GRINDER_PLUGIN_ROOT / "skills" / "ship-tasks").exists()
        )
        self.assertFalse(
            (ISSUE_GRINDER_PLUGIN_ROOT / "skills" / "strategic-explainer").exists()
        )

        skill = (issue_root / "SKILL.md").read_text(encoding="utf-8")
        title_contract = (
            issue_root / "references" / "thread-title.md"
        ).read_text(encoding="utf-8")
        execution_modes = (
            issue_root / "references" / "execution-modes.md"
        ).read_text(encoding="utf-8")
        mode_help = (issue_root / "references" / "mode-help.md").read_text(
            encoding="utf-8"
        )
        metadata = (issue_root / "agents" / "openai.yaml").read_text(
            encoding="utf-8"
        )
        composer = (composer_root / "SKILL.md").read_text(encoding="utf-8")
        reviewer = (reviewer_root / "SKILL.md").read_text(encoding="utf-8")
        reviewer_metadata = (reviewer_root / "agents" / "openai.yaml").read_text(
            encoding="utf-8"
        )
        reviewer_invocation = (
            reviewer_root / "references" / "invocation-and-live-run.md"
        ).read_text(encoding="utf-8")
        runtime = " ".join(
            path.read_text(encoding="utf-8")
            for path in sorted(issue_root.rglob("*.md"))
        )
        normalized = " ".join(runtime.split())

        self.assertIn("name: issue-grinder", skill)
        self.assertIn("candidate blocker → причинное объяснение → reflection", runtime)
        self.assertIn("для каждой причины дай отдельный ответ", normalized)
        self.assertIn(
            "почему она блокирует обязательный результат активного issue, "
            "почему Issue Grinder не может устранить её сам и что "
            "заблокированный шаг даст issue contract и общей цели",
            normalized,
        )
        self.assertIn("не превращай в источник новых Requirements", skill)
        self.assertIn("empty active scope остаётся достаточным", runtime)
        self.assertIn("update_goal(status=blocked)", runtime)
        self.assertIn("Production запрещён полностью", skill)
        self.assertIn("публичный UAT", runtime)
        self.assertIn("`Экономичный` только явно", normalized)
        self.assertIn("не пересчитывай его", normalized)
        self.assertIn("## Выбранный режим — обязательная загрузка", execution_modes)
        for mode, filename in mode_files.items():
            mode_contract = (
                issue_root / "references" / "modes" / filename
            ).read_text(encoding="utf-8")
            self.assertTrue(mode_contract.startswith(f"# {mode}\n"))
            self.assertIn(f"modes/{filename}", execution_modes)
        self.assertNotIn("## Соло", execution_modes)
        economical = (
            issue_root / "references" / "modes" / "economical.md"
        ).read_text(encoding="utf-8")
        self.assertIn("resumable checkpoint", economical)
        routing_guard = (
            issue_root / "scripts" / "model_routing_guard.py"
        ).read_text(encoding="utf-8")
        self.assertIn("Все содержательные решения и работа режима выполняются Luna Max", economical)
        self.assertIn("issue-grinder/model-routing/v2", routing_guard)
        self.assertIn('parser.add_argument("--packet-id", required=True)', routing_guard)
        self.assertIn("dispatch_fingerprint", routing_guard)
        self.assertIn("actual_luna_model_mismatch", routing_guard)
        self.assertNotIn("FORCED_PROFILE_AGENT_TYPES", routing_guard)
        self.assertNotIn("platform_agent_type_bypasses_mode_profile", routing_guard)
        self.assertIn("multi-agent-routing.md", execution_modes)
        multi_agent_routing = (
            issue_root / "references" / "multi-agent-routing.md"
        ).read_text(encoding="utf-8")
        self.assertIn("Имя или тип агента не выбирает профиль режима", multi_agent_routing)
        for name in ("epic-planning.md", "epic-publication.md", "attachments.md"):
            self.assertTrue((composer_root / "references" / name).is_file())
        self.assertIn("Issue Grinder · ...", runtime)
        self.assertIn(
            "set_thread_title` не более одного раза без", title_contract
        )
        self.assertIn("Meaningful title", title_contract)
        self.assertIn("только пользователю в чате", normalized)
        self.assertIn("[краткую справку](references/mode-help.md)", skill)
        self.assertIn("`Экономичный` включается только явно", mode_help)
        self.assertIn("не обращается к Task Manager", mode_help)
        self.assertIn("Sol/controller делает почти всю работу сам", mode_help)
        self.assertIn('value: "task-manager"', metadata)
        self.assertIn("allow_implicit_invocation: true", metadata)
        self.assertIn("$issue-grinder:task-composer", composer)
        self.assertIn("name: scope-reviewer", reviewer)
        self.assertIn("$issue-grinder:scope-reviewer", reviewer)
        self.assertIn('model="gpt-5.6-luna"', reviewer)
        self.assertIn('reasoning_effort="max"', reviewer)
        self.assertIn("Human Requirements\nне изменяй ни при каких обстоятельствах", reviewer)
        self.assertIn('value: "task-manager"', reviewer_metadata)
        self.assertIn("allow_implicit_invocation: true", reviewer_metadata)
        self.assertIn("ревью плана перед запуском", reviewer)
        self.assertIn("активного долгого Issue Grinder", reviewer)
        self.assertIn("Pre-launch review", reviewer_invocation)
        self.assertIn("Active long-run review", reviewer_invocation)
        self.assertIn("передай управление owning Issue Grinder", reviewer_invocation)

        public_manifest = json.dumps(manifest["interface"], ensure_ascii=False)
        self.assertIn("fresh reflection over current primary sources", public_manifest)
        self.assertIn("A terminal handoff lists every confirmed blocking cause", public_manifest)
        self.assertIn("why Issue Grinder cannot resolve it alone", public_manifest)
        self.assertIn("Public UAT", public_manifest)
        self.assertIn("Production remains forbidden", public_manifest)
        self.assertIn("four execution modes", public_manifest)
        self.assertIn(
            "Solo uses the current main profile without execution children",
            public_manifest,
        )
        self.assertIn(
            "Mode topology includes only agents doing Issue Grinder delivery work",
            public_manifest,
        )
        self.assertIn(
            "For ordinary comments and successful final reports, Issue Grinder uses the standalone Strategic Explainer",
            public_manifest,
        )
        self.assertIn("exact gpt-6-luna/max selects Balance regardless of scope size", public_manifest)
        self.assertIn("Other profiles select Solo for small and medium scopes, and Classic for justified large scopes", public_manifest)
        self.assertIn("once per continuous run", public_manifest)
        self.assertIn("delivery-free help path", public_manifest)
        self.assertIn("Classic keeps most implementation on its controller", public_manifest)
        self.assertIn(
            "Economical is explicit-only",
            public_manifest,
        )
        self.assertIn("Includes three independent Task Manager skills", public_manifest)
        self.assertIn("Scope Reviewer turns an exact selected plan", public_manifest)
        self.assertIn("Plan review and every Release review remain read-only", public_manifest)
        self.assertIn("never Human Requirements", public_manifest)
        self.assertIn("plans before launch or active long runs", public_manifest)

    def test_retired_ship_tasks_is_absent(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        names = [plugin["name"] for plugin in marketplace["plugins"]]
        self.assertNotIn("ship-tasks", names)
        self.assertFalse((REPOSITORY_ROOT / "plugins" / "ship-tasks").exists())

    def test_strategic_explainer_is_a_standalone_generic_plugin(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        entry = next(
            plugin
            for plugin in marketplace["plugins"]
            if plugin["name"] == "strategic-explainer"
        )
        manifest = read_json(
            STRATEGIC_EXPLAINER_PLUGIN_ROOT / ".codex-plugin" / "plugin.json"
        )
        skill = (
            STRATEGIC_EXPLAINER_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "SKILL.md"
        ).read_text(encoding="utf-8")
        metadata = (
            STRATEGIC_EXPLAINER_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "agents"
            / "openai.yaml"
        ).read_text(encoding="utf-8")
        provider = (
            STRATEGIC_EXPLAINER_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "references"
            / "provider-contract.md"
        ).read_text(encoding="utf-8")
        entrypoint = (
            STRATEGIC_EXPLAINER_PLUGIN_ROOT
            / "skills"
            / "strategic-explainer"
            / "references"
            / "provider-entrypoint.md"
        ).read_text(encoding="utf-8")

        self.assertEqual(
            entry["source"],
            {"source": "local", "path": "./plugins/strategic-explainer"},
        )
        self.assertEqual(manifest["name"], "strategic-explainer")
        self.assertEqual(manifest["skills"], "./skills/")
        self.assertRegex(manifest["version"], r"^0\.1\.\d+\+codex\.\d{14}$")
        self.assertNotIn("mcpServers", manifest)
        self.assertNotIn("apps", manifest)
        public_manifest = json.dumps(manifest["interface"], ensure_ascii=False)
        for internal_marker in (
            "fork_turns",
            "STRATEGIC_EXPLAINER_PROVIDER_V1",
            "reasoning_effort",
            "главную причинную мысль",
            "Проверь понимание",
            "Редакторская реконструкция",
        ):
            self.assertNotIn(internal_marker, public_manifest)
        self.assertFalse((STRATEGIC_EXPLAINER_PLUGIN_ROOT / ".mcp.json").exists())
        self.assertIn("name: strategic-explainer", skill)
        self.assertIn("семантический facade к изолированному provider-subagent", skill)
        self.assertIn("текущий агент исполняет facade router", skill)
        self.assertIn("references/provider-entrypoint.md", skill)
        self.assertIn("references/provider-contract.md", skill)
        self.assertIn('fork_turns="none"', skill)
        self.assertIn('model="gpt-5.6-luna"', skill)
        self.assertIn('reasoning_effort="max"', skill)
        self.assertIn("Не наследуй current model/effort", skill)
        self.assertNotIn("gpt-6-luna", metadata)
        self.assertNotIn("gpt-5.6-luna", metadata)
        self.assertNotIn("reasoning_effort", metadata)
        self.assertNotIn("STRATEGIC_EXPLAINER_PROVIDER_V1", metadata)
        self.assertIn("одна реальная user-facing formulation", skill)
        self.assertIn("Сразу полностью прочитай", skill)
        self.assertIn("Routine", skill)
        self.assertIn("publication-ready text", skill)
        self.assertIn("STRATEGIC_EXPLAINER_PROVIDER_V1", skill)
        self.assertIn("Никогда не вызывай Strategic Explainer", entrypoint)
        self.assertIn("STRATEGIC_EXPLAINER_INVOCATION_ERROR", entrypoint)
        self.assertIn("не читай `provider-contract.md`", entrypoint)
        self.assertNotIn("одну главную причинную мысль", skill)
        self.assertNotIn("первый смысловой слой", skill)
        self.assertNotIn("неизменяемое смысловое ядро", skill)
        self.assertIn("Установи исходный вопрос и факты", provider)
        self.assertIn("Собери strategic context снизу вверх", provider)
        self.assertIn("одну главную причинную мысль", provider)
        self.assertIn("Проверь понимание", provider)
        self.assertIn("Редакторская реконструкция", provider)
        self.assertIn("неизменяемое ядро", provider)
        self.assertNotIn("PROBLEM_CONTEXT_ERROR", skill)
        self.assertNotIn("CONTEXT_INTEGRITY_ERROR", skill)
        self.assertNotIn("2–4 реально", skill)
        self.assertNotIn("Technical Brief", skill)
        self.assertNotIn("User Brief", skill)
        self.assertNotIn("PARENT NOTES", skill)
        self.assertNotIn("ShipTask", skill)
        self.assertNotIn("Task Manager", skill)
        self.assertIn('display_name: "Strategic Explainer"', metadata)



if __name__ == "__main__":
    unittest.main()
