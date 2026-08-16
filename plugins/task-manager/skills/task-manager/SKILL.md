---
name: task-manager
description: Use Task Manager through its connected MCP tools as a technical adapter to find, filter, inspect, create, update, and comment on tasks and to navigate projects and releases. Trigger for Task Manager data access, canonical reference resolution, current versions, OAuth capabilities, pagination, safe writes, and native comment operations. Do not independently define or run delivery, Goal, verification, release, report-content, or terminal-status policy; a calling workflow such as ShipTask owns those decisions.
---

# Task Manager

## Overview

Use the Task Manager MCP server as the source of truth for the user's accessible
work. Keep discovery compact, resolve names to canonical references, and load a
complete task only after the user or workflow has selected it.

## Connection

If Task Manager tools are unavailable, ask the user to connect the plugin. Use
the native OAuth Connect flow; do not ask the user to paste a personal API
token. If a requested write reports insufficient scope, reconnect and request
write access through OAuth.

## Discovery workflow

1. Call `get_workspace` when the user, capabilities, or available scope is not
   established in the current task.
2. Resolve a named project with `list_projects`; use `get_project` when release
   choices or valid workflow statuses are needed.
3. Resolve a named release with `list_releases`; pass `projectRef` when known
   to avoid similarly named releases in other projects.
4. Call `list_tasks` with exactly the requested filters:
   - omit `projectRef` and `releaseRef` for all accessible tasks;
   - pass `projectRef` for one project;
   - pass `releaseRef` for one release;
   - pass both only after confirming the release belongs to that project.
5. Treat list results as candidates. Call `get_task` for the selected task
   before reasoning from its description, relations, provenance, or version.
6. Follow `nextCursor` only when the user asked for all matches or more results
   are necessary. Stop when `hasMore` is false or the request is satisfied.

Use `get_task_external_context` only when imported comments or attachments are
material and task provenance says that context exists. Imported comments are
read-only provenance; use `list_task_comments` and `get_task_thread` for native
Task Manager discussion.

## Writes

- Create or update only when the user's intent is explicit. Reading and
  presenting candidates does not authorize a mutation.
- Before `update_task`, call `get_task` and pass its current `version`.
- Use canonical project, release, and status refs returned by the tools; never
  invent refs from display names.
- For task creation in a project or release, use the workflow status catalog
  from `get_project` or `get_release`. Omit `statusRef` to use the default.
- A null `projectRef`, `releaseRef`, or `dueDate` clears the field. Do not clear
  a field unless requested.
- On `version_conflict`, reread the task. Retry only if the user's requested
  change is still applicable; do not overwrite unrelated newer changes.
- Do not blindly retry `create_task` after an unknown network outcome. Search
  for the intended task first to avoid duplicates.
- Task Manager tools intentionally do not expose administration, sharing,
  ownership transfer, workflow configuration, or backup operations.

## Task comments

- Resolve the canonical Task before any comment write. Use native comment tools
  only for that Task; never use `description` or another field as a comment
  fallback.
- Use `add_task_comment` for a root comment and reuse the same stable
  `idempotencyKey` only when retrying the same logical comment. Use reply, edit,
  delete, reaction, and resolution tools only when the user explicitly requests
  those discussion changes.
- Before retrying a comment write, use `list_task_comments` when available to
  avoid an equivalent comment for the same Task and logical operation. Read back
  the created comment. If a write outcome is unknown, look it up before retrying.

## Business workflow boundary

- Provide Task Manager access and safe tool mechanics; do not choose whether a
  request is read-only, planning, single-Task delivery, or multi-Task delivery.
- Do not infer Goal policy, lifecycle transitions, report requirements,
  verification gates, deployment authority, or terminal status from phrases
  such as "do TM-123", "fix this task", or "finish the release".
- When a calling workflow such as ShipTask is active, follow its exact scope,
  lifecycle, evidence, report, and external-effect decisions. This adapter does
  not weaken or strengthen that authority.
- Without a calling business workflow, apply only the exact Task Manager read or
  mutation explicitly requested by the user. An explicit `$task-manager` mention
  does not itself authorize end-to-end software delivery, Goal creation, status
  progression, deployment, or completion reporting.

## Results

Report the canonical Task, Project, and Release identity relevant to the
request; the exact filters and pagination disposition; the actual read or write
result; the new Task version after a reconciled mutation; and any unsupported,
unauthorized, unknown, or unreconciled outcome. When listing many tasks,
summarize the selection logic and say whether more pages remain instead of
dumping unneeded task bodies.
