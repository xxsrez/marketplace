# Task Manager Codex marketplace

This marketplace intentionally contains exactly one plugin: Task Manager.

Add the repository as a Codex marketplace, install **Task Manager**, then use
the native OAuth **Connect** action. No personal API token is required.

The plugin connects to:

`https://task-manager.xxsrez-work.chatgpt.site/api/mcp`

Repository layout:

- `.agents/plugins/marketplace.json` — the single marketplace entry;
- `plugins/task-manager/.codex-plugin/plugin.json` — plugin manifest;
- `plugins/task-manager/.mcp.json` — remote MCP and OAuth resource;
- `plugins/task-manager/skills/task-manager/` — agent workflow guidance.
