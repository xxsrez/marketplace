# Srez Marketplace

This is Andrey's extensible Codex plugin marketplace. It contains independent
Task Manager, Ship Tasks, Strategic Explainer, Strategic Explainer Fast, and Mind Diary plugins under
`plugins/`.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/marketplace
```

Then install the plugins you need:

1. Open **Plugins → Task Manager**, **Plugins → Ship Tasks**,
   **Plugins → Strategic Explainer**, **Plugins → Strategic Explainer Fast**, or
   **Plugins → Mind Diary UAT** and click
   **Install**.
2. For Task Manager, complete the native OAuth **Connect** step during install
   or upgrade. Mind Diary authenticates when Codex first connects to its MCP
   server. Ship Tasks has no connector of its own and requires both Task Manager
   and one Strategic Explainer to be installed separately. Ship Tasks prefers
   Fast when it is installed and otherwise uses ordinary Strategic Explainer;
   Task Composer continues to require ordinary Strategic Explainer. Both
   Explainer plugins have no connector or authentication of their own.
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

Ship Tasks is a separate plugin with two coordinated Task Manager skills:

- `task-composer` owns planning-only formulation and Task Manager backlog
  capture. It preserves one independently deliverable outcome as one Task and
  turns compound work into a problem-first Epic with concrete subtasks, live
  existing Labels, native hierarchy and semantic relations. Unknown current
  Release is omitted rather than guessed. Meaningful user-provided bug-report
  attachments are preserved as native attachments on the applicable Task;
  screenshots showing the reported bug are relevant by default, and failed
  attachment binding remains an explicit partial result. The skill never
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
- Ship Tasks prefers
  `$strategic-explainer-fast:strategic-explainer-fast` for mandatory
  explanations and falls back to `$strategic-explainer:strategic-explainer`
  only when Fast is absent. It never calls both for one publication unit. Task
  Composer continues to use ordinary Strategic Explainer. Codex manifests do
  not provide a plugin-to-plugin dependency field, so the selected provider is
  installed separately; missing providers remain fail-closed.

Strategic Explainer is a standalone generic communication plugin. It exposes a
fresh stateless API for one real user-facing comment, report, decision or state
explanation, blocker report, final, or explicit editing request. A new default
subagent with `fork_turns="none"` receives one short task, exact scope, and
resolvable read-only source anchors without inherited process context, caller
analysis, or candidate text. The caller sees only this opaque protocol, accepts
ready text and a separate source basis or an operational refusal, and never
reads or applies the provider method. Invalid invocations and factual
corrections use a new instance. The provider reconstructs the reader's actual
question from evidence, preserves material facts and uncertainty, removes
implementation noise and mixed-language jargon, and makes no status, scope,
authority, or mutation decisions.

Strategic Explainer Fast is a separate standalone generic communication plugin.
It carries the same human-facing outcome, truth, language and authority
requirements, but runs in the current agent context without spawning a
subagent. It treats inherited process history as untrusted framing, rebuilds the
message from current read-only sources, performs an in-context comprehension
pass and returns publication text separately from source basis. It explicitly
does not claim ordinary Strategic Explainer's clean-context or independent-reader
guarantees.

Mentioning `$task-manager` alone still authorizes only the requested adapter
operation. Use `$ship-tasks` or an unambiguous natural-language delivery request
that names an existing Task or selected Task Manager Project/Release/current
scope. A delivery verb alone does not route ordinary code, product, repository,
or plugin work into ShipTask. Create-and-deliver requires an explicit request
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
- `plugins/strategic-explainer/.codex-plugin/plugin.json` — standalone
  Strategic Explainer plugin manifest;
- `plugins/strategic-explainer/skills/strategic-explainer/` — router plus
  admission-gated provider contract used only by a fresh subagent;
- `plugins/strategic-explainer-fast/.codex-plugin/plugin.json` — standalone Fast
  plugin manifest;
- `plugins/strategic-explainer-fast/skills/strategic-explainer-fast/` —
  in-context/no-subagent communication runtime;
- `plugins/mind-diary/.codex-plugin/plugin.json` — Mind Diary plugin manifest;
- `plugins/mind-diary/.mcp.json` — direct Mind Diary MCP and OAuth resource;
- `plugins/mind-diary/bin/` and `local-companion/` — bundled macOS exact-file
  launcher, binaries, source and tests;
- `plugins/mind-diary/assets/` — Mind Diary brand assets;
- `plugins/mind-diary/skills/mind-diary/` — bounded content workflow guidance.

Each future plugin gets its own `plugins/<plugin-name>/` directory and one
catalog entry. Task Manager remains independently installable as
`task-manager@srez-marketplace`; Ship Tasks and Strategic Explainer are
independently installable as `ship-tasks@srez-marketplace`,
`strategic-explainer@srez-marketplace`, and
`strategic-explainer-fast@srez-marketplace`.

Validate plugin changes before publishing:

```bash
python3 -m unittest discover -s tests -p 'test_*.py' -v
(cd plugins/task-manager/local-companion && go test -race ./...)
(cd plugins/mind-diary/local-companion && go test -race ./...)
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/task-manager
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/ship-tasks
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/task-composer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/strategic-explainer/skills/strategic-explainer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/strategic-explainer-fast/skills/strategic-explainer-fast
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/mind-diary/skills/mind-diary
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/task-manager
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/ship-tasks
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/strategic-explainer
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/strategic-explainer-fast
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/mind-diary
jq empty plugins/task-manager/.codex-plugin/plugin.json \
  plugins/task-manager/.mcp.json \
  plugins/ship-tasks/.codex-plugin/plugin.json \
  plugins/strategic-explainer/.codex-plugin/plugin.json \
  plugins/strategic-explainer-fast/.codex-plugin/plugin.json \
  plugins/mind-diary/.codex-plugin/plugin.json \
  plugins/mind-diary/.mcp.json \
  .agents/plugins/marketplace.json
! rg -n 'Single-task delivery|Work completion reports|without a Goal|deliver one' \
  plugins/task-manager/skills/task-manager
diff -qr /Users/andrey/Projects/Home/ShipTask/ship-tasks \
  plugins/ship-tasks/skills/ship-tasks
diff -qr /Users/andrey/Projects/Home/ShipTask/task-composer \
  plugins/ship-tasks/skills/task-composer
diff -qr /Users/andrey/Projects/Home/ShipTask/strategic-explainer \
  plugins/strategic-explainer/skills/strategic-explainer
diff -qr /Users/andrey/Projects/Home/ShipTask/strategic-explainer-fast \
  plugins/strategic-explainer-fast/skills/strategic-explainer-fast
git diff --check
```
