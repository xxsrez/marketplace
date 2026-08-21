package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type staticTokenSource struct {
	token string
}

func (source *staticTokenSource) Token(context.Context) (string, error) {
	return source.token, nil
}

func (source *staticTokenSource) Invalidate() {}

func TestUploadLocalFileSendsOnlySafeMetadataAndVerifiesResponse(t *testing.T) {
	t.Parallel()
	body := []byte("%PDF-1.7\nprivate local fixture\n")
	checksum := sha256.Sum256(body)
	path := filepath.Join(t.TempDir(), "private-directory", "report.pdf")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}

	var captured struct {
		path        string
		filename    string
		contentType string
		idempotency string
		body        []byte
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		captured.path = request.URL.Path
		captured.filename = request.Header.Get("X-File-Filename")
		captured.contentType = request.Header.Get("Content-Type")
		captured.idempotency = request.Header.Get("Idempotency-Key")
		captured.body, _ = io.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(map[string]any{
			"data": map[string]any{
				"ref": "file_test_123", "filename": "report.pdf",
				"mediaType": "application/pdf", "byteSize": len(body),
				"checksumSha256": hex.EncodeToString(checksum[:]),
				"kind":           "file", "state": "ready", "version": 1,
				"readyExpiresAt": "2026-08-21T00:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	client := newUploadClient(server.Client(), server.URL, &staticTokenSource{token: "access-token"})
	result, err := client.UploadLocalFile(context.Background(), localFileInput{
		Path: path, IdempotencyKey: "stable-upload-key",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.FileRef != "file_test_123" || !result.SourceVerified {
		t.Fatalf("result = %#v", result)
	}
	if captured.path != "/api/agent/v1/files" {
		t.Fatalf("request path = %q", captured.path)
	}
	if captured.filename != "report.pdf" || captured.contentType != "application/pdf" {
		t.Fatalf("headers = filename %q content-type %q", captured.filename, captured.contentType)
	}
	if captured.idempotency != "stable-upload-key" || string(captured.body) != string(body) {
		t.Fatal("request payload differs")
	}
	capturedJSON, _ := json.Marshal(captured)
	if strings.Contains(string(capturedJSON), path) || strings.Contains(string(capturedJSON), "private-directory") {
		t.Fatalf("request disclosed local path: %s", capturedJSON)
	}
}

func TestUploadLocalFileRetryUsesStableIdempotencyWithoutDuplicate(t *testing.T) {
	t.Parallel()
	body := []byte("retry-safe")
	checksum := sha256.Sum256(body)
	path := filepath.Join(t.TempDir(), "retry.txt")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	var mutex sync.Mutex
	refsByKey := map[string]string{}
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		requestCount++
		key := request.Header.Get("Idempotency-Key")
		ref := refsByKey[key]
		if ref == "" {
			ref = "file_retry_123"
			refsByKey[key] = ref
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(map[string]any{"data": map[string]any{
			"ref": ref, "filename": "retry.txt", "mediaType": "text/plain",
			"byteSize": len(body), "checksumSha256": hex.EncodeToString(checksum[:]),
			"kind": "file", "state": "ready", "version": 1,
		}})
	}))
	defer server.Close()
	client := newUploadClient(server.Client(), server.URL, &staticTokenSource{token: "token"})
	for range 2 {
		result, err := client.UploadLocalFile(context.Background(), localFileInput{
			Path: path, IdempotencyKey: "same-logical-upload",
		})
		if err != nil || result.FileRef != "file_retry_123" {
			t.Fatalf("result = %#v, error = %v", result, err)
		}
	}
	if requestCount != 2 || len(refsByKey) != 1 {
		t.Fatalf("requests = %d, logical files = %d", requestCount, len(refsByKey))
	}
}

func TestUploadLocalFileRejectsServerMetadataMismatch(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "mismatch.txt")
	if err := os.WriteFile(path, []byte("local bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(writer, `{"data":{"ref":"file_bad","filename":"mismatch.txt","mediaType":"text/plain","byteSize":11,"checksumSha256":"`+strings.Repeat("0", 64)+`","kind":"file","state":"ready","version":1}}`)
	}))
	defer server.Close()
	client := newUploadClient(server.Client(), server.URL, &staticTokenSource{token: "token"})
	_, err := client.UploadLocalFile(context.Background(), localFileInput{
		Path: path, IdempotencyKey: "mismatch",
	})
	assertLocalErrorCode(t, err, "server_metadata_mismatch")
}

func TestAttachLocalFileToTaskUsesAgentRESTAndReturnsAttachmentRef(t *testing.T) {
	t.Parallel()
	var captured struct {
		path          string
		authorization string
		contentType   string
		body          map[string]string
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		captured.path = request.URL.Path
		captured.authorization = request.Header.Get("Authorization")
		captured.contentType = request.Header.Get("Content-Type")
		_ = json.NewDecoder(request.Body).Decode(&captured.body)
		original := "https://task-manager.test/api/agent/v1/tasks/TM-123/attachments/att_12345678/content?variant=original"
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(map[string]any{"data": map[string]any{
			"ref": "att_12345678", "filename": "bound-report.pdf",
			"mediaType": "application/pdf", "byteSize": 41,
			"checksumSha256": strings.Repeat("b", 64), "kind": "file",
			"state": "ready", "version": 1,
			"links": map[string]any{"original": original, "thumbnail": nil},
		}})
	}))
	defer server.Close()

	client := newUploadClient(server.Client(), server.URL, &staticTokenSource{token: "access-token"})
	result, err := client.AttachLocalFileToTask(context.Background(), attachLocalFileInput{
		TaskRef: "TM-123", FileRef: "file_12345678",
		IdempotencyKey: "stable-bind-key", DisplayName: "bound-report.pdf",
	})
	if err != nil {
		t.Fatal(err)
	}
	if captured.path != "/api/agent/v1/tasks/TM-123/attachments" ||
		captured.authorization != "Bearer access-token" ||
		captured.contentType != "application/json" {
		t.Fatalf("request = %#v", captured)
	}
	if captured.body["fileRef"] != "file_12345678" ||
		captured.body["idempotencyKey"] != "stable-bind-key" ||
		captured.body["displayName"] != "bound-report.pdf" {
		t.Fatalf("body = %#v", captured.body)
	}
	if result.TaskRef != "TM-123" || result.AttachmentRef != "att_12345678" ||
		result.ChecksumSHA256 != strings.Repeat("b", 64) || result.OriginalURL == nil {
		t.Fatalf("result = %#v", result)
	}
}

func TestAttachLocalFileToTaskRejectsInvalidInputsBeforeNetwork(t *testing.T) {
	t.Parallel()
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		requestCount++
	}))
	defer server.Close()
	client := newUploadClient(server.Client(), server.URL, &staticTokenSource{token: "token"})
	_, err := client.AttachLocalFileToTask(context.Background(), attachLocalFileInput{
		TaskRef: "", FileRef: "file_12345678", IdempotencyKey: "bind-key",
	})
	assertLocalErrorCode(t, err, "invalid_arguments")
	if requestCount != 0 {
		t.Fatalf("invalid bind performed %d requests", requestCount)
	}
}
