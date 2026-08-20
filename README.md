# Srez Marketplace

This is Andrey's extensible Codex plugin marketplace. It contains independent
Task Manager, Task Manager UAT, Ship Tasks, and Mind Diary plugins under
`plugins/`.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/marketplace
```

Then install the plugins you need:

1. Open **Plugins → Task Manager**, **Plugins → Ship Tasks**, or
   **Plugins → Mind Diary UAT** and click **Install**.
2. For Task Manager or Mind Diary, start using the plugin and authenticate when
   Codex first connects to its MCP server. Ship Tasks has no connector of its
   own and requires Task Manager to be installed separately.
3. Start a new task in Codex after installation or authentication so it loads
   the selected plugin's current skills and tools.

`task-manager-uat@srez-marketplace` is an operator-only release-validation
profile. Install it explicitly only for synthetic UAT smoke. A bundled local
stdio bridge proxies the deployed private UAT MCP and adds exact-path upload, so
`fileRef` upload and the later `attach_file_to_task` bind cannot cross
environments. The outer Sites bypass credential and the ordinary Task Manager
OAuth refresh credential remain separate Keychain items; neither is stored in
plugin configuration. The bridge refuses the production origin and does not
replace or reconfigure the production `task-manager` plugin.

No server URL, client ID, secret, or personal API token is required. Task
Manager and Mind Diary UAT distribute their MCP connections directly and
authenticate with OAuth on first use. Mind Diary remains a restricted UAT
pilot; this package is not a production or public-directory release.

The plugins connect to:

- Task Manager: `https://task-manager.xxsrez-work.chatgpt.site/api/mcp`
- Task Manager UAT validation: `https://task-manager-uat.xxsrez-work.chatgpt.site/api/mcp`
- Mind Diary: `https://mind-diary.xxsrez-work.chatgpt.site/api/mcp`

Task Manager is an adapter-only plugin with two coordinated MCP components:

- `task-manager` provides OAuth/MCP discovery, canonical
  references, pagination, safe writes, optimistic concurrency, and native
  comment mechanics.
- `task-manager-local` is a bundled macOS stdio companion with one
  `upload_local_file` tool. It reads one exact host-authorized regular file,
  uploads a verified snapshot into the common `fileRef` workflow, stores its
  rotating OAuth refresh token in Keychain, and never sends the full path to
  Task Manager. The remote `attach_file_to_task` tool performs the later bind.

Ship Tasks is a separate plugin:

- `ship-tasks` owns delivery intent, lifecycle, Goals, verification, releases,
  report content, and terminal status policy;
- `strategic-explainer` is a generic sibling skill that turns bounded technical
  context into clear outcome-first User Briefs without making decisions or
  performing mutations;
- it depends on the separately installed Task Manager plugin for MCP tools and
  authentication, but does not bundle or duplicate that connector.

Mentioning `$task-manager` alone still authorizes only the requested adapter
operation. Use `$ship-tasks` or an unambiguous natural-language delivery request
that names an existing Task or selected Task Manager Project/Release/current
scope. A delivery verb alone does not route ordinary code, product, repository,
or plugin work into ShipTask. Create-and-deliver requires an explicit request
to create exactly one Task in Task Manager and immediately start it.

Mind Diary bundles one content skill. It selects one explicit Mind and revision,
uses browse/search/fetch for bounded reads, and preserves immutable history,
HEAD CAS and idempotency on explicit writes. It does not treat ChatGPT/Codex
conversation memory as Mind Diary content.

Repository layout:

- `.agents/plugins/marketplace.json` — the ordered marketplace catalog;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.mcp.json` — direct production MCP and OAuth resource;
- `plugins/task-manager/bin/` and `local-companion/` — bundled macOS local-file
  launcher, binaries, source and tests;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance;
- `plugins/task-manager-uat/.codex-plugin/plugin.json` — explicit UAT-only
  validation manifest;
- `plugins/task-manager-uat/.mcp.json` and `bin/` — private UAT MCP/local-file
  bridge package;
- `plugins/ship-tasks/.codex-plugin/plugin.json` — Ship Tasks plugin manifest;
- `plugins/ship-tasks/skills/ship-tasks/` — Task Manager delivery workflow;
- `plugins/ship-tasks/skills/strategic-explainer/` — generic communication
  skill used directly or in a fresh subagent;
- `plugins/mind-diary/.codex-plugin/plugin.json` — Mind Diary plugin manifest;
- `plugins/mind-diary/.mcp.json` — direct Mind Diary MCP and OAuth resource;
- `plugins/mind-diary/assets/` — Mind Diary brand assets;
- `plugins/mind-diary/skills/mind-diary/` — bounded content workflow guidance.

Each future plugin gets its own `plugins/<plugin-name>/` directory and one
catalog entry. Task Manager remains independently installable as
`task-manager@srez-marketplace`; Ship Tasks is independently installable as
`ship-tasks@srez-marketplace`. The UAT validation profile is independently
installable as `task-manager-uat@srez-marketplace` and is never the normal
Task Manager data plane.

Validate plugin changes before publishing:

```bash
python3 -m unittest discover -s tests -p 'test_*.py' -v
(cd plugins/task-manager/local-companion && go test -race ./...)
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/task-manager
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/ship-tasks
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/ship-tasks/skills/strategic-explainer
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/mind-diary/skills/mind-diary
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/task-manager
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/task-manager-uat
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/ship-tasks
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/mind-diary
jq empty plugins/task-manager/.codex-plugin/plugin.json \
  plugins/task-manager/.mcp.json \
  plugins/task-manager-uat/.codex-plugin/plugin.json \
  plugins/task-manager-uat/.mcp.json \
  plugins/ship-tasks/.codex-plugin/plugin.json \
  plugins/mind-diary/.codex-plugin/plugin.json \
  plugins/mind-diary/.mcp.json \
  .agents/plugins/marketplace.json
! rg -n 'Single-task delivery|Work completion reports|without a Goal|deliver one' \
  plugins/task-manager/skills/task-manager
diff -qr /Users/andrey/Projects/Home/ShipTask/ship-tasks \
  plugins/ship-tasks/skills/ship-tasks
diff -qr /Users/andrey/Projects/Home/ShipTask/strategic-explainer \
  plugins/ship-tasks/skills/strategic-explainer
git diff --check
```
