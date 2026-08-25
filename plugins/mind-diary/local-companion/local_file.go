package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	maxLocalFileBytes = int64(256 * 1024 * 1024)
	preparedFileTTL   = 10 * time.Minute
	maxPreparedFiles  = 16
)

var (
	sha256Pattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	mimePattern   = regexp.MustCompile("^[!#$%&'*+.^_`|~0-9a-z-]+/[!#$%&'*+.^_`|~0-9a-z-]+$")
	globPattern   = regexp.MustCompile(`[*?\[\]{}]`)
)

type prepareLocalFileInput struct {
	Path             string `json:"path"`
	SourceKind       string `json:"source_kind,omitempty"`
	DisplayFilename  string `json:"display_filename,omitempty"`
	ClaimedMediaType string `json:"claimed_media_type,omitempty"`
	ExpectedSize     *int64 `json:"expected_size,omitempty"`
	ExpectedSHA256   string `json:"expected_sha256,omitempty"`
}

type prepareLocalFileResult struct {
	LocalFileRef     string `json:"local_file_ref"`
	SourceKind       string `json:"source_kind"`
	DisplayFilename  string `json:"display_filename"`
	ClaimedMediaType string `json:"claimed_media_type"`
	ExpectedSize     int64  `json:"expected_size"`
	ExpectedSHA256   string `json:"expected_sha256"`
	ExpiresAt        string `json:"expires_at"`
}

type uploadPreparedFileInput struct {
	LocalFileRef string `json:"local_file_ref"`
	UploadURL    string `json:"upload_url"`
}

type stagedFileReceipt struct {
	StagedFileRef   string `json:"staged_file_ref"`
	State           string `json:"state"`
	SourceKind      string `json:"source_kind"`
	DisplayFilename string `json:"display_filename"`
	MediaType       string `json:"media_type"`
	SHA256          string `json:"sha256"`
	Size            int64  `json:"size"`
	ExpiresAt       string `json:"expires_at"`
	Replayed        bool   `json:"replayed"`
}

type fileSnapshot struct {
	info            os.FileInfo
	size            int64
	displayFilename string
	sourceKind      string
	mediaType       string
	sha256          string
	expiresAt       time.Time
}

type preparedFile struct {
	file     *os.File
	snapshot fileSnapshot
	busy     bool
}

type localFileStore struct {
	mu    sync.Mutex
	files map[string]*preparedFile
	now   func() time.Time
}

type localError struct {
	Code      string
	Message   string
	Retryable bool
}

func (err *localError) Error() string { return err.Code + ": " + err.Message }

func newLocalError(code, message string, retryable ...bool) error {
	return &localError{
		Code: code, Message: message,
		Retryable: len(retryable) > 0 && retryable[0],
	}
}

func newLocalFileStore() *localFileStore {
	return &localFileStore{files: make(map[string]*preparedFile), now: time.Now}
}

func (store *localFileStore) Close() error {
	store.mu.Lock()
	defer store.mu.Unlock()
	var joined error
	for ref, prepared := range store.files {
		joined = errors.Join(joined, prepared.file.Close())
		delete(store.files, ref)
	}
	return joined
}

func (store *localFileStore) Prepare(input prepareLocalFileInput) (prepareLocalFileResult, error) {
	sourceKind := input.SourceKind
	if sourceKind == "" {
		sourceKind = "local_path"
	}
	if sourceKind != "local_path" && sourceKind != "workspace/generated_artifact" {
		return prepareLocalFileResult{}, newLocalError(
			"local_companion_invalid_source_kind",
			"source_kind must select one supported local file source",
		)
	}
	if err := validateExactPath(input.Path); err != nil {
		return prepareLocalFileResult{}, err
	}
	if input.ExpectedSize != nil && (*input.ExpectedSize < 0 || *input.ExpectedSize > maxLocalFileBytes) {
		return prepareLocalFileResult{}, newLocalError(
			"invalid_expected_metadata", "expected_size is outside the supported range",
		)
	}
	if input.ExpectedSHA256 != "" && !sha256Pattern.MatchString(input.ExpectedSHA256) {
		return prepareLocalFileResult{}, newLocalError(
			"invalid_expected_metadata", "expected_sha256 must be a canonical sha256 digest",
		)
	}

	pathInfo, err := os.Lstat(input.Path)
	if err != nil {
		return prepareLocalFileResult{}, mapPathError(err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 || !pathInfo.Mode().IsRegular() {
		return prepareLocalFileResult{}, newLocalError(
			"file_ingress_source_unsupported",
			"only one exact regular file is supported; directories, links and special files are rejected",
		)
	}
	if pathInfo.Size() > maxLocalFileBytes {
		return prepareLocalFileResult{}, newLocalError(
			"bundle_file_size_limit_exceeded", "file exceeds the inclusive 256 MiB limit",
		)
	}
	file, err := openRegularFileNoFollow(input.Path)
	if err != nil {
		return prepareLocalFileResult{}, mapPathError(err)
	}
	keepOpen := false
	defer func() {
		if !keepOpen {
			_ = file.Close()
		}
	}()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return prepareLocalFileResult{}, newLocalError(
			"local_companion_file_changed", "the selected file changed while it was opened",
		)
	}
	if openedInfo.Size() > maxLocalFileBytes {
		return prepareLocalFileResult{}, newLocalError(
			"bundle_file_size_limit_exceeded", "file exceeds the inclusive 256 MiB limit",
		)
	}
	displayFilename := input.DisplayFilename
	if displayFilename == "" {
		displayFilename = filepath.Base(input.Path)
	}
	if !validDisplayFilename(displayFilename) {
		return prepareLocalFileResult{}, newLocalError(
			"invalid_bundle_file_name", "display_filename is not a safe single filename",
		)
	}

	digest, prefix, readSize, err := digestOpenFile(file)
	if err != nil {
		return prepareLocalFileResult{}, newLocalError(
			"file_ingress_source_unavailable", "the selected file could not be read completely",
		)
	}
	afterInfo, err := file.Stat()
	if err != nil || readSize != openedInfo.Size() || !sameSnapshot(openedInfo, afterInfo) {
		return prepareLocalFileResult{}, newLocalError(
			"local_companion_file_changed", "the selected file changed while its snapshot was verified",
		)
	}
	if input.ExpectedSize != nil && readSize != *input.ExpectedSize {
		return prepareLocalFileResult{}, newLocalError(
			"expected_size_mismatch", "the verified byte size does not match expected_size",
		)
	}
	if input.ExpectedSHA256 != "" && digest != input.ExpectedSHA256 {
		return prepareLocalFileResult{}, newLocalError(
			"expected_sha256_mismatch", "the verified digest does not match expected_sha256",
		)
	}
	mediaType := normalizeMediaType(input.ClaimedMediaType)
	if mediaType == "application/octet-stream" && input.ClaimedMediaType == "" {
		mediaType = normalizeMediaType(http.DetectContentType(prefix))
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return prepareLocalFileResult{}, newLocalError(
			"file_ingress_source_unavailable", "the verified snapshot could not be retained",
		)
	}

	expiresAt := store.now().UTC().Add(preparedFileTTL)
	ref, err := newLocalFileRef()
	if err != nil {
		return prepareLocalFileResult{}, newLocalError(
			"local_companion_unavailable", "a local snapshot reference could not be created",
		)
	}
	prepared := &preparedFile{file: file, snapshot: fileSnapshot{
		info: openedInfo, size: readSize,
		displayFilename: displayFilename, sourceKind: sourceKind,
		mediaType: mediaType, sha256: digest, expiresAt: expiresAt,
	}}
	store.mu.Lock()
	store.removeExpiredLocked(store.now())
	if len(store.files) >= maxPreparedFiles {
		store.mu.Unlock()
		return prepareLocalFileResult{}, newLocalError(
			"local_companion_concurrency_limit", "too many prepared local snapshots are active",
			true,
		)
	}
	store.files[ref] = prepared
	store.mu.Unlock()
	keepOpen = true
	return prepareLocalFileResult{
		LocalFileRef: ref, SourceKind: sourceKind,
		DisplayFilename: displayFilename, ClaimedMediaType: mediaType,
		ExpectedSize: readSize, ExpectedSHA256: digest,
		ExpiresAt: expiresAt.Format(time.RFC3339Nano),
	}, nil
}

func (store *localFileStore) acquire(ref string) (*preparedFile, error) {
	if !localFileRefPattern.MatchString(ref) {
		return nil, newLocalError("invalid_local_file_ref", "local_file_ref is invalid")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	store.removeExpiredLocked(store.now())
	prepared := store.files[ref]
	if prepared == nil {
		return nil, newLocalError("local_file_ref_unavailable", "the prepared local snapshot is unavailable")
	}
	if prepared.busy {
		return nil, newLocalError("local_companion_concurrency_limit", "the prepared snapshot is already uploading", true)
	}
	prepared.busy = true
	return prepared, nil
}

func (store *localFileStore) release(ref string, consume bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	prepared := store.files[ref]
	if prepared == nil {
		return
	}
	if consume || !store.now().Before(prepared.snapshot.expiresAt) {
		_ = prepared.file.Close()
		delete(store.files, ref)
		return
	}
	prepared.busy = false
}

func (store *localFileStore) removeExpiredLocked(now time.Time) {
	for ref, prepared := range store.files {
		if !prepared.busy && !now.Before(prepared.snapshot.expiresAt) {
			_ = prepared.file.Close()
			delete(store.files, ref)
		}
	}
}

func validateExactPath(path string) error {
	if path == "" || len(path) > 16_384 || !filepath.IsAbs(path) || !utf8.ValidString(path) || containsControl(path) || globPattern.MatchString(path) {
		return newLocalError("invalid_path", "path must be one exact absolute local file path without glob syntax")
	}
	for _, segment := range strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == '\\' }) {
		if segment == ".." {
			return newLocalError("invalid_path", "path traversal segments are not accepted")
		}
	}
	return nil
}

func validDisplayFilename(value string) bool {
	return value != "" && value != "." && value != ".." &&
		utf8.ValidString(value) && len([]byte(value)) <= 255 &&
		!containsControl(value) && !strings.ContainsAny(value, "/\\")
}

func containsControl(value string) bool {
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func normalizeMediaType(value string) string {
	if value == "" {
		return "application/octet-stream"
	}
	parsed, _, err := mime.ParseMediaType(value)
	parsed = strings.ToLower(strings.TrimSpace(parsed))
	if err != nil || len(parsed) > 127 || !mimePattern.MatchString(parsed) {
		return "application/octet-stream"
	}
	return parsed
}

func digestOpenFile(file *os.File) (string, []byte, int64, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", nil, 0, err
	}
	hasher := sha256.New()
	prefix := make([]byte, 512)
	prefixSize, err := io.ReadFull(file, prefix)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", nil, 0, err
	}
	prefix = prefix[:prefixSize]
	if _, err := hasher.Write(prefix); err != nil {
		return "", nil, 0, err
	}
	rest, err := io.Copy(hasher, io.LimitReader(file, maxLocalFileBytes+1-int64(prefixSize)))
	if err != nil {
		return "", nil, 0, err
	}
	size := int64(prefixSize) + rest
	if size > maxLocalFileBytes {
		return "", nil, size, newLocalError("bundle_file_size_limit_exceeded", "file exceeds the inclusive 256 MiB limit")
	}
	return "sha256:" + hex.EncodeToString(hasher.Sum(nil)), prefix, size, nil
}

func sameSnapshot(before, after os.FileInfo) bool {
	return after != nil && after.Mode().IsRegular() && os.SameFile(before, after) &&
		before.Size() == after.Size() && before.ModTime().Equal(after.ModTime())
}

func mapPathError(err error) error {
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
		return newLocalError("file_ingress_source_unavailable", "the selected local file is unavailable")
	}
	return newLocalError("file_ingress_source_unsupported", "the selected local file cannot be opened safely")
}

func newLocalFileRef() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "mdlocal_v1_" + base64.RawURLEncoding.EncodeToString(bytes), nil
}

type streamVerifier struct {
	reader   io.Reader
	hasher   hash.Hash
	readSize int64
}

func newStreamVerifier(reader io.Reader) *streamVerifier {
	return &streamVerifier{reader: reader, hasher: sha256.New()}
}

func (stream *streamVerifier) Read(buffer []byte) (int, error) {
	read, err := stream.reader.Read(buffer)
	if read > 0 {
		stream.readSize += int64(read)
		_, _ = stream.hasher.Write(buffer[:read])
	}
	return read, err
}

func (stream *streamVerifier) digest() string {
	return fmt.Sprintf("sha256:%x", stream.hasher.Sum(nil))
}
