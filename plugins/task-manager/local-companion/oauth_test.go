package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

type memoryCredentialStore struct {
	mutex       sync.Mutex
	record      oauthCredential
	hasRecord   bool
	saveCount   int
	deleteCount int
}

func (store *memoryCredentialStore) Load() (oauthCredential, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if !store.hasRecord {
		return oauthCredential{}, errCredentialNotFound
	}
	return store.record, nil
}

func (store *memoryCredentialStore) Save(record oauthCredential) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.record = record
	store.hasRecord = true
	store.saveCount++
	return nil
}

func (store *memoryCredentialStore) Delete() error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.record = oauthCredential{}
	store.hasRecord = false
	store.deleteCount++
	return nil
}

type callbackBrowser struct {
	mutex     sync.Mutex
	openCount int
}

func (browser *callbackBrowser) Open(rawURL string) error {
	browser.mutex.Lock()
	browser.openCount++
	browser.mutex.Unlock()
	authorization, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	callback, err := url.Parse(authorization.Query().Get("redirect_uri"))
	if err != nil {
		return err
	}
	query := callback.Query()
	query.Set("code", "authorization-code")
	query.Set("state", authorization.Query().Get("state"))
	callback.RawQuery = query.Encode()
	go func() {
		response, getErr := http.Get(callback.String())
		if getErr == nil {
			_ = response.Body.Close()
		}
	}()
	return nil
}

func TestOAuthFirstUseRefreshAndRevokedReconnect(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	registrationCount := 0
	refreshCount := 0
	rejectNextRefresh := false
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/oauth/register":
			registrationCount++
			var metadata struct {
				RedirectURIs []string `json:"redirect_uris"`
				TokenMethod  string   `json:"token_endpoint_auth_method"`
			}
			if err := json.NewDecoder(request.Body).Decode(&metadata); err != nil {
				t.Errorf("registration body: %v", err)
			}
			if len(metadata.RedirectURIs) != 1 || metadata.TokenMethod != "none" {
				t.Errorf("registration metadata = %#v", metadata)
			}
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"client_id": "tm_oauth_client_test", "redirect_uris": metadata.RedirectURIs,
			})
		case "/oauth/token":
			if err := request.ParseForm(); err != nil {
				t.Errorf("token form: %v", err)
			}
			grant := request.Form.Get("grant_type")
			if request.Form.Get("client_secret") != "" {
				t.Error("public client sent a secret")
			}
			mutex.Lock()
			if grant == "refresh_token" {
				refreshCount++
				if rejectNextRefresh {
					rejectNextRefresh = false
					mutex.Unlock()
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(writer).Encode(map[string]any{"error": "invalid_grant"})
					return
				}
			}
			sequence := refreshCount + registrationCount
			mutex.Unlock()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"access_token": "tm_oat_test", "refresh_token": "tm_ort_rotated_" + string(rune('0'+sequence)),
				"token_type": "Bearer", "expires_in": 900, "scope": "api:read api:write",
			})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	store := &memoryCredentialStore{}
	browser := &callbackBrowser{}
	manager, err := newOAuthManager(server.URL, server.Client(), store, browser)
	if err != nil {
		t.Fatal(err)
	}
	token, err := manager.Token(context.Background())
	if err != nil || token != "tm_oat_test" {
		t.Fatalf("first token = %q, error = %v", token, err)
	}
	if registrationCount != 1 || browser.openCount != 1 || store.saveCount != 1 {
		t.Fatalf("first use: registrations=%d browser=%d saves=%d", registrationCount, browser.openCount, store.saveCount)
	}

	manager.Invalidate()
	if _, err := manager.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	if refreshCount != 1 || browser.openCount != 1 || store.saveCount != 2 {
		t.Fatalf("refresh: refreshes=%d browser=%d saves=%d", refreshCount, browser.openCount, store.saveCount)
	}

	manager.Invalidate()
	mutex.Lock()
	rejectNextRefresh = true
	mutex.Unlock()
	if _, err := manager.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	if refreshCount != 2 || registrationCount != 2 || browser.openCount != 2 || store.deleteCount != 1 {
		t.Fatalf("reconnect: refreshes=%d registrations=%d browser=%d deletes=%d", refreshCount, registrationCount, browser.openCount, store.deleteCount)
	}
}

func TestOAuthRejectsNonHTTPSNonLoopbackOrigin(t *testing.T) {
	t.Parallel()
	_, err := newOAuthManager("http://example.test", http.DefaultClient, &memoryCredentialStore{}, &callbackBrowser{})
	if err == nil {
		t.Fatal("expected unsafe origin rejection")
	}
	var localErr *localError
	if !errors.As(err, &localErr) || localErr.Code != "invalid_origin" {
		t.Fatalf("error = %v", err)
	}
}

func TestOAuthRegistrationEdgeFailureIsActionableAndBounded(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth/register" {
			http.NotFound(writer, request)
			return
		}
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`<html><body>Sign in required: private edge detail</body></html>`))
	}))
	defer server.Close()

	manager, err := newOAuthManager(
		server.URL, server.Client(), &memoryCredentialStore{}, &callbackBrowser{},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = manager.Token(context.Background())
	var localErr *localError
	if !errors.As(err, &localErr) {
		t.Fatalf("error type = %T: %v", err, err)
	}
	if localErr.Code != "oauth_registration_blocked" {
		t.Fatalf("error code = %q", localErr.Code)
	}
	if !strings.Contains(localErr.Message, "HTTP 401") {
		t.Fatalf("error message = %q", localErr.Message)
	}
	if strings.Contains(localErr.Message, "private edge detail") || strings.Contains(localErr.Message, "<html>") {
		t.Fatalf("error disclosed edge response body: %q", localErr.Message)
	}
}
