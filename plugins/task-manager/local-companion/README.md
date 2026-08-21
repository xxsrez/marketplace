# Task Manager local file companion

This bundled stdio MCP server exposes two staged local tools:

`upload_local_file(path, idempotencyKey, expectedByteSize?, expectedSha256?, displayFilename?)`

`attach_local_file_to_task(taskRef, fileRef, idempotencyKey, displayName?)`

It reads one absolute host-authorized path, rejects final-component symlinks and
non-regular files, snapshots at most 25 MiB through a stable open handle, verifies optional
size/SHA-256 expectations, and uploads the snapshot to Task Manager
`POST /api/agent/v1/files`. Only the basename or explicit display filename,
verified MIME, stable idempotency key and bytes leave the Mac. The local path is
not sent to Task Manager and is not returned by the tool.

The returned `fileRef` is deliberately unbound. Call
`attach_local_file_to_task` with an independent bind idempotency key to receive
the durable Task-scoped `attachmentRef`. Both operations use the canonical
Agent REST contract; hosted Codex and hosted MCP are not part of the local-file
workflow. The companion never proxies remote MCP discovery or tool calls.

## Authentication and lifecycle

Starting the process, MCP initialization, ping and `tools/list` perform no
network, browser, Keychain or `/usr/bin/security` operation. The first actual
upload or bind call opens the system browser for Task Manager Authorization
Code + PKCE consent through a random loopback callback and a dynamically
registered public client. Access and refresh tokens stay only in process memory;
a new companion process authorizes again. A revoked refresh token reconnects in
the same bounded operation. The stdio companion does not contain a private-UAT
hosting credential.

## Owner-only UAT

The UAT plugin points machine requests at `http://127.0.0.1:47821`. Nothing is
started there automatically. Immediately before a UAT local-file smoke, the
operator starts the same binary in bounded ingress mode and supplies the Sites
bypass value once through stdin:

```zsh
read -rs 'TM_UAT_SITES_TOKEN?UAT Sites token: '
printf '\n'
exec {TM_UAT_TOKEN_FD}< <(printf '%s' "$TM_UAT_SITES_TOKEN")
unset TM_UAT_SITES_TOKEN
./bin/task-manager-local-launcher --serve-private-uat-ingress <&$TM_UAT_TOKEN_FD
exec {TM_UAT_TOKEN_FD}<&-
unset TM_UAT_TOKEN_FD
```

The value is not an argument, environment variable, config entry or Keychain
item. The short-lived process-substitution writer feeds an inherited descriptor,
then the parent shell variable is cleared before the foreground ingress starts.
The ingress binds exact `127.0.0.1:47821`, keeps the value only in process
memory, exposes only OAuth plus staged-file/TaskAttachment REST routes, rejects
MCP/UI/arbitrary proxy traffic and disappears when the process stops. Browser
authorization still opens the owner-only UAT Site directly. Ordinary UAT UI use
does not require this ingress.

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
