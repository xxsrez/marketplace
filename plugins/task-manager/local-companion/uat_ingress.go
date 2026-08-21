package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	defaultPrivateUATIngressListen = "127.0.0.1:47821"
	sitesAuthorizationHeader       = "OAI-Sites-Authorization"
	maxIngressControlBody          = int64(1024 * 1024)
	maxIngressFileBody             = int64(26 * 1024 * 1024)
)

var (
	stagedFilePath               = regexp.MustCompile(`^/api/agent/v1/files/[A-Za-z0-9_-]{8,200}$`)
	taskAttachmentCollectionPath = regexp.MustCompile(
		`^/api/agent/v1/tasks/[A-Za-z0-9_-]{1,200}/attachments$`,
	)
	taskAttachmentPath = regexp.MustCompile(
		`^/api/agent/v1/tasks/[A-Za-z0-9_-]{1,200}/attachments/[A-Za-z0-9_-]{8,200}$`,
	)
	taskAttachmentContentPath = regexp.MustCompile(
		`^/api/agent/v1/tasks/[A-Za-z0-9_-]{1,200}/attachments/[A-Za-z0-9_-]{8,200}/content$`,
	)
)

type privateUATIngress struct {
	upstream *url.URL
	token    string
	client   *http.Client
}

func readSitesBypassToken(reader io.Reader) (string, error) {
	body, err := io.ReadAll(io.LimitReader(reader, 4097))
	if err != nil || len(body) == 0 || len(body) > 4096 {
		return "", newLocalError("uat_ingress_secret_invalid", "private UAT ingress token is missing or invalid")
	}
	token := strings.TrimSpace(string(body))
	if token == "" || strings.ContainsAny(token, "\r\n\x00") {
		return "", newLocalError("uat_ingress_secret_invalid", "private UAT ingress token is missing or invalid")
	}
	return token, nil
}

func runPrivateUATIngress(ctx context.Context, listenAddress, token string) error {
	if listenAddress == "" {
		listenAddress = defaultPrivateUATIngressListen
	}
	if err := validateLoopbackListenAddress(listenAddress); err != nil {
		return err
	}
	upstream, _ := url.Parse(uatTaskManagerOrigin)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	client := &http.Client{
		Timeout:   15 * time.Minute,
		Transport: transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	ingress := &privateUATIngress{upstream: upstream, token: token, client: client}
	listener, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		return newLocalError("uat_ingress_listen_failed", "private UAT ingress loopback port is unavailable")
	}
	server := &http.Server{
		Handler:           ingress,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ErrorLog:          log.New(io.Discard, "", 0),
	}
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	err = server.Serve(listener)
	if !errors.Is(err, http.ErrServerClosed) {
		return newLocalError("uat_ingress_failed", "private UAT ingress stopped unexpectedly")
	}
	<-shutdownDone
	return nil
}

func validateLoopbackListenAddress(address string) error {
	if address != defaultPrivateUATIngressListen {
		return newLocalError(
			"uat_ingress_listen_invalid", "private UAT ingress must bind exact 127.0.0.1:47821",
		)
	}
	return nil
}

func (ingress *privateUATIngress) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Host != defaultPrivateUATIngressListen {
		http.NotFound(writer, request)
		return
	}
	if !allowedPrivateUATRequest(request.Method, request.URL) {
		http.NotFound(writer, request)
		return
	}
	limit := maxIngressControlBody
	if request.Method == http.MethodPost && request.URL.Path == "/api/agent/v1/files" {
		limit = maxIngressFileBody
	}
	body, err := io.ReadAll(io.LimitReader(request.Body, limit+1))
	if err != nil || int64(len(body)) > limit {
		http.Error(writer, "request is too large", http.StatusRequestEntityTooLarge)
		return
	}
	target := *ingress.upstream
	target.Path = request.URL.Path
	target.RawQuery = request.URL.RawQuery
	outgoing, err := http.NewRequestWithContext(
		request.Context(), request.Method, target.String(), bytes.NewReader(body),
	)
	if err != nil {
		http.Error(writer, "upstream request failed", http.StatusBadGateway)
		return
	}
	copyIngressRequestHeaders(outgoing.Header, request.Header)
	outgoing.Header.Set(sitesAuthorizationHeader, "Bearer "+ingress.token)
	outgoing.Host = ingress.upstream.Host
	response, err := ingress.client.Do(outgoing)
	if err != nil || response == nil {
		http.Error(writer, "private UAT is unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 && response.StatusCode < 400 {
		http.Error(writer, "unexpected private UAT redirect", http.StatusBadGateway)
		return
	}
	responseLimit := maxIngressControlBody
	if request.Method == http.MethodGet && taskAttachmentContentPath.MatchString(request.URL.Path) {
		responseLimit = maxIngressFileBody
	}
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, responseLimit+1))
	if err != nil || int64(len(responseBody)) > responseLimit {
		http.Error(writer, "private UAT response is too large", http.StatusBadGateway)
		return
	}
	copyIngressResponseHeaders(writer.Header(), response.Header)
	writer.WriteHeader(response.StatusCode)
	_, _ = writer.Write(responseBody)
}

func allowedPrivateUATRequest(method string, target *url.URL) bool {
	if target == nil || target.RawPath != "" || target.Fragment != "" {
		return false
	}
	if target.RawQuery != "" {
		if method != http.MethodGet || !taskAttachmentContentPath.MatchString(target.Path) {
			return false
		}
		values, err := url.ParseQuery(target.RawQuery)
		if err != nil || len(values) != 1 {
			return false
		}
		variant := values.Get("variant")
		if variant != "original" && variant != "thumbnail" {
			return false
		}
	}
	switch target.Path {
	case "/oauth/register", "/oauth/token", "/oauth/revoke":
		return method == http.MethodPost
	case "/api/agent/v1/files":
		return method == http.MethodPost
	}
	if stagedFilePath.MatchString(target.Path) {
		return method == http.MethodGet || method == http.MethodDelete
	}
	if taskAttachmentCollectionPath.MatchString(target.Path) {
		return method == http.MethodGet || method == http.MethodPost
	}
	if taskAttachmentPath.MatchString(target.Path) {
		return method == http.MethodGet || method == http.MethodPatch || method == http.MethodDelete
	}
	if taskAttachmentContentPath.MatchString(target.Path) {
		return method == http.MethodGet
	}
	return false
}

func copyIngressRequestHeaders(target, source http.Header) {
	for _, name := range []string{
		"Accept", "Authorization", "Content-Type", "Idempotency-Key",
		"X-File-Filename", "X-File-Version", "X-Attachment-Version", "Range", "If-Range",
	} {
		for _, value := range source.Values(name) {
			target.Add(name, value)
		}
	}
	target.Set("User-Agent", "task-manager-local-uat-ingress/1")
}

func copyIngressResponseHeaders(target, source http.Header) {
	for _, name := range []string{
		"Accept-Ranges", "Cache-Control", "Content-Disposition", "Content-Range",
		"Content-Type", "ETag", "X-Content-Type-Options", "WWW-Authenticate",
	} {
		for _, value := range source.Values(name) {
			target.Add(name, value)
		}
	}
	target.Del("Set-Cookie")
	target.Del(sitesAuthorizationHeader)
}

func privateUATIngressCommandHint() string {
	return fmt.Sprintf(
		"start the on-demand private UAT ingress on %s and retry the same operation",
		defaultPrivateUATIngressListen,
	)
}

func exitWithLocalError(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
}
