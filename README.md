# Task Manager Codex marketplace

This marketplace intentionally contains exactly one plugin: Task Manager.

Add the marketplace once:

```bash
codex plugin marketplace add xxsrez/task-manager-codex-connector
```

Then:

1. Open **Plugins → Task Manager** in Codex Desktop and click **Install**.
2. Click **Authenticate** and sign in to Task Manager.
3. Start a new task in Codex.

No server URL, client ID, secret, or personal API token is required. The plugin
uses the registered Task Manager app connector and OAuth.

The plugin connects to:

`https://task-manager.xxsrez-work.chatgpt.site/api/mcp`

Repository layout:

- `.agents/plugins/marketplace.json` — the single marketplace entry;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.app.json` — registered Task Manager app connector;
- `plugins/task-manager/.mcp.json` — remote MCP and OAuth resource;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance.
