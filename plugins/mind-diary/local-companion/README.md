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

`source_kind: workspace/generated_artifact` is accepted only when the opened
descriptor is the same file as a canonical path inside one of the trusted
process-configured roots in `MIND_DIARY_WORKSPACE_ROOTS`. The value is a macOS
path list supplied to the Codex process/launcher and forwarded explicitly by
the packaged `.mcp.json`; it is never accepted as a tool argument. Missing,
relative, nonexistent or non-directory roots fail closed for workspace
provenance. `source_kind: local_path` remains an explicit one-file authority.
Canonical path evidence is rebound to the exact opened descriptor before the
root check, so a parent-symlink swap cannot relabel an outside file as a
workspace artifact.

The adapter emits only the shared runtime error vocabulary. Expected metadata
uses `bundle_file_size_mismatch` and `bundle_file_digest_mismatch`; local refs
use `local_companion_ref_not_found`, `local_companion_ref_expired` and
`local_companion_ref_in_use`. Unknown internal adapter failures collapse to
`file_ingress_transport_unavailable`, and messages remain path-free.
Readable-path failures include a bounded, path-free remediation: choose or copy
one regular file into an absolute path readable by the current Codex host and
prepare it again. A changed snapshot asks the caller to wait for a stable file;
an expired or unavailable local ref asks for a fresh prepare. The adapter never
labels an ordinary missing path as proven cross-host or temporary-file expiry.

The process uses no browser, OAuth flow, Keychain, daemon or proxy. The hosted
Mind Diary MCP remains the authority for the principal's current writable mount
and for minting the one-use intent; the local client supplies no destination
generation or credential-owned target identifier. This binary is pinned to the
Mind Diary UAT origin and has no runtime origin override.

Build and test from this directory:

```bash
gofmt -w *.go
go test -race ./...
./build-release.sh 0.1.0+codex.YYYYMMDDhhmmss ../bin
```

`build-release.sh` is the only release build entrypoint. It builds both macOS
architectures with `-trimpath -buildvcs=false`, so a dirty checkout cannot
stamp a stale Git revision into either artifact. The required immutable plugin
version is supplied explicitly outside Go VCS metadata and embedded in MCP
`serverInfo.version`. Marketplace provenance tests rebuild the exact checked-in
source with this script and require byte-for-byte equality with both packaged
binaries.
