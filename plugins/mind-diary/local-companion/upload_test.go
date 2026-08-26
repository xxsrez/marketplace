package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

const testCapability = "mdupload_v1_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

type intentFixture struct {
	mu                sync.Mutex
	prepared          prepareLocalFileResult
	body              []byte
	getCount          int
	putCount          int
	staged            bool
	failPutAfterStage bool
	rejectionCode     string
	requestsClean     bool
}

func (fixture *intentFixture) handler(responseWriter http.ResponseWriter, request *http.Request) {
	fixture.mu.Lock()
	defer fixture.mu.Unlock()
	responseWriter.Header().Set("Content-Type", "application/json")
	if request.URL.Path != uploadIntentRoutePrefix+testCapability || request.URL.RawQuery != "" {
		responseWriter.WriteHeader(http.StatusNotFound)
		_, _ = responseWriter.Write([]byte(`{"ok":false,"error":{"code":"file_ingress_source_unavailable","retryable":false}}`))
		return
	}
	if request.Header.Get("Authorization") != "" || request.Header.Get("Cookie") != "" || request.Referer() != "" {
		fixture.requestsClean = false
	}
	switch request.Method {
	case http.MethodGet:
		fixture.getCount++
		if fixture.rejectionCode != "" {
			_ = json.NewEncoder(responseWriter).Encode(map[string]any{
				"ok":   true,
				"data": map[string]any{"status": "rejected", "code": fixture.rejectionCode},
			})
			return
		}
		if fixture.staged {
			fixture.writeStaged(responseWriter)
			return
		}
		_, _ = responseWriter.Write([]byte(`{"ok":true,"data":{"status":"pending","expires_at":"2026-08-26T00:10:00Z"}}`))
	case http.MethodPut:
		fixture.putCount++
		if request.Header.Get("Content-Type") != "application/octet-stream" {
			fixture.requestsClean = false
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			responseWriter.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		fixture.body = body
		fixture.staged = true
		if fixture.failPutAfterStage {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = responseWriter.Write([]byte(`{"ok":false,"error":{"code":"file_ingress_transport_unavailable","retryable":true}}`))
			return
		}
		responseWriter.WriteHeader(http.StatusCreated)
		fixture.writeStaged(responseWriter)
	default:
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (fixture *intentFixture) writeStaged(responseWriter http.ResponseWriter) {
	receipt := map[string]any{
		"ok": true,
		"data": map[string]any{
			"status": "staged",
			"staged_file": map[string]any{
				"staged_file_ref":  "staged_fixture",
				"state":            "verified",
				"source_kind":      fixture.prepared.SourceKind,
				"display_filename": fixture.prepared.DisplayFilename,
				"media_type":       fixture.prepared.ClaimedMediaType,
				"sha256":           fixture.prepared.ExpectedSHA256,
				"size":             fixture.prepared.ExpectedSize,
				"expires_at":       "2026-08-26T01:00:00Z",
				"replayed":         false,
			},
		},
	}
	_ = json.NewEncoder(responseWriter).Encode(receipt)
}

func testUploadService(t *testing.T, fixture *intentFixture) (*localFileServiceImpl, *httptest.Server) {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(fixture.handler))
	service, err := newLocalFileService(server.Client(), server.URL)
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = service.Close()
		server.Close()
	})
	return service, server
}

func prepareUploadFixture(t *testing.T, service *localFileServiceImpl, fixture *intentFixture, content []byte) prepareLocalFileResult {
	t.Helper()
	path := writeFixture(t, "fixture.bin", content)
	prepared, err := service.PrepareLocalFile(context.Background(), prepareLocalFileInput{
		Path: path, ClaimedMediaType: "application/octet-stream",
	})
	if err != nil {
		t.Fatal(err)
	}
	fixture.mu.Lock()
	fixture.prepared = prepared
	fixture.requestsClean = true
	fixture.mu.Unlock()
	return prepared
}

func TestUploadPreparedFileStreamsCredentiallessAndConsumesRef(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	fixture := &intentFixture{}
	service, server := testUploadService(t, fixture)
	content := []byte(strings.Repeat("opaque-file-stream-", 4096))
	prepared := prepareUploadFixture(t, service, fixture, content)
	result, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
		LocalFileRef: prepared.LocalFileRef,
		UploadURL:    server.URL + uploadIntentRoutePrefix + testCapability,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.StagedFileRef != "staged_fixture" || result.SHA256 != prepared.ExpectedSHA256 {
		t.Fatalf("unexpected staged receipt: %#v", result)
	}
	fixture.mu.Lock()
	if fixture.getCount != 1 || fixture.putCount != 1 || !fixture.requestsClean || string(fixture.body) != string(content) {
		t.Fatalf("unexpected transport: get=%d put=%d clean=%v bytes=%d", fixture.getCount, fixture.putCount, fixture.requestsClean, len(fixture.body))
	}
	fixture.mu.Unlock()
	_, err = service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
		LocalFileRef: prepared.LocalFileRef,
		UploadURL:    server.URL + uploadIntentRoutePrefix + testCapability,
	})
	assertLocalCode(t, err, "local_companion_ref_not_found")
}

func TestUploadPreparedFileReconcilesUnknownPut(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	fixture := &intentFixture{failPutAfterStage: true}
	service, server := testUploadService(t, fixture)
	prepared := prepareUploadFixture(t, service, fixture, []byte("reconcile exact snapshot"))
	result, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
		LocalFileRef: prepared.LocalFileRef,
		UploadURL:    server.URL + uploadIntentRoutePrefix + testCapability,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.StagedFileRef != "staged_fixture" {
		t.Fatalf("unexpected reconciliation result: %#v", result)
	}
	fixture.mu.Lock()
	defer fixture.mu.Unlock()
	if fixture.getCount != 2 || fixture.putCount != 1 {
		t.Fatalf("expected GET-before-PUT plus reconciliation, got GET=%d PUT=%d", fixture.getCount, fixture.putCount)
	}
}

func TestUploadPreparedFileRejectsChangedSnapshotBeforePut(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	fixture := &intentFixture{}
	service, server := testUploadService(t, fixture)
	path := writeFixture(t, "mutable.bin", []byte("before"))
	prepared, err := service.PrepareLocalFile(context.Background(), prepareLocalFileInput{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	fixture.prepared = prepared
	if err := os.WriteFile(path, []byte("after-change"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
		LocalFileRef: prepared.LocalFileRef,
		UploadURL:    server.URL + uploadIntentRoutePrefix + testCapability,
	})
	assertLocalCode(t, err, "local_companion_file_changed")
	fixture.mu.Lock()
	defer fixture.mu.Unlock()
	if fixture.getCount != 1 || fixture.putCount != 0 {
		t.Fatalf("changed file reached PUT: GET=%d PUT=%d", fixture.getCount, fixture.putCount)
	}
}

func TestUploadPreparedFileRejectsForeignOrNonExactURLWithoutNetwork(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	fixture := &intentFixture{}
	service, server := testUploadService(t, fixture)
	prepared := prepareUploadFixture(t, service, fixture, []byte("fixture"))
	for name, uploadURL := range map[string]string{
		"foreign":          "https://example.com" + uploadIntentRoutePrefix + testCapability,
		"query":            server.URL + uploadIntentRoutePrefix + testCapability + "?leak=1",
		"wrong route":      server.URL + "/api/other/" + testCapability,
		"short capability": server.URL + uploadIntentRoutePrefix + "mdupload_v1_short",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
				LocalFileRef: prepared.LocalFileRef, UploadURL: uploadURL,
			})
			assertLocalCode(t, err, "invalid_upload_url")
		})
	}
	fixture.mu.Lock()
	defer fixture.mu.Unlock()
	if fixture.getCount != 0 || fixture.putCount != 0 {
		t.Fatalf("invalid URL reached network: GET=%d PUT=%d", fixture.getCount, fixture.putCount)
	}
}

func TestUploadPreparedFileRetainsOnlyRetryableRejectedSnapshot(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("local companion is macOS-only")
	}
	t.Run("retryable", func(t *testing.T) {
		fixture := &intentFixture{rejectionCode: "capacity_soft_limit"}
		service, server := testUploadService(t, fixture)
		prepared := prepareUploadFixture(t, service, fixture, []byte("retryable fixture"))
		uploadURL := server.URL + uploadIntentRoutePrefix + testCapability
		_, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
			LocalFileRef: prepared.LocalFileRef, UploadURL: uploadURL,
		})
		assertLocalCode(t, err, "capacity_soft_limit")
		fixture.mu.Lock()
		fixture.rejectionCode = ""
		fixture.mu.Unlock()
		if _, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
			LocalFileRef: prepared.LocalFileRef, UploadURL: uploadURL,
		}); err != nil {
			t.Fatalf("retryable rejection consumed the exact local snapshot: %v", err)
		}
	})
	t.Run("definitive", func(t *testing.T) {
		fixture := &intentFixture{rejectionCode: "bundle_file_digest_mismatch"}
		service, server := testUploadService(t, fixture)
		prepared := prepareUploadFixture(t, service, fixture, []byte("definitive fixture"))
		uploadURL := server.URL + uploadIntentRoutePrefix + testCapability
		_, err := service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
			LocalFileRef: prepared.LocalFileRef, UploadURL: uploadURL,
		})
		assertLocalCode(t, err, "bundle_file_digest_mismatch")
		fixture.mu.Lock()
		fixture.rejectionCode = ""
		fixture.mu.Unlock()
		_, err = service.UploadPreparedFile(context.Background(), uploadPreparedFileInput{
			LocalFileRef: prepared.LocalFileRef, UploadURL: uploadURL,
		})
		assertLocalCode(t, err, "local_companion_ref_not_found")
	})
}
