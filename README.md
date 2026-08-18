# Srez Marketplace

This is Andrey's extensible Codex plugin marketplace. It contains independent
Task Manager and Mind Diary plugins under `plugins/`.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/marketplace
```

Then install either plugin:

1. Open **Plugins → Task Manager** or **Plugins → Mind Diary** and click
   **Install**.
2. For Task Manager, start using the plugin and authenticate when Codex first
   connects to its MCP server. Mind Diary requests authentication during
   installation.
3. Start a new task in Codex after installation or authentication so it loads
   the selected plugin's current skills and tools.

No server URL, client ID, secret, or personal API token is required. Task
Manager distributes its production MCP connection directly and authenticates
on first use. Mind Diary uses its registered app connector and authenticates on
install.

The plugins connect to:

- Task Manager: `https://task-manager.xxsrez-work.chatgpt.site/api/mcp`
- Mind Diary: `https://mind-diary.xxsrez-work.chatgpt.site/api/mcp`

Task Manager bundles two skills behind one OAuth-protected MCP connection:

- `task-manager` is the technical adapter for OAuth/MCP discovery, canonical
  references, pagination, safe writes, optimistic concurrency, and native
  comment mechanics;
- `ship-tasks` owns delivery intent, lifecycle, Goals, verification, releases,
  report content, and terminal status policy.

Mentioning `$task-manager` alone still authorizes only the requested adapter
operation. Use `$ship-tasks` or an unambiguous natural-language delivery request
to execute work.

Mind Diary bundles one content skill. It selects one explicit Mind and revision,
uses browse/search/fetch for bounded reads, and preserves immutable history,
HEAD CAS and idempotency on explicit writes. It does not treat ChatGPT/Codex
conversation memory as Mind Diary content.

Repository layout:

- `.agents/plugins/marketplace.json` — the ordered marketplace catalog;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.mcp.json` — direct production MCP and OAuth resource;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance;
- `plugins/task-manager/skills/ship-tasks/` — Task Manager delivery workflow;
- `plugins/mind-diary/.codex-plugin/plugin.json` — Mind Diary plugin manifest;
- `plugins/mind-diary/.app.json` — registered Mind Diary app connector;
- `plugins/mind-diary/.mcp.json` — Mind Diary MCP and OAuth resource;
- `plugins/mind-diary/assets/` — Mind Diary brand assets;
- `plugins/mind-diary/skills/mind-diary/` — bounded content workflow guidance.

Each future plugin gets its own `plugins/<plugin-name>/` directory and one
catalog entry. Task Manager remains independently installable as
`task-manager@srez-marketplace`.

Validate plugin changes before publishing:

```bash
python3 -m unittest discover -s tests -p 'test_*.py' -v
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/task-manager
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/ship-tasks
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/mind-diary/skills/mind-diary
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/task-manager
python3 /Users/andrey/.codex/skills/.system/plugin-creator/scripts/validate_plugin.py \
  plugins/mind-diary
jq empty plugins/task-manager/.codex-plugin/plugin.json \
  plugins/task-manager/.mcp.json \
  plugins/mind-diary/.codex-plugin/plugin.json \
  plugins/mind-diary/.app.json plugins/mind-diary/.mcp.json \
  .agents/plugins/marketplace.json
! rg -n 'Single-task delivery|Work completion reports|without a Goal|deliver one' \
  plugins/task-manager/skills/task-manager
diff -qr /Users/andrey/Projects/Home/ShipTask/ship-tasks \
  plugins/task-manager/skills/ship-tasks
git diff --check
```
