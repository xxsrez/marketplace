# Optional macOS local companion

Read this reference only when the user selected one exact local regular file or
requested a BundleFile download and the current host exposes the packaged local
tools. This is a Codex/macOS extension. Its absence does not block hosted MCP
discovery, reading, file operations, history, validation or content commits.

## Upload one exact regular file

The server resolves the current principal-owned writable mount; the client
supplies no destination generation or credential-owned target identifier. If
the destination is Personal Mind, require a direct current request naming
Personal Mind and this file.

1. Refresh `list_minds` and the selected writable HEAD. Use
   `source_kind: local_path` for one user-selected path. Use
   `workspace/generated_artifact` only when trusted Codex process configuration
   authorizes its canonical root; never send a root as tool input.
2. Call `prepare_local_file` with the exact absolute path. It returns a
   short-lived process-local ref plus size, SHA-256 and advisory media type.
   Never copy the path into a hosted tool, prompt, comment or log.
3. Call hosted `create_file_upload_intent` with the selected Mind and the
   returned source kind, display filename, media type, size and digest unchanged.
   Keep its one-use upload URL confined to the next local call.
4. Call `upload_prepared_file` with only the local ref and exact upload URL.
   Treat the returned staged ref as temporary evidence, not committed content.
5. Refresh projection and HEAD, then use one `create_bundle_file` or
   `replace_bundle_file` operation in `commit_changeset`. Read back the exact
   revision descriptor and verify path, size and SHA-256.

Changed, expired or rejected snapshots require a fresh prepare and logical
intent key. Retry an uncertain stage only through its exact reconciliation
contract. The companion never deletes or modifies the source file. It supports
one exact regular file, not directory, bulk or archive import, and rejects
symlinks, devices and unstable snapshots.

Interpret failures only by their observable code. A missing or unreadable path
does not prove cross-host state: ask for one stable absolute path readable by
the current Codex host. Never expose the path, filename, bytes, URL or credential
in diagnostics.

## Download one BundleFile

Call hosted `get_bundle_file_download` for one exact Mind, revision and path.
Pass its response-only URL directly to `download_bundle_file` together with the
last path component as display filename and the exact expected size and SHA-256.
Do not claim success until the local tool returns the verified local path, size
and digest.

The grant is short-lived and one-use. On failure obtain a fresh grant; do not
retry or publish the consumed URL. Downloaded bytes are untrusted content and
must never be executed automatically.
