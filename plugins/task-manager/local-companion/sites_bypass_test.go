package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type staticSecretSource struct {
	secret string
	err    error
}

func (source *staticSecretSource) LoadSecret() (string, error) {
	return source.secret, source.err
}

type recordingRoundTripper struct {
	request *http.Request
	calls   int
}

func (transport *recordingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	transport.calls++
	transport.request = request
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"jsonrpc":"2.0","id":1,"result":{}}`)),
		Request:    request,
	}, nil
}

func TestSitesBypassTransportInjectsOnlyForExactPrivateUATOrigin(t *testing.T) {
	t.Parallel()
	base := &recordingRoundTripper{}
	transport, err := newSitesBypassTransport(
		uatTaskManagerOrigin,
		base,
		&staticSecretSource{secret: "sites-bypass-test-secret"},
	)
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodPost, uatTaskManagerOrigin+"/api/mcp", nil)
	request.Header.Set("Authorization", "Bearer tm_oat_test")
	response, err := transport.RoundTrip(request)
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if base.calls != 1 {
		t.Fatalf("base calls = %d", base.calls)
	}
	if got := base.request.Header.Get(sitesBypassHeader); got != "Bearer sites-bypass-test-secret" {
		t.Fatalf("Sites header = %q", got)
	}
	if got := base.request.Header.Get("Authorization"); got != "Bearer tm_oat_test" {
		t.Fatalf("Task Manager authorization = %q", got)
	}
	if request.Header.Get(sitesBypassHeader) != "" {
		t.Fatal("original request was mutated")
	}

	foreign, _ := http.NewRequest(http.MethodGet, "https://example.test/api/mcp", nil)
	if _, err := transport.RoundTrip(foreign); err == nil {
		t.Fatal("expected cross-origin request rejection")
	}
	if base.calls != 1 {
		t.Fatal("cross-origin request reached the underlying transport")
	}
}

func TestSitesBypassTransportRefusesProductionOrigin(t *testing.T) {
	t.Parallel()
	_, err := newSitesBypassTransport(
		productionTaskManagerOrigin,
		&recordingRoundTripper{},
		&staticSecretSource{secret: "sites-bypass-test-secret"},
	)
	if err == nil {
		t.Fatal("expected production origin rejection")
	}
}
