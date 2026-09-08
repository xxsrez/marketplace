package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var downloadPathPattern = regexp.MustCompile(`^/api/bundle-download/mdg_v1_[A-Za-z0-9_-]{43}$`)

type downloadBundleFileInput struct {
	DownloadURL     string `json:"download_url"`
	DisplayFilename string `json:"display_filename"`
	ExpectedSize    *int64 `json:"expected_size"`
	ExpectedSHA256  string `json:"expected_sha256"`
}
type downloadBundleFileResult struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

func (service *localFileServiceImpl) DownloadBundleFile(ctx context.Context, input downloadBundleFileInput) (downloadBundleFileResult, error) {
	empty := downloadBundleFileResult{}
	u, err := url.Parse(input.DownloadURL)
	if err != nil || len(input.DownloadURL) > 4096 || u.Scheme+"://"+u.Host != service.publicOrigin || u.User != nil || u.RawQuery != "" || u.ForceQuery || strings.ContainsAny(input.DownloadURL, "?#") || u.Fragment != "" || u.RawPath != "" || !downloadPathPattern.MatchString(u.Path) {
		return empty, newLocalError("invalid_download_url", "use the exact download URL from get_bundle_file_download")
	}
	if input.ExpectedSize == nil || *input.ExpectedSize < 0 || !sha256Pattern.MatchString(input.ExpectedSHA256) {
		return empty, newLocalError("invalid_request", "expected size and SHA-256 from the hosted descriptor are required")
	}
	if *input.ExpectedSize > maxLocalFileBytes {
		return empty, newLocalError("bundle_file_size_limit_exceeded", "file exceeds the local download limit")
	}
	if !validDisplayFilename(input.DisplayFilename) {
		return empty, newLocalError("invalid_bundle_file_name", "display_filename must be one safe filename")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.DownloadURL, nil)
	if err != nil {
		return empty, newLocalError("invalid_download_url", "download URL is invalid")
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "MindDiary-Local/"+buildVersion)
	response, err := service.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return empty, newLocalError("local_companion_cancelled", "download was cancelled; obtain a fresh hosted grant")
		}
		return empty, newLocalError("file_download_unavailable", "download failed; obtain a fresh grant from get_bundle_file_download")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return empty, newLocalError("file_download_unavailable", "download was rejected; obtain a fresh grant from get_bundle_file_download")
	}
	if response.ContentLength >= 0 && response.ContentLength != *input.ExpectedSize {
		return empty, newLocalError("bundle_file_size_mismatch", "download size differs from the hosted descriptor")
	}
	dir, err := os.MkdirTemp("", "mind-diary-download-")
	if err != nil {
		return empty, newLocalError("local_download_failed", "could not create a private download directory")
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(dir)
		}
	}()
	path := filepath.Join(dir, input.DisplayFilename)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return empty, newLocalError("local_download_failed", "could not create the download file")
	}
	defer file.Close()
	digest := sha256.New()
	size, err := io.Copy(io.MultiWriter(file, digest), io.LimitReader(response.Body, *input.ExpectedSize+1))
	if ctx.Err() != nil {
		return empty, newLocalError("local_companion_cancelled", "download was cancelled; obtain a fresh hosted grant")
	}
	if err != nil {
		return empty, newLocalError("file_download_unavailable", "download was interrupted; obtain a fresh hosted grant")
	}
	if size != *input.ExpectedSize {
		return empty, newLocalError("bundle_file_size_mismatch", "download size differs from the hosted descriptor")
	}
	hash := fmt.Sprintf("sha256:%x", digest.Sum(nil))
	if hash != input.ExpectedSHA256 {
		return empty, newLocalError("bundle_file_digest_mismatch", "download SHA-256 differs from the hosted descriptor")
	}
	if err := file.Sync(); err != nil {
		return empty, newLocalError("local_download_failed", "could not persist the verified file")
	}
	if err := file.Close(); err != nil {
		return empty, newLocalError("local_download_failed", "could not close the verified file")
	}
	success = true
	return downloadBundleFileResult{Path: path, Size: size, SHA256: hash}, nil
}

func downloadBundleFileTool() map[string]any {
	return map[string]any{
		"name": "download_bundle_file", "title": "Download and verify one Mind Diary file",
		"description": "Download one exact hosted get_bundle_file_download capability into a new private local directory. Requires the same descriptor's size and SHA-256; returns a path only after byte verification. No credentials, redirects, arbitrary origins or overwrite. Obtain a fresh hosted grant after failure.",
		"inputSchema": map[string]any{"type": "object", "additionalProperties": false,
			"required": []string{"download_url", "display_filename", "expected_size", "expected_sha256"},
			"properties": map[string]any{
				"download_url":     map[string]any{"type": "string", "format": "uri", "maxLength": 4096},
				"display_filename": map[string]any{"type": "string", "minLength": 1, "maxLength": 255},
				"expected_size":    map[string]any{"type": "integer", "minimum": 0, "maximum": maxLocalFileBytes},
				"expected_sha256":  map[string]any{"type": "string", "pattern": "^sha256:[0-9a-f]{64}$"},
			}},
		"outputSchema": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"path", "size", "sha256"}, "properties": map[string]any{
			"path": map[string]any{"type": "string", "minLength": 1}, "size": map[string]any{"type": "integer", "minimum": 0, "maximum": maxLocalFileBytes}, "sha256": map[string]any{"type": "string", "pattern": "^sha256:[0-9a-f]{64}$"},
		}},
		"annotations": map[string]any{"readOnlyHint": false, "destructiveHint": false, "idempotentHint": false, "openWorldHint": true},
	}
}
