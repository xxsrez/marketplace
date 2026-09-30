# Task Manager local file companion

This bundled stdio MCP server exposes two staged local tools:

`upload_local_file(path, idempotencyKey, expectedByteSize?, expectedSha256?, displayFilename?)`

`attach_local_file_to_task(taskRef, fileRef, idempotencyKey, displayName?)`

It reads one absolute host-authorized path, rejects final-component symlinks and
non-regular files, snapshots at most 25 MiB through a stable open handle, verifies optional
size/SHA-256 expectations, and uploads the snapshot to Task Manager
`POST /api/agent/v1/files`. Only the basename or explicit display filename,
verified MIME, stable idempotency key and bytes leave the host. The local path is
not sent to Task Manager and is not returned by the tool.

The returned `fileRef` is deliberately unbound. Call
`attach_local_file_to_task` with an independent bind idempotency key to receive
the durable Task-scoped `attachmentRef`. Both operations use the canonical
Agent REST contract; hosted Codex and hosted MCP are not part of the local-file
workflow. The companion never proxies remote MCP discovery or tool calls.

## Authentication and lifecycle

Starting the process, MCP initialization, ping and `tools/list` perform no
network, browser, Keychain or `/usr/bin/security` operation. On macOS, the first actual
upload or bind call opens the system browser for Task Manager Authorization
Code + PKCE consent through a random loopback callback and a dynamically
registered public client. Access and refresh tokens stay only in process memory;
a new companion process authorizes again. A revoked refresh token reconnects in
the same bounded operation. The stdio companion contains no Sites bypass or
private-UAT hosting credential.

On Linux (including headless containers), upload and bind instead use the
existing Task Manager personal-token REST contract. Provision a personal API
token with `api:write` as the **runtime** secret `TASK_MANAGER_LOCAL_TOKEN` for
the MCP process. The package explicitly forwards only this named credential
through `env_vars`; it never reads another Codex connection's OAuth cache.
The token owner must have access to the target Task. The hosted MCP connection
remains separately authenticated and should use the same Task Manager account.

No browser, loopback listener, OAuth registration, token refresh or credential
file is used on Linux. Missing/invalid secrets fail before HTTP; a rejected
token is not retried or replaced by browser authentication. Rotate the secret
and restart the process after expiry/revocation. The companion neither creates
tokens nor writes them to disk, logs or tool results. Do not put tokens in chat,
tool arguments, plugin manifests, install scripts or repository files.

For Codex Cloud, configure secret injection into the actual MCP execution
environment. A secret available only to environment setup is insufficient;
do not copy it into a snapshot to work around that boundary. If the host cannot
inject a runtime credential, use its native file bridge when available; report
the missing runtime-credential capability otherwise. Merely installing the
package does not establish either authentication path.

Release acceptance uses the production origin only after explicit production
approval. A bounded canary may require one first-use browser PKCE consent, but
never stores credentials in Keychain or starts a proxy/daemon. The shipped
binary has no origin or transport override; production is its only data plane.

## Packaging and portability

The plugin ships self-contained Go binaries for macOS and Linux arm64 and x86_64 plus a
POSIX launcher. It does not depend on system Node, npm, Python or a package
manager. Windows is not enabled. The launcher chooses by both OS and architecture.
Linux uses `O_NOFOLLOW`, nonblocking open, device/inode identity and nanosecond
mtime/ctime checks around the bounded read, preserving the regular-file and
change-detection contract. A Linux build does not imply a live Cloud upload was
tested.

Build and test from this directory:

```bash
gofmt -w *.go
go test -race ./...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-darwin-arm64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-darwin-amd64 .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-linux-amd64 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/task-manager-local-linux-arm64 .
```

Cloud acceptance: initialize and list tools without credentials; confirm both
local tools are callable; configure the runtime secret; upload one authorized
file with expected size/SHA-256; bind with an independent key; read back the
attachment and compare metadata. Real production canaries still need explicit
production authority. Local mock tests do not establish Cloud integration.
