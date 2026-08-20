package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestPrepareLocalFileVerifiesExactRegularFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "evidence.pdf")
	body := []byte("%PDF-1.7\nlocal companion fixture\n")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	wantHash := sha256.Sum256(body)

	prepared, err := prepareLocalFile(localFileInput{
		Path:             path,
		IdempotencyKey:   "tm-290-pdf",
		ExpectedByteSize: int64Ptr(int64(len(body))),
		ExpectedSHA256:   hex.EncodeToString(wantHash[:]),
	}, fileReadHooks{})
	if err != nil {
		t.Fatal(err)
	}
	if prepared.Filename != "evidence.pdf" {
		t.Fatalf("filename = %q", prepared.Filename)
	}
	if prepared.MediaType != "application/pdf" {
		t.Fatalf("media type = %q", prepared.MediaType)
	}
	if prepared.ByteSize != int64(len(body)) {
		t.Fatalf("byte size = %d", prepared.ByteSize)
	}
	if prepared.ChecksumSHA256 != hex.EncodeToString(wantHash[:]) {
		t.Fatalf("checksum = %q", prepared.ChecksumSHA256)
	}
	if string(prepared.Bytes) != string(body) {
		t.Fatal("prepared bytes differ")
	}
}

func TestPrepareLocalFileInsideWorkspace(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".tm-local-inside-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	path, err := filepath.Abs(filepath.Join(dir, "inside.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("inside workspace"), 0o600); err != nil {
		t.Fatal(err)
	}
	prepared, err := prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "inside-workspace",
	}, fileReadHooks{})
	if err != nil {
		t.Fatal(err)
	}
	if prepared.MediaType != "text/plain" || prepared.ByteSize != 16 {
		t.Fatalf("prepared = %#v", prepared)
	}
}

func TestPrepareLocalFileRejectsUnsafeSources(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	regular := filepath.Join(dir, "regular.txt")
	if err := os.WriteFile(regular, []byte("safe"), 0o600); err != nil {
		t.Fatal(err)
	}
	symlink := filepath.Join(dir, "link.txt")
	if err := os.Symlink(regular, symlink); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		code string
	}{
		{name: "relative", path: "regular.txt", code: "invalid_local_path"},
		{name: "directory", path: dir, code: "unsupported_local_file"},
		{name: "symlink", path: symlink, code: "unsupported_local_file"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := prepareLocalFile(localFileInput{
				Path: test.path, IdempotencyKey: "safe-key",
			}, fileReadHooks{})
			assertLocalErrorCode(t, err, test.code)
		})
	}
}

func TestPrepareLocalFileRejectsOversizeAndExpectedMetadataMismatch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	oversized := filepath.Join(dir, "oversized.bin")
	file, err := os.Create(oversized)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(maxLocalFileBytes + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = prepareLocalFile(localFileInput{
		Path: oversized, IdempotencyKey: "oversized",
	}, fileReadHooks{})
	assertLocalErrorCode(t, err, "local_file_too_large")

	path := filepath.Join(dir, "small.txt")
	if err := os.WriteFile(path, []byte("small"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "size-mismatch", ExpectedByteSize: int64Ptr(4),
	}, fileReadHooks{})
	assertLocalErrorCode(t, err, "local_metadata_mismatch")
	_, err = prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "hash-mismatch", ExpectedSHA256: strings.Repeat("0", 64),
	}, fileReadHooks{})
	assertLocalErrorCode(t, err, "local_metadata_mismatch")
}

func TestPrepareLocalFileRejectsChangedDuringRead(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "changing.txt")
	if err := os.WriteFile(path, []byte("before"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "changing",
	}, fileReadHooks{
		AfterRead: func() error {
			return os.WriteFile(path, []byte("after!"), 0o600)
		},
	})
	assertLocalErrorCode(t, err, "local_file_changed")
}

func TestPrepareLocalFileDoesNotExposePathInErrors(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "secret", "missing.pdf")
	_, err := prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "missing",
	}, fileReadHooks{})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), path) || strings.Contains(err.Error(), "secret") {
		t.Fatalf("error disclosed local path: %v", err)
	}
}

func TestPrepareLocalFileMapsPermissionFailureWithoutPath(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "private.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.7\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "permission",
	}, fileReadHooks{
		Open: func(string) (*os.File, error) {
			return nil, os.ErrPermission
		},
	})
	assertLocalErrorCode(t, err, "local_path_denied")
	if strings.Contains(err.Error(), path) {
		t.Fatalf("error disclosed local path: %v", err)
	}
}

func TestAuthorizedExternalFixtureFromEnvironment(t *testing.T) {
	path := os.Getenv("TASK_MANAGER_LOCAL_COMPANION_TEST_FILE")
	if path == "" {
		t.Skip("no explicitly authorized external fixture")
	}
	expectedSize, err := strconv.ParseInt(
		os.Getenv("TASK_MANAGER_LOCAL_COMPANION_TEST_SIZE"), 10, 64,
	)
	if err != nil {
		t.Fatal("TASK_MANAGER_LOCAL_COMPANION_TEST_SIZE must be an integer")
	}
	prepared, err := prepareLocalFile(localFileInput{
		Path: path, IdempotencyKey: "authorized-external-fixture",
		ExpectedByteSize: &expectedSize,
		ExpectedSHA256:   os.Getenv("TASK_MANAGER_LOCAL_COMPANION_TEST_SHA256"),
	}, fileReadHooks{})
	if err != nil {
		t.Fatal(err)
	}
	if prepared.MediaType != os.Getenv("TASK_MANAGER_LOCAL_COMPANION_TEST_MIME") {
		t.Fatalf("media type = %q", prepared.MediaType)
	}
}

func assertLocalErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s error", code)
	}
	var localErr *localError
	if !errors.As(err, &localErr) {
		t.Fatalf("error type = %T: %v", err, err)
	}
	if localErr.Code != code {
		t.Fatalf("error code = %q, want %q", localErr.Code, code)
	}
}

func int64Ptr(value int64) *int64 { return &value }
