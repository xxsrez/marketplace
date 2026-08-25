# Mind Diary local file companion

This bundled macOS stdio MCP server adds two path-safe tools beside the hosted
Mind Diary content tools:

1. `prepare_local_file` opens one exact host-authorized regular-file path,
   keeps a stable descriptor in this process, stream-hashes at most 256 MiB and
   returns only a ten-minute `local_file_ref` plus the metadata required by
   hosted `create_file_upload_intent`.
2. `upload_prepared_file` accepts that opaque ref and the exact same-origin
   `upload_url`, performs credentialless GET-before-PUT, streams the retained
   descriptor, reconciles unknown outcomes and returns the verified
   `staged_file_ref` receipt.

The local path, bytes, OAuth bearer, cookie and upload capability never appear
in a successful result. The companion rejects directories, globs, traversal,
final-component symlinks, special files, redirects, foreign origins, query or
fragment variants, inline base64 and arbitrary URLs. A successful upload
consumes the local ref. A retryable or unknown transport result keeps the exact
snapshot only until its original expiry; changed, expired or definitively
rejected snapshots are closed.

The process uses no browser, OAuth flow, Keychain, daemon or proxy. The hosted
Mind Diary MCP remains the authority for binding selection and for minting the
one-use intent. This binary is pinned to the Mind Diary UAT origin and has no
runtime origin override.

Build and test from this directory:

```bash
gofmt -w *.go
go test -race ./...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/mind-diary-local-darwin-arm64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath \
  -ldflags='-s -w' -o ../bin/mind-diary-local-darwin-amd64 .
```
