import json
import os
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TASK_MANAGER_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "task-manager"
SHIP_TASKS_PLUGIN_ROOT = REPOSITORY_ROOT / "plugins" / "ship-tasks"
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
                "ship-tasks",
                "strategic-explainer",
                "issue-grinder",
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
            "Рой": "swarm.md",
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
        self.assertIn("top-level `gpt-5.6-luna` при любом effort", normalized)
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
        balance = (
            issue_root / "references" / "modes" / "balance.md"
        ).read_text(encoding="utf-8")
        swarm = (
            issue_root / "references" / "modes" / "swarm.md"
        ).read_text(encoding="utf-8")
        routing_guard = (
            issue_root / "scripts" / "model_routing_guard.py"
        ).read_text(encoding="utf-8")
        self.assertIn("До первой source mutation", balance)
        self.assertIn("Best-of-M wave с `M >= 2`", swarm)
        self.assertIn("Все содержательные решения и работа режима выполняются Luna Max", economical)
        self.assertIn("issue-grinder/model-routing/v1", routing_guard)
        self.assertIn("actual_luna_model_mismatch", routing_guard)
        self.assertNotIn("FORCED_PROFILE_AGENT_TYPES", routing_guard)
        self.assertNotIn("platform_agent_type_bypasses_mode_profile", routing_guard)
        self.assertIn("Имя или тип агента не выбирает профиль режима", execution_modes)
        self.assertIn("Issue Grinder · ...", runtime)
        self.assertIn(
            "set_thread_title` не более одного раза без", title_contract
        )
        self.assertIn("Meaningful title", title_contract)
        self.assertIn("только пользователю в чате", normalized)
        self.assertIn("[краткую справку](references/mode-help.md)", skill)
        self.assertIn("`По умолчанию` — не шестой режим", mode_help)
        self.assertIn("не обращается к Task Manager", mode_help)
        self.assertIn("Sol/controller: он делает почти всё", mode_help)
        self.assertIn(
            "основной объём лёгкой и средней ограниченной implementation",
            mode_help,
        )
        self.assertIn("возвращает пакет Sol/controller-у", mode_help)
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

        public_manifest = json.dumps(manifest["interface"], ensure_ascii=False)
        self.assertIn("fresh reflection over current primary sources", public_manifest)
        self.assertIn("A terminal handoff lists every confirmed blocking cause", public_manifest)
        self.assertIn("why Issue Grinder cannot resolve it alone", public_manifest)
        self.assertIn("Public UAT", public_manifest)
        self.assertIn("Production remains forbidden", public_manifest)
        self.assertIn("five stable execution modes", public_manifest)
        self.assertIn(
            "Solo keeps all Issue Grinder delivery work in one execution lane",
            public_manifest,
        )
        self.assertIn(
            "Mode topology includes only agents doing Issue Grinder delivery work",
            public_manifest,
        )
        self.assertIn(
            "In all five modes Issue Grinder uses the standalone Strategic Explainer",
            public_manifest,
        )
        self.assertIn("issue count alone does not select it", public_manifest)
        self.assertIn("once per continuous run", public_manifest)
        self.assertIn("delivery-free help path", public_manifest)
        self.assertIn("controller does almost all work", public_manifest)
        self.assertIn("preferred executor for light and medium", public_manifest)
        self.assertIn("returns problems to the controller", public_manifest)
        self.assertIn("Includes three independent Task Manager skills", public_manifest)
        self.assertIn("Scope Reviewer turns an exact selected plan", public_manifest)
        self.assertIn("Plan review and every Release review remain read-only", public_manifest)
        self.assertIn("never Human Requirements", public_manifest)

    def test_ship_tasks_is_a_separate_skill_only_plugin(self) -> None:
        marketplace = read_json(
            REPOSITORY_ROOT / ".agents" / "plugins" / "marketplace.json"
        )
        names = [plugin["name"] for plugin in marketplace["plugins"]]
        self.assertEqual(names.count("task-manager"), 1)
        self.assertEqual(names.count("ship-tasks"), 1)
        self.assertEqual(names.count("strategic-explainer"), 1)

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
                / "task-composer"
                / "SKILL.md"
            ).is_file()
        )
        self.assertFalse(
            (SHIP_TASKS_PLUGIN_ROOT / "skills" / "strategic-explainer").exists()
        )
        self.assertFalse(
            (SHIP_TASKS_PLUGIN_ROOT / "skills" / "scope-reviewer").exists()
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
        self.assertRegex(manifest["version"], r"^0\.1\.\d+\+codex\.\d{14}$")
        self.assertIn("selected Task Manager", manifest["description"])
        self.assertIn("A delivery verb alone", manifest["interface"]["longDescription"])
        self.assertIn("preserves natural-language subagent rules", manifest["interface"]["longDescription"])
        self.assertIn("resumes unfinished task-owned Git work", manifest["interface"]["longDescription"])
        self.assertIn("one fresh read-only critic", manifest["interface"]["longDescription"])
        self.assertIn("stale or inconclusive reviews remain In Review", manifest["interface"]["longDescription"])
        self.assertIn("selected provider result as reflection input", manifest["interface"]["longDescription"])
        self.assertIn("ordinary `$strategic-explainer:strategic-explainer` is installed", manifest["interface"]["longDescription"])
        self.assertIn("pass it only the purpose, bounded scope", manifest["interface"]["longDescription"])
        self.assertIn("facade owns all internal execution and editing behavior", manifest["interface"]["longDescription"])
        self.assertIn("independent optional installation", manifest["interface"]["longDescription"])
        self.assertIn("clients neither receive nor reproduce those instructions", manifest["interface"]["longDescription"])
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
        self.assertIn("одного delivery-глагола недостаточно", frontmatter.lower())
        self.assertIn("создать и сразу выполнить одну Task", frontmatter)
        self.assertNotIn(
            "по естественным просьбам выполнить, исправить или довести",
            frontmatter,
        )
        self.assertIn("## 1. Выбери mode и exact scope", skill)
        self.assertIn("Обычная просьба исправить код/продукт", skill)
        self.assertIn("planning/backlog capture с Task", skill)
        self.assertIn("принадлежит Task Composer", skill)
        self.assertIn("current parent chain до ближайшего relevant Epic", skill)
        self.assertIn("bounded context любому implementation/review packet", skill)
        self.assertIn("не расширяет selector/scope", skill)

    def test_task_composer_plans_without_delivery_or_taxonomy_mutation(self) -> None:
        root = SHIP_TASKS_PLUGIN_ROOT / "skills" / "task-composer"
        skill = (root / "SKILL.md").read_text(encoding="utf-8")
        metadata = (root / "agents" / "openai.yaml").read_text(encoding="utf-8")
        normalized_skill = " ".join(skill.lower().split())

        self.assertIn("name: task-composer", skill)
        self.assertIn("canonical status `Backlog`", skill)
        self.assertIn("current unreleased Release", skill)
        self.assertIn("Не создавай Epic с одной формальной подзадачей", skill)
        self.assertIn("$strategic-explainer:strategic-explainer", skill)
        self.assertIn("как semantic facade отдельной", skill)
        self.assertIn("назначение description", skill)
        self.assertNotIn('model="gpt-5.6-luna"', skill)
        self.assertNotIn('reasoning_effort="max"', skill)
        self.assertNotIn("STRATEGIC_EXPLAINER_PROVIDER_V1", skill)
        self.assertIn("Не превращай шаги исходного плана в Tasks механически", skill)
        self.assertIn("самодостаточную проекцию", skill)
        self.assertIn("Epic context не расширяет scope child", skill)
        self.assertIn("taxonomy не расширяй", skill)
        self.assertIn("не строй последовательную цепочку", normalized_skill)
        self.assertIn("bounded duplicate search", skill)
        self.assertIn("Type/classification выражай native Label", skill)
        self.assertIn("`BUG:`, `EPIC:`, `[Bug]`, `Epic —`, `Feature:`", skill)
        self.assertIn(
            "Legacy-prefixed и clean outcome title считай одним duplicate candidate",
            skill,
        )
        self.assertIn("не дублирует type/classification label", normalized_skill)
        self.assertIn("не повторяй create вслепую", normalized_skill)
        self.assertIn(
            "формат сам по себе не создаёт презумпцию уместности",
            skill,
        )
        self.assertIn(
            "один критерий к screenshot, документу, логу, записи",
            normalized_skill,
        )
        self.assertNotIn("считай осмысленным по умолчанию", skill)
        self.assertIn("native attachment самой конкретной создаваемой Task", skill)
        self.assertIn("attachment disposition", skill)
        self.assertIn("каждый обязательный attachment", normalized_skill)
        self.assertIn("не реализуй", skill.lower())
        self.assertIn('display_name: "Task Composer"', metadata)
        self.assertIn('value: "task-manager"', metadata)
        self.assertIn("allow_implicit_invocation: true", metadata)

    def test_ship_tasks_classifies_review_outcomes_without_tool_choreography(self) -> None:
        skill = (
            SHIP_TASKS_PLUGIN_ROOT / "skills" / "ship-tasks" / "SKILL.md"
        ).read_text(encoding="utf-8")
        report = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "references"
            / "delivery-report.md"
        ).read_text(encoding="utf-8")
        handoff = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "references"
            / "strategic-explainer.md"
        ).read_text(encoding="utf-8")
        run_report = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "references"
            / "run-report.md"
        ).read_text(encoding="utf-8")
        critical_review = (
            SHIP_TASKS_PLUGIN_ROOT
            / "skills"
            / "ship-tasks"
            / "references"
            / "critical-codebase-review.md"
        ).read_text(encoding="utf-8")

        for outcome in (
            "task-contract-conflict",
            "verified-success",
            "verified-failure",
            "verification-blocked",
            "critical-codebase-accepted",
        ):
            self.assertIn(outcome, skill)
        self.assertIn("[critical codebase review](references/critical-codebase-review.md)", skill)
        self.assertIn("`To Do == 0`, `In Progress == 0`, `In Review > 0`", critical_review)
        self.assertIn("человека как содержательного verifier", critical_review)
        self.assertIn("bounded unlocker", critical_review)
        self.assertIn("ровно одного read-only subagent role `critic`", critical_review)
        self.assertIn('`fork_turns="none"`', critical_review)
        self.assertIn("Mere absence of findings", critical_review)
        self.assertIn("stale", critical_review)
        self.assertIn("residual", critical_review)
        self.assertIn("In Review → In Progress", skill)
        self.assertIn("batch-implementation`, с Goal", skill)
        self.assertIn("минимум двух", skill)
        self.assertIn("сам release новый Goal не создаёт", skill)
        self.assertIn("Release-only run Goal не создаёт", skill)
        self.assertNotIn("запускает batch по project memory", skill)
        self.assertIn("не доказывает defect каждой", skill)
        self.assertIn("Перед любым существенным status transition", skill)
        self.assertIn("Доказательство важнее выбранного способа", skill)
        self.assertIn("Сбой одного выбранного", skill)
        self.assertIn("не делает его обязательным", skill)
        self.assertIn("как обеспечить это, решает агент", skill)
        self.assertIn("продолжай rework в этом же", skill)
        self.assertIn("Native comment create/list/read", skill)
        self.assertIn("Всегда создай и перечитай обязательный comment", skill)
        self.assertIn("немедленно сообщи в Codex chat", skill)
        self.assertIn("opening Task comment", skill)
        self.assertIn("каждые 10 минут", skill)
        self.assertIn("краткий перечень существенных инцидентов", skill)
        self.assertNotIn("сначала восстанови", skill)
        self.assertNotIn("после material repair", skill)
        self.assertNotIn("communication remainder", report)
        self.assertIn("Инвариант effects", report)
        self.assertIn("До связанного существенного status transition", report)
        self.assertIn("Client protocol выбора Strategic Explainer", handoff)
        self.assertIn("Для каждой publication unit сделай один semantic call", handoff)
        self.assertIn(
            "внутренним исполнением полностью владеет facade Strategic Explainer",
            " ".join(handoff.split()),
        )
        self.assertIn(
            "Обычный `To Do → In Progress` не запускает publication unit",
            " ".join(skill.split()),
        )
        self.assertNotIn("2–4 способа", handoff)
        self.assertIn("найденный и устранённый", run_report)
        self.assertIn("не совместим с формулировкой полного успеха", run_report)
        self.assertIn("authoritative source anchors и factual inventory", run_report)
        self.assertIn("без caller rewrite", run_report)
        for retired in (
            "строгого tool threshold",
            "строгого model-tool threshold",
            "blocker occurrence",
            "self-recovery",
        ):
            self.assertNotIn(retired, skill)

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

    def test_ship_tasks_honors_natural_language_topology_and_resume(self) -> None:
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
        normalized_skill = " ".join(skill.split())
        normalized_handoff = " ".join(handoff.split())

        self.assertIn("$strategic-explainer:strategic-explainer", skill)
        self.assertIn("$strategic-explainer:strategic-explainer", handoff)
        self.assertIn("Сначала выведи effective topology rule", skill)
        self.assertIn("Exact count — обязательное число", skill)
        self.assertIn("Root не входит в явно названное число субагентов", normalized_skill)
        self.assertIn("условие по длительности, сложности", skill)
        self.assertIn("Только без применимого rule сам решай", skill)
        self.assertIn("единственный integration owner", skill)
        self.assertIn("Task Manager comments/status/version writes", normalized_skill)
        self.assertIn("собственную feature branch и собственный Git worktree", normalized_skill)
        self.assertIn("Один writable worktree принадлежит одному writer", skill)
        self.assertIn("Общий no-subagent rule означает ноль субагентов", normalized_skill)
        self.assertIn("В начале run выбери communication mode", skill)
        self.assertIn("разрешён, иначе native", normalized_skill)
        self.assertIn("communication mode становится native", normalized_skill)
        self.assertIn("не делай второй editorial rewrite", normalized_skill)
        self.assertIn("exact unfinished worktree/branch существует", normalized_skill)
        self.assertIn("active/unknown ownership не перехватывай", normalized_skill)
        self.assertIn(
            "Обычный `To Do → In Progress` не запускает publication unit",
            normalized_skill,
        )
        self.assertIn("Client protocol выбора Strategic Explainer", handoff)
        self.assertNotIn('model="gpt-5.6-luna"', handoff)
        self.assertNotIn('reasoning_effort="max"', handoff)
        self.assertNotIn("STRATEGIC_EXPLAINER_PROVIDER_V1", handoff)
        self.assertIn("semantic call", handoff)
        self.assertIn("Матрица выбора", normalized_handoff)
        self.assertIn("Opt-out", normalized_handoff)
        self.assertNotIn("одну главную причинную мысль", handoff)
        self.assertNotIn("Проверка понимания", handoff)
        self.assertEqual(skill.count('fork_turns="none"'), 1)
        self.assertEqual(handoff.count('fork_turns="none"'), 0)
        self.assertIn("reflection input", normalized_handoff)
        self.assertNotIn("не вызывать отдельно", handoff)
        self.assertNotIn("сначала восстанови", skill)

    def test_ship_tasks_titles_only_a_proven_fresh_codex_task(self) -> None:
        root = SHIP_TASKS_PLUGIN_ROOT / "skills" / "ship-tasks"
        skill = (root / "SKILL.md").read_text(encoding="utf-8")
        title_contract = (root / "references" / "thread-title.md").read_text(
            encoding="utf-8"
        )
        normalized_contract = " ".join(title_contract.lower().split())

        self.assertIn("[title contract](references/thread-title.md)", skill)
        self.assertIn("ровно одного кандидата calling task", title_contract)
        self.assertIn("не выбирай просто самый свежий task", normalized_contract)
        self.assertIn("history не paginated", title_contract)
        self.assertIn("до первой Task Manager mutation", title_contract)
        self.assertIn("без `threadId`", title_contract)
        self.assertIn("не передавай discovery candidate id", title_contract)
        self.assertIn("task-title=renamed", title_contract)
        self.assertIn("task-title=preserved", title_contract)
        self.assertIn("task-title=not-available", title_contract)


if __name__ == "__main__":
    unittest.main()
