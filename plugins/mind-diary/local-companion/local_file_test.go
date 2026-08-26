package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

func writeFixture(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestPrepareLocalFileReturnsOnlyPathFreeVerifiedMetadata(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	path := writeFixture(t, "fixture.epub", []byte("arbitrary opaque fixture"))
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	result, err := store.Prepare(prepareLocalFileInput{
		Path: path, SourceKind: "local_path", ClaimedMediaType: "Application/EPUB+ZIP; version=3",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !localFileRefPattern.MatchString(result.LocalFileRef) {
		t.Fatalf("unexpected local ref: %q", result.LocalFileRef)
	}
	if result.SourceKind != "local_path" || result.DisplayFilename != "fixture.epub" {
		t.Fatalf("unexpected metadata: %#v", result)
	}
	if result.ClaimedMediaType != "application/epub+zip" {
		t.Fatalf("unexpected normalized MIME: %q", result.ClaimedMediaType)
	}
	if result.ExpectedSize != int64(len("arbitrary opaque fixture")) ||
		!sha256Pattern.MatchString(result.ExpectedSHA256) {
		t.Fatalf("unverified metadata: %#v", result)
	}
	serialized, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(serialized), path) || strings.Contains(string(serialized), filepath.Dir(path)) {
		t.Fatalf("local path leaked in result: %s", serialized)
	}
}

func TestPrepareRejectsNonExactOrUnsupportedPaths(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	directory := t.TempDir()
	regular := writeFixture(t, "regular.bin", []byte("fixture"))
	symlink := filepath.Join(t.TempDir(), "link.bin")
	if err := os.Symlink(regular, symlink); err != nil {
		t.Fatal(err)
	}
	for name, path := range map[string]string{
		"relative":  "fixture.bin",
		"glob":      filepath.Join(directory, "*.bin"),
		"traversal": directory + "/../" + filepath.Base(directory) + "/fixture.bin",
		"directory": directory,
		"symlink":   symlink,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := store.Prepare(prepareLocalFileInput{Path: path})
			if err == nil {
				t.Fatal("expected rejection")
			}
			var localErr *localError
			if !errorsAs(err, &localErr) {
				t.Fatalf("unexpected error type: %v", err)
			}
		})
	}
}

func TestReadablePathErrorsExposeObservedCauseAndSafeRemediation(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	assertSafe := func(
		t *testing.T,
		err error,
		code, remediation string,
		forbidden ...string,
	) {
		t.Helper()
		var localErr *localError
		if !errors.As(err, &localErr) {
			t.Fatalf("unexpected error type: %v", err)
		}
		if localErr.Code != code || localErr.Remediation != remediation {
			t.Fatalf("unexpected mapped error: %#v", localErr)
		}
		result := toolError(
			localErr.Code,
			localErr.Message,
			localErr.Retryable,
			localErr.Remediation,
		)
		serialized, marshalErr := json.Marshal(result)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		for _, value := range forbidden {
			if value != "" && strings.Contains(string(serialized), value) {
				t.Fatalf("private source data leaked in diagnostic: %s", serialized)
			}
		}
		details := result["structuredContent"].(map[string]any)["error"].(map[string]any)
		if details["message"] != localErr.Message || details["remediation"] != remediation {
			t.Fatalf("structured remediation missing: %#v", details)
		}
	}

	t.Run("missing path is only unavailable on this host", func(t *testing.T) {
		store := newLocalFileStore()
		t.Cleanup(func() { _ = store.Close() })
		missing := filepath.Join(t.TempDir(), "private-missing-name.bin")
		_, err := store.Prepare(prepareLocalFileInput{Path: missing})
		assertSafe(
			t,
			err,
			"file_ingress_source_unavailable",
			readablePathRemediation,
			missing,
			filepath.Base(missing),
		)
	})

	t.Run("unsupported authority remains distinct", func(t *testing.T) {
		store := newLocalFileStore()
		t.Cleanup(func() { _ = store.Close() })
		path := writeFixture(t, "private-workspace-name.bin", []byte("fixture"))
		_, err := store.Prepare(prepareLocalFileInput{
			Path: path, SourceKind: "workspace/generated_artifact",
		})
		assertSafe(
			t,
			err,
			"file_ingress_source_unsupported",
			readablePathRemediation,
			path,
			filepath.Base(path),
		)
	})

	t.Run("changed snapshot has stable-file remediation", func(t *testing.T) {
		selected := writeFixture(t, "private-selected-name.bin", []byte("selected"))
		replacement := writeFixture(t, "private-replacement-name.bin", []byte("replacement"))
		store := newLocalFileStore()
		t.Cleanup(func() { _ = store.Close() })
		store.hooks.open = func(string) (*os.File, error) { return os.Open(replacement) }
		_, err := store.Prepare(prepareLocalFileInput{Path: selected})
		assertSafe(
			t,
			err,
			"local_companion_file_changed",
			changedFileRemediation,
			selected,
			replacement,
			filepath.Base(selected),
			filepath.Base(replacement),
		)
	})

	t.Run("expired local ref asks for a fresh prepare", func(t *testing.T) {
		store := newLocalFileStore()
		t.Cleanup(func() { _ = store.Close() })
		now := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
		store.now = func() time.Time { return now }
		path := writeFixture(t, "private-expired-name.bin", []byte("fixture"))
		prepared, err := store.Prepare(prepareLocalFileInput{Path: path})
		if err != nil {
			t.Fatal(err)
		}
		now = now.Add(preparedFileTTL + time.Second)
		_, err = store.acquire(prepared.LocalFileRef)
		assertSafe(
			t,
			err,
			"local_companion_ref_expired",
			expiredRefRemediation,
			path,
			filepath.Base(path),
		)
	})
}

func TestPrepareAllowsLiteralGlobCharactersInExistingFilename(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	path := writeFixture(t, "report[final].bin", []byte("fixture"))
	result, err := store.Prepare(prepareLocalFileInput{Path: path})
	if err != nil {
		t.Fatalf("expected exact filename with literal glob characters to pass: %v", err)
	}
	if result.DisplayFilename != "report[final].bin" {
		t.Fatalf("unexpected display filename: %#v", result)
	}
}

func TestPrepareEnforcesExpectedMetadataAndInclusiveLimit(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	path := writeFixture(t, "fixture.bin", []byte("fixture"))
	wrongSize := int64(1)
	_, err := store.Prepare(prepareLocalFileInput{Path: path, ExpectedSize: &wrongSize})
	assertLocalCode(t, err, "bundle_file_size_mismatch")
	_, err = store.Prepare(prepareLocalFileInput{
		Path:           path,
		ExpectedSHA256: "sha256:" + strings.Repeat("0", 64),
	})
	assertLocalCode(t, err, "bundle_file_digest_mismatch")
	oversize := filepath.Join(t.TempDir(), "oversize.bin")
	file, err := os.Create(oversize)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(maxLocalFileBytes + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = store.Prepare(prepareLocalFileInput{Path: oversize})
	assertLocalCode(t, err, "bundle_file_size_limit_exceeded")
}

func TestWorkspaceSourceRequiresTrustedCanonicalRoot(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	workspace := t.TempDir()
	outside := t.TempDir()
	insidePath := filepath.Join(workspace, "inside.bin")
	outsidePath := filepath.Join(outside, "outside.bin")
	if err := os.WriteFile(insidePath, []byte("inside"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outsidePath, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	canonicalRoots, err := canonicalWorkspaceRoots([]string{workspace})
	if err != nil {
		t.Fatal(err)
	}
	store := newLocalFileStoreWithWorkspaceRoots(canonicalRoots)
	t.Cleanup(func() { _ = store.Close() })
	prepared, err := store.Prepare(prepareLocalFileInput{
		Path: insidePath, SourceKind: "workspace/generated_artifact",
	})
	if err != nil {
		t.Fatal(err)
	}
	if prepared.SourceKind != "workspace/generated_artifact" {
		t.Fatalf("wrong provenance: %#v", prepared)
	}
	_, err = store.Prepare(prepareLocalFileInput{
		Path: outsidePath, SourceKind: "workspace/generated_artifact",
	})
	assertLocalCode(t, err, "file_ingress_source_unsupported")
	if _, err := store.Prepare(prepareLocalFileInput{
		Path: outsidePath, SourceKind: "local_path",
	}); err != nil {
		t.Fatalf("local_path should remain an explicit one-file authority: %v", err)
	}
	withoutRoots := newLocalFileStore()
	t.Cleanup(func() { _ = withoutRoots.Close() })
	_, err = withoutRoots.Prepare(prepareLocalFileInput{
		Path: insidePath, SourceKind: "workspace/generated_artifact",
	})
	assertLocalCode(t, err, "file_ingress_source_unsupported")
}

func TestWorkspaceAuthorityRejectsCanonicalSymlinkEscape(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	workspace := t.TempDir()
	outside := t.TempDir()
	outsidePath := filepath.Join(outside, "outside.bin")
	if err := os.WriteFile(outsidePath, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(workspace, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	canonicalRoots, err := canonicalWorkspaceRoots([]string{workspace})
	if err != nil {
		t.Fatal(err)
	}
	store := newLocalFileStoreWithWorkspaceRoots(canonicalRoots)
	t.Cleanup(func() { _ = store.Close() })
	_, err = store.Prepare(prepareLocalFileInput{
		Path:       filepath.Join(link, "outside.bin"),
		SourceKind: "workspace/generated_artifact",
	})
	assertLocalCode(t, err, "file_ingress_source_unsupported")
}

func TestWorkspaceAuthorityBindsCanonicalEvidenceToOpenedDescriptor(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	workspace := t.TempDir()
	outside := t.TempDir()
	for _, directory := range []string{workspace, outside} {
		if err := os.WriteFile(filepath.Join(directory, "artifact.bin"), []byte(directory), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	selectorRoot := t.TempDir()
	selector := filepath.Join(selectorRoot, "selected")
	if err := os.Symlink(outside, selector); err != nil {
		t.Fatal(err)
	}
	canonicalRoots, err := canonicalWorkspaceRoots([]string{workspace})
	if err != nil {
		t.Fatal(err)
	}
	store := newLocalFileStoreWithWorkspaceRoots(canonicalRoots)
	t.Cleanup(func() { _ = store.Close() })
	realEval := store.hooks.evalSymlinks
	store.hooks.evalSymlinks = func(path string) (string, error) {
		if err := os.Remove(selector); err != nil {
			return "", err
		}
		if err := os.Symlink(workspace, selector); err != nil {
			return "", err
		}
		return realEval(path)
	}
	_, err = store.Prepare(prepareLocalFileInput{
		Path:       filepath.Join(selector, "artifact.bin"),
		SourceKind: "workspace/generated_artifact",
	})
	assertLocalCode(t, err, "file_ingress_source_unsupported")
}

func TestCanonicalRuntimeErrorContractParity(t *testing.T) {
	expected := []string{
		"bundle_file_digest_mismatch",
		"bundle_file_size_limit_exceeded",
		"bundle_file_size_mismatch",
		"capacity_accounting_untrusted",
		"capacity_fairness_limit",
		"capacity_hard_limit",
		"capacity_soft_limit",
		"file_ingress_intent_conflict",
		"file_ingress_intent_expired",
		"file_ingress_source_unavailable",
		"file_ingress_source_unsupported",
		"file_ingress_transport_unavailable",
		"invalid_bundle_file_name",
		"invalid_path",
		"invalid_request",
		"invalid_upload_url",
		"local_companion_cancelled",
		"local_companion_concurrency_limit",
		"local_companion_file_changed",
		"local_companion_invalid_source_kind",
		"local_companion_ref_expired",
		"local_companion_ref_in_use",
		"local_companion_ref_not_found",
		"staging_quota_exceeded",
	}
	actual := make([]string, 0, len(canonicalRuntimeErrorCodes))
	for code := range canonicalRuntimeErrorCodes {
		actual = append(actual, code)
	}
	slices.Sort(actual)
	if !slices.Equal(actual, expected) {
		t.Fatalf("Go/runtime error contract diverged:\nactual=%v\nexpected=%v", actual, expected)
	}
	for _, retired := range []string{
		"invalid_expected_metadata", "expected_size_mismatch", "expected_sha256_mismatch",
		"invalid_local_file_ref", "local_file_ref_unavailable", "invalid_arguments", "internal_error",
	} {
		if _, exists := canonicalRuntimeErrorCodes[retired]; exists {
			t.Fatalf("retired adapter-only code remains canonical: %s", retired)
		}
	}
}

func TestWorkspaceRootsComeOnlyFromValidatedProcessConfiguration(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	first := t.TempDir()
	second := t.TempDir()
	t.Setenv(workspaceRootsEnvironment, strings.Join([]string{first, second}, string(filepath.ListSeparator)))
	roots, err := configuredWorkspaceRoots()
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 2 || !withinWorkspaceRoots(filepath.Join(roots[0], "file.bin"), roots) {
		t.Fatalf("trusted roots were not canonicalized: %#v", roots)
	}
	t.Setenv(workspaceRootsEnvironment, "relative-root")
	if _, err := configuredWorkspaceRoots(); err == nil {
		t.Fatal("relative process configuration was accepted")
	}
}

func TestPreparedRefExpiresAndBusyCleanupDoesNotCloseActiveDescriptor(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	now := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	store.now = func() time.Time { return now }
	path := writeFixture(t, "fixture.bin", []byte("fixture"))
	prepared, err := store.Prepare(prepareLocalFileInput{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	active, err := store.acquire(prepared.LocalFileRef)
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(preparedFileTTL + time.Second)
	store.mu.Lock()
	store.removeExpiredLocked(now)
	_, retained := store.files[prepared.LocalFileRef]
	store.mu.Unlock()
	if !retained {
		t.Fatal("expiry cleanup closed a busy descriptor")
	}
	if _, err := active.file.Seek(0, 0); err != nil {
		t.Fatalf("busy descriptor was closed: %v", err)
	}
	store.release(prepared.LocalFileRef, false)
	store.mu.Lock()
	_, retained = store.files[prepared.LocalFileRef]
	store.mu.Unlock()
	if retained {
		t.Fatal("expired descriptor survived release")
	}
}

func TestPreparedRefUsesCanonicalNotFoundExpiredAndInUseCodes(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	now := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	store.now = func() time.Time { return now }
	path := writeFixture(t, "fixture.bin", []byte("fixture"))
	prepared, err := store.Prepare(prepareLocalFileInput{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.acquire(prepared.LocalFileRef); err != nil {
		t.Fatal(err)
	}
	_, err = store.acquire(prepared.LocalFileRef)
	assertLocalCode(t, err, "local_companion_ref_in_use")
	store.release(prepared.LocalFileRef, false)
	now = now.Add(preparedFileTTL + time.Second)
	_, err = store.acquire(prepared.LocalFileRef)
	assertLocalCode(t, err, "local_companion_ref_expired")
	_, err = store.acquire(prepared.LocalFileRef)
	assertLocalCode(t, err, "local_companion_ref_not_found")
	_, err = store.acquire("invalid")
	assertLocalCode(t, err, "local_companion_ref_not_found")
}

func TestProtocolListsExactTwoStepTools(t *testing.T) {
	result, protocolErr := handleMCPRequest(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/list",
	}, nil)
	if protocolErr != nil {
		t.Fatal(protocolErr)
	}
	tools := result.(map[string]any)["tools"].([]any)
	if len(tools) != 2 {
		t.Fatalf("expected two tools, got %d", len(tools))
	}
	prepare := tools[0].(map[string]any)
	upload := tools[1].(map[string]any)
	if prepare["name"] != "prepare_local_file" || upload["name"] != "upload_prepared_file" {
		t.Fatalf("unexpected tool names: %#v %#v", prepare["name"], upload["name"])
	}
	uploadSchema := upload["inputSchema"].(map[string]any)
	properties := uploadSchema["properties"].(map[string]any)
	if len(properties) != 2 || properties["local_file_ref"] == nil || properties["upload_url"] == nil {
		t.Fatalf("upload schema admits unexpected authority: %#v", properties)
	}
	outputSchema := upload["outputSchema"].(map[string]any)
	outputProperties := outputSchema["properties"].(map[string]any)
	for _, forbidden := range []string{"path", "upload_url", "local_file_ref"} {
		if outputProperties[forbidden] != nil {
			t.Fatalf("upload output leaks local/capability authority through %q", forbidden)
		}
	}
}

func errorsAs(err error, target any) bool {
	return errors.As(err, target)
}

func assertLocalCode(t *testing.T, err error, code string) {
	t.Helper()
	var localErr *localError
	if !errorsAs(err, &localErr) || localErr.Code != code {
		t.Fatalf("expected %s, got %v", code, err)
	}
}
