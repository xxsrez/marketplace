# Task Manager Codex marketplace

This marketplace intentionally contains exactly one plugin: Task Manager.

Install it in Codex Desktop:

```bash
codex plugin marketplace add xxsrez/task-manager-codex-connector
codex plugin add task-manager@task-manager
```

Then:

1. Restart Codex Desktop and start a new task.
2. Open **Plugins → Task Manager → Connect** (or run
   `codex mcp login task-manager` in Terminal).
3. Sign in to Task Manager with ChatGPT and approve read/write access.

No personal API token is required. Codex registers its OAuth client
automatically on first connection.

The plugin connects to:

`https://task-manager.xxsrez-work.chatgpt.site/api/mcp`

Repository layout:

- `.agents/plugins/marketplace.json` — the single marketplace entry;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.mcp.json` — remote MCP and OAuth resource;
- `plugins/task-manager/assets/` — card icons and screenshot;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance.
