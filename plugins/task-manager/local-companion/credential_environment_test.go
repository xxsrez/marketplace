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
	"testing"
)

func TestEnvironmentTokenUploadAndBind(t *testing.T) {
	const fixtureToken = "tm_pat_test-fixture-not-a-real-credential"
	body := []byte("container fixture")
	digest := sha256.Sum256(body)
	checksum := hex.EncodeToString(digest[:])
	uploaded, bound := false, false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+fixtureToken {
			t.Error("wrong authorization")
		}
		ref := "file_fixture"
		switch r.URL.Path {
		case "/api/agent/v1/files":
			received, _ := io.ReadAll(r.Body)
			if string(received) != string(body) || r.Header.Get("Idempotency-Key") != "upload" {
				t.Error("wrong upload")
			}
			uploaded = true
		case "/api/agent/v1/tasks/TM-1/attachments":
			var payload map[string]string
			if json.NewDecoder(r.Body).Decode(&payload) != nil || payload["fileRef"] != "file_fixture" || payload["idempotencyKey"] != "bind" || !uploaded {
				t.Error("wrong bind")
			}
			ref, bound = "attachment_fixture", true
		default:
			t.Error("unexpected route, including OAuth")
			w.WriteHeader(404)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{
			"ref": ref, "filename": "fixture.txt", "mediaType": "text/plain",
			"byteSize": len(body), "checksumSha256": checksum, "state": "ready", "kind": "file", "version": 1,
		}})
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "fixture.txt")
	if err := os.WriteFile(path, body, 0600); err != nil {
		t.Fatal(err)
	}
	source := &environmentTokenSource{lookup: func(string) string { return fixtureToken }}
	client := newUploadClient(server.Client(), server.URL, source)
	file, err := client.UploadLocalFile(context.Background(), localFileInput{Path: path, IdempotencyKey: "upload", ExpectedSHA256: checksum})
	if err != nil {
		t.Fatal(err)
	}
	attachment, err := client.AttachLocalFileToTask(context.Background(), attachLocalFileInput{TaskRef: "TM-1", FileRef: file.FileRef, IdempotencyKey: "bind"})
	if err != nil {
		t.Fatal(err)
	}
	if !file.SourceVerified || !bound || attachment.AttachmentRef != "attachment_fixture" {
		t.Fatal("upload/bind incomplete")
	}
}

func TestEnvironmentTokenValidation(t *testing.T) {
	for _, test := range []struct{ value, code string }{
		{"", "environment_token_required"},
		{"tm_pat_", "environment_token_invalid"},
		{"tm_oat_fixture", "environment_token_invalid"},
		{"tm_pat_fixture\n", "environment_token_invalid"},
		{"tm_pat_" + strings.Repeat("a", 4096), "environment_token_invalid"},
	} {
		source := &environmentTokenSource{lookup: func(name string) string {
			if name != localTokenEnvironment {
				t.Fatal("unexpected secret name")
			}
			return test.value
		}}
		_, err := source.Token(context.Background())
		assertLocalErrorCode(t, err, test.code)
		if test.value != "" && strings.Contains(err.Error(), test.value) {
			t.Fatal("error exposed secret")
		}
	}
}

func TestEnvironmentTokenRejectsUnauthorizedWithoutRetryingCredential(t *testing.T) {
	const fixtureToken = "tm_pat_test-fixture-not-a-real-credential"
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("Authorization") != "Bearer "+fixtureToken {
			t.Error("missing expected authorization")
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	source := &environmentTokenSource{lookup: func(string) string { return fixtureToken }}
	client := newUploadClient(server.Client(), server.URL, source)
	path := filepath.Join(t.TempDir(), "fixture.txt")
	if err := os.WriteFile(path, []byte("test fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := client.UploadLocalFile(context.Background(), localFileInput{Path: path, IdempotencyKey: "fixture"})
	assertLocalErrorCode(t, err, "environment_token_rejected")
	_, err = client.AttachLocalFileToTask(context.Background(), attachLocalFileInput{TaskRef: "TM-1", FileRef: "fixture", IdempotencyKey: "bind"})
	assertLocalErrorCode(t, err, "environment_token_rejected")
	if calls != 1 {
		t.Fatalf("requests = %d, want one", calls)
	}
}

func TestEnvironmentTokenMissingDoesNotMakeNetworkRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("missing credential must not reach network")
	}))
	defer server.Close()
	source := &environmentTokenSource{lookup: func(string) string { return "" }}
	client := newUploadClient(server.Client(), server.URL, source)
	_, err := client.AttachLocalFileToTask(context.Background(), attachLocalFileInput{TaskRef: "TM-1", FileRef: "fixture", IdempotencyKey: "bind"})
	assertLocalErrorCode(t, err, "environment_token_required")
}
