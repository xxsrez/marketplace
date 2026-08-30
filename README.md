# Srez Marketplace

This is Andrey's extensible Codex plugin marketplace. It contains independent
Task Manager, Issue Grinder, legacy Ship Tasks, Strategic Explainer, and Mind Diary plugins under
`plugins/`.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/marketplace
```

Then install the plugins you need:

1. Open **Plugins → Task Manager**, **Plugins → Issue Grinder**,
   **Plugins → Strategic Explainer**, or
   **Plugins → Mind Diary UAT** and click
   **Install**.
2. For Task Manager, complete the native OAuth **Connect** step during install
   or upgrade. Mind Diary authenticates when Codex first connects to its MCP
   server. Issue Grinder has no connector of its own and requires Task Manager
   to be installed separately. Strategic Explainer is optional for Issue
   Grinder: it uses ordinary when installed and allowed, and otherwise writes
   natively without
   blocking comments or lifecycle transitions. Task Composer continues to
   require ordinary Strategic Explainer for its Explainer-backed path. Legacy
   Ship Tasks remains available only as a rollback package and should not be
   installed together with Issue Grinder. The Explainer plugin has no connector
   or authentication of its own.
3. Start a new task in Codex after installation or authentication so it loads
   the selected plugin's current skills and tools.

No server URL, client ID, secret, or personal API token is required. Task
Manager and Mind Diary UAT distribute their MCP connections directly. Task
Manager authenticates during install so an upgraded client cannot silently
fall back to an obsolete registered app connector; after that, new server tools
are discovered from the live MCP connection without Developer mode or connector
re-registration. Mind Diary remains a restricted UAT pilot; this package is not
a production or public-directory release.

The plugins connect to:

- Task Manager: `https://task-manager.xxsrez-work.chatgpt.site/api/mcp`
- Mind Diary Codex transport:
  `https://mind-diary.xxsrez-work.chatgpt.site/api/mcp/2025-11-25`
  with canonical OAuth resource
  `https://mind-diary.xxsrez-work.chatgpt.site/api/mcp`

Task Manager is an adapter-only plugin with two coordinated MCP components:

- `task-manager` provides OAuth/MCP discovery, canonical
  references, pagination, safe writes, optimistic concurrency, and native
  comment mechanics.
- `task-manager-local` is a bundled macOS stdio companion with
  `upload_local_file` and `attach_local_file_to_task`. It reads one exact
  host-authorized regular file,
  uploads a verified snapshot into the common `fileRef` workflow, keeps OAuth
  only in process memory after first use, and never sends the full path to Task
  Manager. The local bind tool performs the second staged Agent REST operation.

Issue Grinder is the current skill-only delivery plugin with two independent
Task Manager skills:

- `issue-grinder` delivers selected issue, Release, Project, or current-Release
  scope from `To Do`, `In Progress`, and `In Review` through five execution
  modes. `Соло` uses the current model to perform one issue or packet at a time
  with no subagents and native publication; it can handle one or many issues
  but is never selected from issue count alone. `Классический`, `Баланс` and `Рой` preserve a terminal promise with
  different amounts of economical work; `Экономичный` may instead leave one
  honest resumable candidate without false `Done` or Goal completion. Explicit
  mode selection wins; otherwise any top-level Luna selects `Экономичный` and
  another model selects `Классический`, once per continuous run. It keeps scope
  live. A pure question about modes, the default resolver, their differences or
  selection uses a delivery-free help path without Task Manager, Goal, title
  mutation or subagents. Delivery creates a strategic Goal only for an
  explicitly invoked multi-issue run, comments every non-trivial lifecycle
  transition, and treats `blocked by` as implementation readiness rather than a
  status lock;
- before any terminal blocker, it turns the candidate into a causal explanation
  and performs a fresh reflection over current primary sources. A verified safe
  action cancels the blocker and resumes delivery. A terminal handoff lists all
  confirmed causes and gives a separate answer for each one: why it blocks the
  goal, why Issue Grinder cannot resolve it alone, and what the blocked step
  contributes to the goal. Public UAT is ordinary non-production work;
  Production remains forbidden;
- independent writers and intentional `Рой` candidates use separate feature
  branches and Git worktrees. One integration owner reduces work to one exact
  candidate. `Классический` gives Luna Max only strict-simple packets;
  `Баланс`, `Рой` and `Экономичный` use it more broadly under their mode
  contracts, while `Соло` does not route work to Luna or any other child.
  Before fresh work, a new run inventories related worktrees,
  branches, commits and local changes, then resumes a proven quiescent
  checkpoint instead of creating an accidental parallel replacement;
- `task-composer` keeps the existing planning-only Backlog workflow. Both skills
  depend on the separately installed Task Manager adapter. Issue Grinder uses
  standalone Strategic Explainer when available and otherwise writes natively.

Ship Tasks is the preserved rollback plugin with two coordinated Task Manager
skills:

- `task-composer` owns planning-only formulation and Task Manager backlog
  capture. It preserves one independently deliverable outcome as one Task and
  turns compound work into a problem-first Epic with concrete subtasks, live
  existing Labels, native hierarchy and semantic relations. Unknown current
  Release is omitted rather than guessed. User-provided bug-report attachments
  are preserved as native attachments on the applicable Task when their content
  is relevant to that Task and useful to its executor; screenshots and other
  file types receive no relevance presumption. Failed attachment binding
  remains an explicit partial result. The skill never
  implements the created work or expands Label taxonomy;

- `ship-tasks` owns delivery intent, lifecycle, Goals, verification, releases,
  report content, and terminal status policy. It adaptively fills independent
  safe lanes with subagents only when no user topology rule applies. Natural-
  language rules for exact or relative counts, roles, opt-outs, and conditions
  such as duration or complexity keep their meaning; the root agent is not
  counted as a named subagent. Concurrent implementation writers receive
  separate feature branches and Git worktrees. After interruption or a new
  session, the skill resumes an existing unfinished task-owned worktree or
  branch when prior ownership is safely inactive and exclusive, then rechecks
  the checkpoint instead of restarting the work. Only genuinely simple,
  predictable packets use Luna Max; uncertainty is handed back to the current
  integration owner without a cheap retry. The agent remains free to choose
  planning, tools, implementation, diagnostics, bounded context, and acceptance
  methods while requiring sufficient evidence, a native comment before every
  meaningful lifecycle transition, and factual separation of Task conflicts,
  proven defects, genuine verification blockers, and proven success;
- both skills depend on the separately installed Task Manager plugin for MCP
  tools and authentication, but do not bundle or duplicate that connector;
- Ship Tasks uses `$strategic-explainer:strategic-explainer` when it is installed
  and allowed, and otherwise uses native writing. A selected provider failure
  switches to native.
  Native mode keeps mandatory comments and lifecycle transitions working while
  neither imitating provider methods nor claiming equivalent quality. Task
  Composer continues to use ordinary Strategic Explainer under its own
  contract. Codex manifests do not provide a plugin-to-plugin dependency field,
  so Strategic Explainer remains a separate optional installation for Ship Tasks.

Strategic Explainer is a standalone generic communication plugin with a
semantic facade for one real user-facing comment, report, decision or state
explanation, blocker report, final, or explicit editing request. A client sends
only the purpose, bounded scope, language and material constraints, plus
resolvable read-only source anchors. The facade owns all internal invocation,
isolation, profile, validation, retry, and editing behavior; clients neither
receive nor reproduce that information. It returns publication-ready text with
a separate source basis or operational unavailability and makes no status,
scope, authority, or mutation decisions.

Mentioning `$task-manager` alone still authorizes only the requested adapter
operation. Use `$issue-grinder` or an unambiguous natural-language delivery
request
that names an existing Task or selected Task Manager Project/Release/current
scope. A delivery verb alone does not route ordinary code, product, repository,
or plugin work into Issue Grinder. Create-and-deliver requires an explicit request
to create exactly one Task in Task Manager and immediately start it.

Mind Diary bundles one content skill and two coordinated MCP components. The
hosted content server selects one explicit Mind and revision, while the bundled
macOS `mind-diary-local` companion prepares one exact regular-file path and
streams it through a one-use hosted intent without disclosing the path or
buffering the full file. The workflow preserves immutable history, HEAD CAS and
idempotency on explicit writes. It does not treat ChatGPT/Codex conversation
memory as Mind Diary content.

Repository layout:

- `.agents/plugins/marketplace.json` — the ordered marketplace catalog;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.mcp.json` — direct production MCP and OAuth resource;
- `plugins/task-manager/bin/` and `local-companion/` — bundled macOS local-file
  launcher, binaries, source and tests;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance;
- `plugins/ship-tasks/.codex-plugin/plugin.json` — Ship Tasks plugin manifest;
- `plugins/ship-tasks/skills/ship-tasks/` — Task Manager delivery workflow;
- `plugins/ship-tasks/skills/task-composer/` — Task Manager planning and
  backlog-composition workflow;
- `plugins/issue-grinder/.codex-plugin/plugin.json` — current Issue Grinder
  plugin manifest;
- `plugins/issue-grinder/skills/issue-grinder/` — current Task Manager delivery
  workflow;
- `plugins/issue-grinder/skills/task-composer/` — canonical planning-only
  Task Composer copy;
- `plugins/strategic-explainer/.codex-plugin/plugin.json` — standalone
  Strategic Explainer plugin manifest;
- `plugins/strategic-explainer/skills/strategic-explainer/` — deterministic role
  resolver, provider-only admission, and internal contract used only by an
  admitted fresh terminal provider;
- `plugins/mind-diary/.codex-plugin/plugin.json` — Mind Diary plugin manifest;
- `plugins/mind-diary/.mcp.json` — direct Mind Diary MCP and OAuth resource;
- `plugins/mind-diary/bin/` and `local-companion/` — bundled macOS exact-file
  launcher, binaries, source and tests;
- `plugins/mind-diary/assets/` — Mind Diary brand assets;
- `plugins/mind-diary/skills/mind-diary/` — bounded content workflow guidance.

Each future plugin gets its own `plugins/<plugin-name>/` directory and one
catalog entry. Task Manager remains independently installable as
`task-manager@srez-marketplace`; Issue Grinder and Strategic Explainer are
independently installable as `issue-grinder@srez-marketplace`,
`strategic-explainer@srez-marketplace`.

Validate plugin changes before publishing:

```bash
python3 -m unittest discover -s tests -p 'test_*.py' -v
(cd plugins/task-manager/local-companion && go test -race ./...)
(cd plugins/mind-diary/local-companion && go test -race ./...)
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/task-manager
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/issue-grinder/skills/issue-grinder
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/issue-grinder/skills/task-composer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/ship-tasks
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/task-composer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/strategic-explainer/skills/strategic-explainer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/mind-diary/skills/mind-diary
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/task-manager
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/issue-grinder
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/ship-tasks
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/strategic-explainer
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/mind-diary
jq empty plugins/task-manager/.codex-plugin/plugin.json \
  plugins/task-manager/.mcp.json \
  plugins/issue-grinder/.codex-plugin/plugin.json \
  plugins/ship-tasks/.codex-plugin/plugin.json \
  plugins/strategic-explainer/.codex-plugin/plugin.json \
  plugins/mind-diary/.codex-plugin/plugin.json \
  plugins/mind-diary/.mcp.json \
  .agents/plugins/marketplace.json
! rg -n 'Single-task delivery|Work completion reports|without a Goal|deliver one' \
  plugins/task-manager/skills/task-manager
diff -qr /Users/andrey/Projects/Home/ShipTask/ship-tasks \
  plugins/ship-tasks/skills/ship-tasks
diff -qr /Users/andrey/Projects/Home/ShipTask/issue-grinder \
  plugins/issue-grinder/skills/issue-grinder
diff -qr /Users/andrey/Projects/Home/ShipTask/task-composer \
  plugins/issue-grinder/skills/task-composer
diff -qr /Users/andrey/Projects/Home/ShipTask/strategic-explainer \
  plugins/strategic-explainer/skills/strategic-explainer
git diff --check
```
