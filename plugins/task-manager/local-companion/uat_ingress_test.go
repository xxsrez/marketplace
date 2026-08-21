package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPrivateUATIngressAllowsOnlyBoundedFileRoutesAndInjectsSecret(t *testing.T) {
	t.Parallel()
	var captured struct {
		path          string
		authorization string
		sitesToken    string
		cookie        string
		body          string
	}
	upstream := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		captured.path = request.URL.Path
		captured.authorization = request.Header.Get("Authorization")
		captured.sitesToken = request.Header.Get(sitesAuthorizationHeader)
		captured.cookie = request.Header.Get("Cookie")
		body, _ := io.ReadAll(request.Body)
		captured.body = string(body)
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Set-Cookie", "private=secret")
		writer.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(writer, `{"data":{"ref":"att_12345678"}}`)
	}))
	defer upstream.Close()
	upstreamURL, _ := url.Parse(upstream.URL)
	ingress := &privateUATIngress{
		upstream: upstreamURL, token: "sites-secret",
		client: upstream.Client(),
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"http://127.0.0.1:47821/api/agent/v1/tasks/TM-123/attachments",
		strings.NewReader(`{"fileRef":"file_12345678","idempotencyKey":"bind"}`),
	)
	request.Header.Set("Authorization", "Bearer task-token")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "must-not-cross=true")
	request.Header.Set(sitesAuthorizationHeader, "attacker-value")
	response := httptest.NewRecorder()
	ingress.ServeHTTP(response, request)

	if response.Code != http.StatusCreated || captured.path != request.URL.Path ||
		captured.authorization != "Bearer task-token" || captured.sitesToken != "Bearer sites-secret" {
		t.Fatalf("response=%d captured=%#v", response.Code, captured)
	}
	if captured.cookie != "" || response.Header().Get("Set-Cookie") != "" ||
		strings.Contains(response.Body.String(), "sites-secret") {
		t.Fatalf("secret-bearing headers crossed ingress: captured=%#v headers=%#v", captured, response.Header())
	}
	if captured.body != `{"fileRef":"file_12345678","idempotencyKey":"bind"}` {
		t.Fatalf("body = %q", captured.body)
	}
}

func TestPrivateUATIngressRejectsMCPUIAndUnexpectedQueries(t *testing.T) {
	t.Parallel()
	for _, item := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/mcp"},
		{http.MethodGet, "/workspace"},
		{http.MethodGet, "/oauth/authorize"},
		{http.MethodPost, "/api/agent/v1/tasks/TM-123"},
		{http.MethodPost, "/api/agent/v1/files?unexpected=1"},
		{http.MethodGet, "/api/agent/v1/tasks/TM-123/attachments/att_12345678/content?variant=raw"},
	} {
		requestURL, _ := url.Parse("http://127.0.0.1:47821" + item.path)
		if allowedPrivateUATRequest(item.method, requestURL) {
			t.Fatalf("allowed %s %s", item.method, item.path)
		}
	}
	for _, item := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/oauth/register"},
		{http.MethodPost, "/oauth/token"},
		{http.MethodPost, "/api/agent/v1/files"},
		{http.MethodGet, "/api/agent/v1/files/file_12345678"},
		{http.MethodPost, "/api/agent/v1/tasks/TM-123/attachments"},
		{http.MethodGet, "/api/agent/v1/tasks/TM-123/attachments/att_12345678"},
		{http.MethodGet, "/api/agent/v1/tasks/TM-123/attachments/att_12345678/content?variant=original"},
	} {
		requestURL, _ := url.Parse("http://127.0.0.1:47821" + item.path)
		if !allowedPrivateUATRequest(item.method, requestURL) {
			t.Fatalf("rejected %s %s", item.method, item.path)
		}
	}
}

func TestPrivateUATIngressRejectsRedirectsAndOversizedBodies(t *testing.T) {
	t.Parallel()
	upstream := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Location", "https://example.com/leak")
		writer.WriteHeader(http.StatusFound)
	}))
	defer upstream.Close()
	upstreamURL, _ := url.Parse(upstream.URL)
	client := upstream.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	ingress := &privateUATIngress{upstream: upstreamURL, token: "token", client: client}
	redirect := httptest.NewRecorder()
	ingress.ServeHTTP(redirect, httptest.NewRequest(
		http.MethodPost, "http://127.0.0.1:47821/oauth/token", strings.NewReader("grant_type=refresh_token"),
	))
	if redirect.Code != http.StatusBadGateway || redirect.Header().Get("Location") != "" {
		t.Fatalf("redirect response = %d %#v", redirect.Code, redirect.Header())
	}

	oversized := httptest.NewRecorder()
	ingress.ServeHTTP(oversized, httptest.NewRequest(
		http.MethodPost, "http://127.0.0.1:47821/oauth/token",
		strings.NewReader(strings.Repeat("x", int(maxIngressControlBody)+1)),
	))
	if oversized.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized response = %d", oversized.Code)
	}
}

func TestPrivateUATIngressSecretAndTransportOriginsFailClosed(t *testing.T) {
	t.Parallel()
	if token, err := readSitesBypassToken(strings.NewReader("  sites-token  \n")); err != nil || token != "sites-token" {
		t.Fatalf("token=%q error=%v", token, err)
	}
	if _, err := readSitesBypassToken(strings.NewReader("\n")); err == nil {
		t.Fatal("expected empty token rejection")
	}
	if _, err := validateTransportOrigin("http://example.test:47821", uatTaskManagerOrigin); err == nil {
		t.Fatal("expected non-loopback transport rejection")
	}
	if got, err := validateTransportOrigin("http://127.0.0.1:47821", uatTaskManagerOrigin); err != nil ||
		got != "http://127.0.0.1:47821" {
		t.Fatalf("transport=%q error=%v", got, err)
	}
	if err := validateLoopbackListenAddress("0.0.0.0:47821"); err == nil {
		t.Fatal("expected wildcard listen rejection")
	}
	if err := validateLoopbackListenAddress("127.0.0.1:47822"); err == nil {
		t.Fatal("expected alternate loopback port rejection")
	}
}

func TestRunPrivateUATIngressStopsWithContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// Address validation occurs before any network operation. A canceled context
	// must not weaken the fixed-loopback requirement.
	if err := runPrivateUATIngress(ctx, "0.0.0.0:47821", "token"); err == nil {
		t.Fatal("expected invalid listen address")
	}
}
