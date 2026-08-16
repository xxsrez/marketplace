# Srez Marketplace

This is Andrey's extensible Codex plugin marketplace. It currently contains one
plugin, Task Manager; additional independent plugins can be added under
`plugins/` over time.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/marketplace
```

Then:

1. Open **Plugins → Task Manager** in Codex Desktop and click **Install**.
2. Click **Authenticate** and sign in to Task Manager.
3. Start a new task in Codex.

No server URL, client ID, secret, or personal API token is required. The plugin
uses the registered Task Manager app connector and OAuth.

The plugin connects to:

`https://task-manager.xxsrez-work.chatgpt.site/api/mcp`

The bundled `task-manager` skill is a technical adapter: it defines OAuth/MCP
discovery, canonical references, pagination, safe writes, optimistic
concurrency, and native comment mechanics. A calling workflow owns software
delivery lifecycle, Goals, verification, releases, report content, and terminal
status policy; mentioning `$task-manager` alone does not grant those effects.

Repository layout:

- `.agents/plugins/marketplace.json` — the ordered marketplace catalog;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.app.json` — registered Task Manager app connector;
- `plugins/task-manager/.mcp.json` — remote MCP and OAuth resource;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance.

Each future plugin gets its own `plugins/<plugin-name>/` directory and one
catalog entry. Task Manager remains independently installable as
`task-manager@srez-marketplace`.

Validate a Task Manager plugin change before publishing:

```bash
python3 /Users/andrey/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  plugins/task-manager/skills/task-manager
jq empty plugins/task-manager/.codex-plugin/plugin.json \
  plugins/task-manager/.app.json plugins/task-manager/.mcp.json \
  .agents/plugins/marketplace.json
! rg -n 'Single-task delivery|Work completion reports|without a Goal|deliver one' \
  plugins/task-manager/skills/task-manager \
  plugins/task-manager/.codex-plugin/plugin.json
git diff --check
```
