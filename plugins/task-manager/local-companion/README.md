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
the same bounded operation. The stdio companion contains no Sites bypass or
private-UAT hosting credential.

Release acceptance uses the production origin only after explicit production
approval. A bounded canary may require one first-use browser PKCE consent, but
never stores credentials in Keychain or starts a proxy/daemon.

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
