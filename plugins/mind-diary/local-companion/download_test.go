package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadVerifiedBytesAndFailures(t *testing.T) {
	for _, tc := range []struct {
		name, body, code string
		size             int64
		badHash, chunked bool
		status           int
	}{
		{name: "binary", body: "\x00\xff\x01", size: 3}, {name: "empty", size: 0},
		{name: "wrong-hash", body: "abc", size: 3, badHash: true, code: "bundle_file_digest_mismatch"},
		{name: "short", body: "a", size: 3, chunked: true, code: "bundle_file_size_mismatch"},
		{name: "long", body: "abcd", size: 3, chunked: true, code: "bundle_file_size_mismatch"},
		{name: "denied", size: 0, status: 403, code: "file_download_unavailable"},
		{name: "redirect", size: 0, status: 302, code: "file_download_unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("TMPDIR", root)
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "" || r.Header.Get("Cookie") != "" || r.Header.Get("Accept-Encoding") != "identity" {
					t.Error("unexpected request credentials or encoding")
				}
				if tc.status != 0 {
					w.Header().Set("Location", "https://invalid.example/")
					w.WriteHeader(tc.status)
					return
				}
				if tc.chunked {
					w.(http.Flusher).Flush()
				}
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			service, err := newLocalFileService(server.Client(), server.URL)
			if err != nil {
				t.Fatal(err)
			}
			defer service.Close()
			hash := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(tc.body)))
			if tc.badHash {
				hash = "sha256:" + strings.Repeat("0", 64)
			}
			result, err := service.DownloadBundleFile(context.Background(), downloadBundleFileInput{DownloadURL: server.URL + "/api/bundle-download/mdg_v1_" + strings.Repeat("A", 43), DisplayFilename: "fixture.bin", ExpectedSize: &tc.size, ExpectedSHA256: hash})
			if tc.code != "" {
				assertLocalCode(t, err, tc.code)
				entries, _ := os.ReadDir(root)
				if len(entries) != 0 {
					t.Fatal("partial output leaked")
				}
				if result.Path != "" {
					t.Fatal("failed result exposed path")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			body, err := os.ReadFile(result.Path)
			if err != nil || string(body) != tc.body || result.Size != tc.size || result.SHA256 != hash {
				t.Fatal("incorrect verified output")
			}
			info, _ := os.Stat(result.Path)
			if info.Mode().Perm() != 0600 {
				t.Fatal("unsafe file permissions")
			}
			info, _ = os.Stat(filepath.Dir(result.Path))
			if info.Mode().Perm() != 0700 {
				t.Fatal("unsafe directory permissions")
			}
			_ = service.Close()
			if _, err := os.Stat(result.Path); err != nil {
				t.Fatal("successful download removed on close")
			}
		})
	}
}

func TestDownloadRejectsAuthorityAndMissingMetadata(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(500) }))
	defer server.Close()
	service, _ := newLocalFileService(server.Client(), server.URL)
	defer service.Close()
	size := int64(0)
	base := server.URL + "/api/bundle-download/mdg_v1_" + strings.Repeat("A", 43)
	valid := downloadBundleFileInput{DownloadURL: base, DisplayFilename: "file.bin", ExpectedSize: &size, ExpectedSHA256: "sha256:" + strings.Repeat("0", 64)}
	for _, u := range []string{base + "?", base + "?x=y", base + "#fragment", base + "/more", strings.Replace(base, "/api/", "/%61pi/", 1), strings.Replace(base, "https://", "https://user:password@", 1), "https://example.com" + strings.TrimPrefix(base, server.URL)} {
		input := valid
		input.DownloadURL = u
		_, err := service.DownloadBundleFile(context.Background(), input)
		assertLocalCode(t, err, "invalid_download_url")
	}
	for _, name := range []string{"../escape", "/absolute", "", "a/b"} {
		input := valid
		input.DisplayFilename = name
		_, err := service.DownloadBundleFile(context.Background(), input)
		assertLocalCode(t, err, "invalid_bundle_file_name")
	}
	input := valid
	input.ExpectedSize = nil
	_, err := service.DownloadBundleFile(context.Background(), input)
	assertLocalCode(t, err, "invalid_request")
	if calls != 0 {
		t.Fatal("invalid input reached network")
	}
	args, _ := json.Marshal(map[string]any{"name": "download_bundle_file", "arguments": map[string]any{"download_url": base, "unexpected": "secret"}})
	result, rpcErr := handleToolCall(context.Background(), jsonRPCRequest{Params: args}, service)
	if rpcErr != nil || result.(map[string]any)["isError"] != true {
		t.Fatal("unknown argument accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = service.DownloadBundleFile(ctx, valid)
	assertLocalCode(t, err, "local_companion_cancelled")
}
