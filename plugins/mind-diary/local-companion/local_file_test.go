package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
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

func TestPrepareEnforcesExpectedMetadataAndInclusiveLimit(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	store := newLocalFileStore()
	t.Cleanup(func() { _ = store.Close() })
	path := writeFixture(t, "fixture.bin", []byte("fixture"))
	wrongSize := int64(1)
	_, err := store.Prepare(prepareLocalFileInput{Path: path, ExpectedSize: &wrongSize})
	assertLocalCode(t, err, "expected_size_mismatch")
	_, err = store.Prepare(prepareLocalFileInput{
		Path:           path,
		ExpectedSHA256: "sha256:" + strings.Repeat("0", 64),
	})
	assertLocalCode(t, err, "expected_sha256_mismatch")
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
