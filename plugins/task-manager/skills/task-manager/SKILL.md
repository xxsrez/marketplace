---
name: task-manager
description: Use Task Manager through its connected MCP tools as a technical adapter to find, filter, inspect, create, update, relate, and comment on tasks; work with native attachments; and navigate projects and releases. Trigger for Task Manager data access, canonical reference resolution, current versions, OAuth capabilities, pagination, safe writes, relations, hierarchy, labels, attachments, and native comments. Do not independently define or run delivery, Goal, verification, release, report-content, or terminal-status policy; a calling workflow such as ShipTask owns those decisions.
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

### Labels, hierarchy, relations, and moves

- Treat Label, parent Task, relation, Project, Release, and workflow status refs
  as canonical opaque identifiers. Resolve them through reads before writes.
- Use the dedicated Label tools for one assignment change and
  `replace_task_labels` for an atomic full replacement. Do not infer catalog
  administration from assignment intent.
- Pass the current child Task version to `set_task_parent` and
  `create_subtask`. Same-Project and cycle rules are server-enforced.
- Relation create requires one stable idempotency key. Relation update/delete
  uses the current relation version; changing to `duplicate_of` also requires
  the current source Task version. `blocked_by` and `duplicates` are relative
  read presentations, not stored relation types.
- Use `move_task` rather than generic `update_task` for a Project change. Read
  the target Project first and explicitly clear or replace incompatible Release
  or assignee values; the returned identifier is authoritative.

## Native attachments

- Call `list_task_attachments` only after selecting a Task. Use
  `get_task_attachment` for metadata and protected preview/original URLs, and
  `download_task_attachment` when the client needs the protected original as a
  resource link.
- Route the source before upload:
  - a native OpenAI file object can go to `upload_file`; retain its verified
    `fileRef`, then call `attach_file_to_task` after the target Task exists;
  - an exact user-authorized absolute filesystem path can go only to the local
    `upload_local_file` companion. Pass a stable upload key and, when known,
    expected byte size/SHA-256. Verify its returned metadata, then pass only the
    resulting `fileRef` to local `attach_local_file_to_task` with an independent
    bind key;
  - when the target Task already exists and no staged workflow is needed,
    `upload_task_attachment` remains the one-call compatibility wrapper;
  - an existing `fileRef` goes only to `get_file`, `delete_file`, or
    `attach_file_to_task`; it is not a durable body reference;
  - an existing `attachmentRef` goes to Task attachment/comment tools, not back
    through upload;
  - never pass a local path, base64 payload, DOM URL, or arbitrary remote URL to
    hosted MCP. Local paths enter only through `upload_local_file`; otherwise
    report that the source has not entered native upload.
- `upload_local_file` and `attach_local_file_to_task` are separate local stdio
  tools, not hosted Task Manager tools. Upload requires one absolute path and
  still obeys the Codex host sandbox and
  approval policy. It rejects relative paths, globs, directories,
  final-component symlinks and special files; computes a stable bounded
  snapshot and SHA-256 before network
  I/O; and sends only basename/display filename, MIME, idempotency key and bytes.
  It returns no local path. On first use, its browser PKCE consent is separate
  from the remote MCP connection. Initialization and tool discovery perform no
  network, browser or credential-store access; OAuth metadata and tokens remain
  only in process memory, so a new companion process authorizes again. Bind
  calls the canonical Agent REST endpoint and returns `attachmentRef`; hosted
  Codex and hosted MCP are not part of this local workflow.
- The operator-only UAT profile uses an exact loopback transport started only
  for a local-file smoke. If it is absent, report the actionable local-ingress
  requirement; do not fall back to production or hosted MCP. Never request or
  pass a Sites token in tool arguments or chat. The ingress accepts it only on
  its manual stdin, keeps it in memory, and never uses Keychain or
  `/usr/bin/security`.
- Reuse a local upload idempotency key only for identical bytes and metadata. A
  network error is not evidence that upload failed: retry the exact source with
  the same key. Never change the path, display filename or expected metadata
  while reusing that key.
- `upload_file` and `upload_task_attachment` advertise the native `file`
  parameter through `_meta["openai/fileParams"]`. Their `file_id` and temporary
  `download_url` are transport inputs, not storage identity. Reuse an upload
  idempotency key only for identical bytes and metadata; use an independent
  stable key for bind.
- File-first example: `upload_file(file, uploadKey)` → verify
  filename/MIME/byteSize/SHA-256/state → create or resolve Task →
  `attach_file_to_task(taskRef, fileRef, bindKey)` → use the returned
  `attachmentRef` in a comment or description → read back both Task attachment
  metadata and the comment/thread.
- Local file-first example: `upload_local_file(path, uploadKey,
  expectedByteSize, expectedSha256)` → verify the local/server metadata and
  `sourceVerified=true` → create or resolve Task →
  `attach_local_file_to_task(taskRef, fileRef, bindKey)` → use only the returned
  `attachmentRef` in description/comment → read back attachment and thread.
- Attachment refs are opaque. To embed an image or link a file, insert the
  returned `attachment:v1:<ref>` token into the Task description or native
  Comment body through the ordinary versioned write. Do not copy protected
  URLs into durable text.
- Comment composition is upload first, then `add_task_comment`,
  `reply_to_task_comment`, or `edit_task_comment`, followed by comment and
  attachment read-back. Comment reads return bounded refs, not binary data.
- Before `delete_task_attachment`, list the current attachment and pass its
  version. Remove every live description/comment ref first; deletion is soft
  and the server rejects referenced or inaccessible attachments.

File errors deliberately disclose little. Invalid source URL, content, MIME,
size, or filename is an input error. Reusing an upload/bind key for a changed
payload is a conflict. Unknown, foreign, expired, deleted, already-bound, and
post-bind `fileRef` reads are unavailable without confirming which condition
occurred. If bind fails, keep the verified `fileRef` only when it remains
readable through `get_file`; otherwise reread the Task attachments before any
retry so an unknown write outcome cannot create duplicate work.

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
