---
name: task-manager
description: Use Task Manager through its connected MCP tools to find, filter, inspect, create, update, and comment on tasks and to navigate projects and releases. Trigger for requests about the user's Task Manager work, current tasks, project or release scope, choosing or changing a task, or executing work tied to an existing task that needs a durable completion, failure, rework, or blocker report.
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
  only for that Task; never use `description` or another field as a report
  fallback.
- Use `add_task_comment` for a root report and reuse the same stable
  `idempotencyKey` only when retrying the same logical report. Use reply, edit,
  delete, reaction, and resolution tools only when the user explicitly requests
  those discussion changes.
- Before a report write, use `list_task_comments` when available to avoid an
  equivalent report for the same Task, state, and exact result. Read back the
  created report. If a write outcome is unknown, look it up before retrying.

## Work completion reports

When the user asks to execute, fix, investigate, deliver, or otherwise do
substantive work for an existing canonical Task, a native report comment is part
of that Task's outcome. Read-only discovery and planning-only requests do not
create reports.

On success, after required checks and external effects pass and before any
requested terminal status transition, publish one `COMPLETED` report that:

- leads with the user-visible outcome;
- explains the main runtime/data flow and important implementation decisions in
  plain language;
- includes exact verification and external-effect evidence;
- states limitations, remaining risk, and a practical review path;
- uses one or two compact text/Markdown diagrams for non-trivial or
  cross-component work, or a concise before/after for a trivial change.

When work stops in material rework, failure, or a blocker, publish `REWORK
REQUIRED`, `FAILED`, or `BLOCKED` with the impact, last safe checkpoint,
evidence, cause and confidence, recovery already performed, remaining risk, and
the exact next step or decision needed. Do not comment on transient red/green
iterations resolved inside the same run.

If the required native comment cannot be published or reconciled, do not claim
the Task is complete and do not move it to a terminal status as part of the
same request. Keep its status truthful, report the comment-delivery blocker,
and continue unrelated safe work when possible.

## Results

Report the task identifier and title, relevant project/release, status, the
change made, and the native report disposition (`published`, `not-required`, or
blocked with the exact reason). When listing many tasks, summarize the selection
logic and say whether more pages remain instead of dumping unneeded task bodies.
