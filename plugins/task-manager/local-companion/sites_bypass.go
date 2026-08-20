package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const sitesBypassHeader = "OAI-Sites-Authorization"

type secretSource interface {
	LoadSecret() (string, error)
}

type sitesBypassTransport struct {
	base    http.RoundTripper
	origin  *url.URL
	secrets secretSource
}

func newSitesBypassTransport(
	rawOrigin string,
	base http.RoundTripper,
	secrets secretSource,
) (*sitesBypassTransport, error) {
	origin, err := validateOAuthOrigin(rawOrigin)
	if err != nil || origin != uatTaskManagerOrigin || base == nil || secrets == nil {
		return nil, newLocalError(
			"invalid_private_uat_bridge",
			"Sites bypass is allowed only for the exact Task Manager UAT origin",
		)
	}
	parsed, _ := url.Parse(origin)
	return &sitesBypassTransport{base: base, origin: parsed, secrets: secrets}, nil
}

func (transport *sitesBypassTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request == nil || request.URL == nil || request.URL.Scheme != transport.origin.Scheme ||
		request.URL.Host != transport.origin.Host {
		return nil, errors.New("private UAT bridge rejected a request outside the exact UAT origin")
	}
	secret, err := transport.secrets.LoadSecret()
	if err != nil {
		return nil, err
	}
	secret = strings.TrimSpace(secret)
	if secret == "" || strings.ContainsAny(secret, "\r\n") {
		return nil, newLocalError("sites_bypass_invalid", "private UAT access credential is invalid")
	}
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	clone.Header.Set(sitesBypassHeader, "Bearer "+secret)
	return transport.base.RoundTrip(clone)
}
