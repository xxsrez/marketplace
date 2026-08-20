package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const maxLocalFileBytes int64 = 25 * 1024 * 1024

var sha256Pattern = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

type localFileInput struct {
	Path             string `json:"path"`
	IdempotencyKey   string `json:"idempotencyKey"`
	ExpectedByteSize *int64 `json:"expectedByteSize,omitempty"`
	ExpectedSHA256   string `json:"expectedSha256,omitempty"`
	DisplayFilename  string `json:"displayFilename,omitempty"`
}

type preparedLocalFile struct {
	Bytes          []byte
	Filename       string
	MediaType      string
	ByteSize       int64
	ChecksumSHA256 string
}

type fileReadHooks struct {
	Open      func(string) (*os.File, error)
	AfterRead func() error
}

type localError struct {
	Code    string
	Message string
}

func (err *localError) Error() string { return err.Code + ": " + err.Message }

func newLocalError(code, message string) error {
	return &localError{Code: code, Message: message}
}

func prepareLocalFile(input localFileInput, hooks fileReadHooks) (preparedLocalFile, error) {
	if !filepath.IsAbs(input.Path) {
		return preparedLocalFile{}, newLocalError(
			"invalid_local_path", "path must be an absolute path to one authorized file",
		)
	}
	if err := validateIdempotencyKey(input.IdempotencyKey); err != nil {
		return preparedLocalFile{}, err
	}
	if input.ExpectedByteSize != nil && *input.ExpectedByteSize < 0 {
		return preparedLocalFile{}, newLocalError(
			"invalid_expected_metadata", "expectedByteSize must be non-negative",
		)
	}
	if input.ExpectedSHA256 != "" && !sha256Pattern.MatchString(input.ExpectedSHA256) {
		return preparedLocalFile{}, newLocalError(
			"invalid_expected_metadata", "expectedSha256 must be 64 hexadecimal characters",
		)
	}

	cleanPath := filepath.Clean(input.Path)
	pathInfo, err := os.Lstat(cleanPath)
	if err != nil {
		return preparedLocalFile{}, mapLocalPathError(err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 || !pathInfo.Mode().IsRegular() {
		return preparedLocalFile{}, newLocalError(
			"unsupported_local_file", "only a regular file is accepted; links and special files are rejected",
		)
	}
	if pathInfo.Size() > maxLocalFileBytes {
		return preparedLocalFile{}, newLocalError(
			"local_file_too_large", "file exceeds the 25 MiB local upload limit",
		)
	}
	if input.ExpectedByteSize != nil && pathInfo.Size() != *input.ExpectedByteSize {
		return preparedLocalFile{}, newLocalError(
			"local_metadata_mismatch", "local byte size does not match expectedByteSize",
		)
	}

	open := hooks.Open
	if open == nil {
		open = openLocalFileNoFollow
	}
	file, err := open(cleanPath)
	if err != nil {
		return preparedLocalFile{}, mapLocalPathError(err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return preparedLocalFile{}, newLocalError(
			"local_path_unavailable", "authorized file metadata could not be read",
		)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return preparedLocalFile{}, newLocalError(
			"local_file_changed", "authorized file changed while it was being opened",
		)
	}
	before, err := localFileIdentity(openedInfo)
	if err != nil {
		return preparedLocalFile{}, err
	}

	bytes, err := io.ReadAll(io.LimitReader(file, maxLocalFileBytes+1))
	if err != nil {
		return preparedLocalFile{}, newLocalError(
			"local_read_failed", "authorized file could not be read completely",
		)
	}
	if int64(len(bytes)) > maxLocalFileBytes {
		return preparedLocalFile{}, newLocalError(
			"local_file_too_large", "file exceeds the 25 MiB local upload limit",
		)
	}
	if hooks.AfterRead != nil {
		if err := hooks.AfterRead(); err != nil {
			return preparedLocalFile{}, newLocalError(
				"local_file_changed", "authorized file changed during the stable read",
			)
		}
	}
	afterInfo, err := file.Stat()
	if err != nil {
		return preparedLocalFile{}, newLocalError(
			"local_file_changed", "authorized file changed during the stable read",
		)
	}
	after, err := localFileIdentity(afterInfo)
	if err != nil || before != after || int64(len(bytes)) != before.Size {
		return preparedLocalFile{}, newLocalError(
			"local_file_changed", "authorized file changed during the stable read",
		)
	}
	if len(bytes) == 0 {
		return preparedLocalFile{}, newLocalError("empty_local_file", "empty files are not accepted")
	}

	digest := sha256.Sum256(bytes)
	checksum := hex.EncodeToString(digest[:])
	if input.ExpectedByteSize != nil && int64(len(bytes)) != *input.ExpectedByteSize {
		return preparedLocalFile{}, newLocalError(
			"local_metadata_mismatch", "local byte size does not match expectedByteSize",
		)
	}
	if input.ExpectedSHA256 != "" && !strings.EqualFold(input.ExpectedSHA256, checksum) {
		return preparedLocalFile{}, newLocalError(
			"local_metadata_mismatch", "local SHA-256 does not match expectedSha256",
		)
	}

	filename := input.DisplayFilename
	if filename == "" {
		filename = filepath.Base(cleanPath)
	}
	filename, err = validateDisplayFilename(filename)
	if err != nil {
		return preparedLocalFile{}, err
	}
	mediaType, err := detectSafeMediaType(bytes)
	if err != nil {
		return preparedLocalFile{}, err
	}
	return preparedLocalFile{
		Bytes: bytes, Filename: filename, MediaType: mediaType,
		ByteSize: int64(len(bytes)), ChecksumSHA256: checksum,
	}, nil
}

func validateIdempotencyKey(value string) error {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 200 {
		return newLocalError(
			"invalid_idempotency_key", "idempotencyKey must contain 1 to 200 safe characters",
		)
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return newLocalError(
				"invalid_idempotency_key", "idempotencyKey must contain 1 to 200 safe characters",
			)
		}
	}
	return nil
}

func validateDisplayFilename(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "." || value == ".." || !utf8.ValidString(value) ||
		strings.ContainsAny(value, `/\\`) || len([]byte(value)) > 255 {
		return "", newLocalError(
			"invalid_display_filename", "displayFilename must be a single safe filename",
		)
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return "", newLocalError(
				"invalid_display_filename", "displayFilename must be a single safe filename",
			)
		}
	}
	return value, nil
}

func detectSafeMediaType(bytes []byte) (string, error) {
	sampleLength := min(len(bytes), 4096)
	prefix := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(string(bytes[:min(len(bytes), 1024)]), "\ufeff")))
	if strings.HasPrefix(prefix, "<!doctype html") || strings.HasPrefix(prefix, "<html") ||
		strings.HasPrefix(prefix, "<svg") ||
		(strings.HasPrefix(prefix, "<?xml") && (strings.Contains(prefix, "<svg") || strings.Contains(prefix, "<html"))) {
		return "", newLocalError("unsafe_local_content", "active HTML and SVG content is not accepted")
	}
	mediaType := http.DetectContentType(bytes[:sampleLength])
	mediaType = strings.Split(mediaType, ";")[0]
	switch mediaType {
	case "application/pdf", "image/png", "image/jpeg", "image/gif", "text/plain", "application/zip":
		return mediaType, nil
	default:
		return "application/octet-stream", nil
	}
}

func mapLocalPathError(err error) error {
	if errors.Is(err, os.ErrPermission) {
		return newLocalError("local_path_denied", "host access to the authorized path was denied")
	}
	return newLocalError("local_path_unavailable", "authorized path is not an available regular file")
}
