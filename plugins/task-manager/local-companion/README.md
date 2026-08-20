# Task Manager local file companion

This bundled stdio MCP server normally exposes one tool:

`upload_local_file(path, idempotencyKey, expectedByteSize?, expectedSha256?, displayFilename?)`

It reads one absolute host-authorized path, rejects final-component symlinks and
non-regular files, snapshots at most 25 MiB through a stable open handle, verifies optional
size/SHA-256 expectations, and uploads the snapshot to Task Manager
`POST /api/agent/v1/files`. Only the basename or explicit display filename,
verified MIME, stable idempotency key and bytes leave the Mac. The local path is
not sent to Task Manager and is not returned by the tool.

The returned `fileRef` is deliberately unbound. Use the remote
`attach_file_to_task` tool with a separate bind idempotency key, then use its
`attachmentRef` in Task descriptions or native comments.

For the explicitly installed `task-manager-uat` validation profile, the same
binary runs with `--private-uat-bridge`. That mode is fail-closed to the exact
`task-manager-uat` origin and proxies the deployed remote MCP through the local
stdio connection, preserving remote tool schemas and
`_meta["openai/fileParams"]` while adding `upload_local_file`. It refuses the
production origin.

## Authentication and lifecycle

First use opens the system browser for Task Manager Authorization Code + PKCE
consent through a random loopback callback and a dynamically registered public
client. The access token remains in process memory. Client metadata and the
rotating refresh token are stored as one macOS Keychain generic-password item;
they are never accepted as tool arguments or plugin configuration. A revoked or
expired refresh token causes a fresh browser reconnect. Reinstall/update leaves
the Keychain item available for the same origin; uninstall removes the plugin
files and no background process remains. The user can revoke the OAuth grant in
Task Manager Settings, and a later use reconnects.

Private UAT has a second, operator-only credential at the outer Sites boundary.
It is read only from the macOS Keychain service
`com.xxsrez.task-manager.uat.sites-bypass`, account
`https://task-manager-uat.xxsrez-work.chatgpt.site`, and is injected only into
requests to that exact HTTPS origin. It is never accepted as a flag, tool
argument, plugin field, environment variable or transcript value. This Sites
credential only passes the private hosting gate; ordinary Task Manager OAuth
still establishes the user, scopes and ACL context.

## Packaging and portability

The plugin ships self-contained Go binaries for macOS arm64 and x86_64 plus a
POSIX launcher. It does not depend on system Node, npm, Python or a package
manager. Windows and Linux are not enabled in this version; the source includes
an explicit unsupported-platform boundary instead of silently weakening file or
credential handling.

Build and test from this directory:

```bash
gofmt -w *.go
go test -race ./...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-darwin-arm64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-darwin-amd64 .
```
