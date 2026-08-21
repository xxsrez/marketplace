package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var errCredentialNotFound = errors.New("OAuth credential not found")

type oauthCredential struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	RefreshToken string `json:"refresh_token"`
}

type credentialStore interface {
	Load() (oauthCredential, error)
	Save(oauthCredential) error
	Delete() error
}

type browserOpener interface {
	Open(string) error
}

type oauthManager struct {
	mutex           sync.Mutex
	origin          string
	transportOrigin string
	resource        string
	httpClient      *http.Client
	credentials     credentialStore
	browser         browserOpener
	accessToken     string
	accessExpiry    time.Time
}

type oauthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type oauthProtocolError struct {
	Code    string
	Message string
	Status  int
}

func (err *oauthProtocolError) Error() string {
	if err.Message == "" {
		return err.Code
	}
	return err.Code + ": " + err.Message
}

func newOAuthManager(
	rawOrigin string,
	rawTransportOrigin string,
	httpClient *http.Client,
	credentials credentialStore,
	browser browserOpener,
) (*oauthManager, error) {
	origin, err := validateOAuthOrigin(rawOrigin)
	if err != nil {
		return nil, err
	}
	transportOrigin, err := validateTransportOrigin(rawTransportOrigin, origin)
	if err != nil {
		return nil, err
	}
	return &oauthManager{
		origin: origin, transportOrigin: transportOrigin,
		resource: origin + "/api/mcp", httpClient: httpClient,
		credentials: credentials, browser: browser,
	}, nil
}

func (manager *oauthManager) Token(ctx context.Context) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if manager.accessToken != "" && time.Until(manager.accessExpiry) > 30*time.Second {
		return manager.accessToken, nil
	}
	credential, err := manager.credentials.Load()
	if err == nil {
		tokens, refreshErr := manager.refresh(ctx, credential)
		if refreshErr == nil {
			credential.RefreshToken = tokens.RefreshToken
			if err := manager.credentials.Save(credential); err != nil {
				return "", err
			}
			manager.remember(tokens)
			return manager.accessToken, nil
		}
		var protocolErr *oauthProtocolError
		if !errors.As(refreshErr, &protocolErr) || protocolErr.Code != "invalid_grant" {
			return "", refreshErr
		}
		if err := manager.credentials.Delete(); err != nil {
			return "", err
		}
	} else if !errors.Is(err, errCredentialNotFound) {
		return "", err
	}
	credential, tokens, err := manager.authorize(ctx)
	if err != nil {
		return "", err
	}
	if err := manager.credentials.Save(credential); err != nil {
		return "", err
	}
	manager.remember(tokens)
	return manager.accessToken, nil
}

func (manager *oauthManager) Invalidate() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.accessToken = ""
	manager.accessExpiry = time.Time{}
}

func (manager *oauthManager) remember(tokens oauthTokens) {
	manager.accessToken = tokens.AccessToken
	expiresIn := tokens.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 15 * 60
	}
	manager.accessExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (manager *oauthManager) authorize(ctx context.Context) (oauthCredential, oauthTokens, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return oauthCredential{}, oauthTokens{}, newLocalError(
			"oauth_loopback_unavailable", "a loopback callback could not be opened",
		)
	}
	defer listener.Close()
	redirectURI := "http://" + listener.Addr().String() + "/oauth/callback"
	clientID, err := manager.register(ctx, redirectURI)
	if err != nil {
		return oauthCredential{}, oauthTokens{}, err
	}
	verifier, err := randomBase64URL(64)
	if err != nil {
		return oauthCredential{}, oauthTokens{}, newLocalError("oauth_random_failed", "OAuth PKCE state could not be created")
	}
	state, err := randomBase64URL(32)
	if err != nil {
		return oauthCredential{}, oauthTokens{}, newLocalError("oauth_random_failed", "OAuth state could not be created")
	}
	challengeDigest := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(challengeDigest[:])

	resultChannel := make(chan callbackResult, 1)
	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method != http.MethodGet || request.URL.Path != "/oauth/callback" {
				http.NotFound(writer, request)
				return
			}
			if request.URL.Query().Get("state") != state {
				http.Error(writer, "OAuth callback state did not match.", http.StatusBadRequest)
				return
			}
			result := callbackResult{
				Code: request.URL.Query().Get("code"), State: request.URL.Query().Get("state"),
				Error: request.URL.Query().Get("error"),
			}
			select {
			case resultChannel <- result:
			default:
			}
			writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			writer.Header().Set("Cache-Control", "no-store")
			_, _ = io.WriteString(writer, "Task Manager authorization received. You may close this tab.\n")
		}),
	}
	go func() { _ = server.Serve(listener) }()
	defer server.Shutdown(context.Background())

	authorizationURL := manager.origin + "/oauth/authorize?" + url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"resource":              {manager.resource},
		"scope":                 {"api:write"},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}.Encode()
	if err := manager.browser.Open(authorizationURL); err != nil {
		return oauthCredential{}, oauthTokens{}, newLocalError(
			"oauth_browser_failed", "the system browser could not open Task Manager authorization",
		)
	}
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	var callback callbackResult
	select {
	case <-waitCtx.Done():
		return oauthCredential{}, oauthTokens{}, newLocalError(
			"oauth_callback_timeout", "Task Manager authorization did not complete within ten minutes",
		)
	case callback = <-resultChannel:
	}
	if callback.Error != "" {
		return oauthCredential{}, oauthTokens{}, newLocalError(
			"oauth_access_denied", "Task Manager authorization was denied or canceled",
		)
	}
	if callback.State != state || callback.Code == "" {
		return oauthCredential{}, oauthTokens{}, newLocalError(
			"oauth_callback_invalid", "Task Manager authorization callback did not match this request",
		)
	}
	tokens, err := manager.exchangeCode(ctx, clientID, redirectURI, verifier, callback.Code)
	if err != nil {
		return oauthCredential{}, oauthTokens{}, err
	}
	return oauthCredential{
		ClientID: clientID, RedirectURI: redirectURI, RefreshToken: tokens.RefreshToken,
	}, tokens, nil
}

type callbackResult struct {
	Code  string
	State string
	Error string
}

func (manager *oauthManager) register(ctx context.Context, redirectURI string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"client_name":                "Task Manager local Codex companion",
		"redirect_uris":              []string{redirectURI},
		"grant_types":                []string{"authorization_code", "refresh_token"},
		"response_types":             []string{"code"},
		"token_endpoint_auth_method": "none",
	})
	request, err := http.NewRequestWithContext(
		ctx, http.MethodPost, manager.transportOrigin+"/oauth/register", strings.NewReader(string(body)),
	)
	if err != nil {
		return "", newLocalError("oauth_registration_failed", "OAuth registration request could not be created")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	responseBody, status, err := manager.doBounded(request)
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		return "", registrationError(status)
	}
	var response struct {
		ClientID string `json:"client_id"`
	}
	if json.Unmarshal(responseBody, &response) != nil || response.ClientID == "" {
		return "", newLocalError("oauth_registration_failed", "OAuth registration returned invalid metadata")
	}
	return response.ClientID, nil
}

func registrationError(status int) error {
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return newLocalError(
			"oauth_registration_blocked",
			fmt.Sprintf("Task Manager OAuth registration is blocked before user consent (HTTP %d)", status),
		)
	}
	return newLocalError(
		"oauth_registration_failed",
		fmt.Sprintf("Task Manager OAuth registration was rejected (HTTP %d)", status),
	)
}

func (manager *oauthManager) exchangeCode(
	ctx context.Context,
	clientID, redirectURI, verifier, code string,
) (oauthTokens, error) {
	return manager.requestTokens(ctx, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"resource":      {manager.resource},
		"code_verifier": {verifier},
	})
}

func (manager *oauthManager) refresh(ctx context.Context, credential oauthCredential) (oauthTokens, error) {
	if credential.ClientID == "" || credential.RefreshToken == "" {
		return oauthTokens{}, &oauthProtocolError{Code: "invalid_grant", Message: "stored OAuth credential is incomplete"}
	}
	return manager.requestTokens(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {credential.RefreshToken},
		"client_id":     {credential.ClientID},
		"resource":      {manager.resource},
	})
}

func (manager *oauthManager) requestTokens(ctx context.Context, form url.Values) (oauthTokens, error) {
	request, err := http.NewRequestWithContext(
		ctx, http.MethodPost, manager.transportOrigin+"/oauth/token", strings.NewReader(form.Encode()),
	)
	if err != nil {
		return oauthTokens{}, newLocalError("oauth_token_failed", "OAuth token request could not be created")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	body, status, err := manager.doBounded(request)
	if err != nil {
		return oauthTokens{}, err
	}
	if status != http.StatusOK {
		return oauthTokens{}, parseOAuthError(status, body)
	}
	var tokens oauthTokens
	if json.Unmarshal(body, &tokens) != nil || tokens.AccessToken == "" ||
		tokens.RefreshToken == "" || !strings.EqualFold(tokens.TokenType, "Bearer") {
		return oauthTokens{}, newLocalError("oauth_token_failed", "OAuth token endpoint returned invalid metadata")
	}
	return tokens, nil
}

func (manager *oauthManager) doBounded(request *http.Request) ([]byte, int, error) {
	response, err := manager.httpClient.Do(request)
	if err != nil {
		message := "Task Manager OAuth endpoint is unavailable"
		if transport, parseErr := url.Parse(manager.transportOrigin); parseErr == nil &&
			transport.Scheme == "http" && transport.Hostname() == "127.0.0.1" {
			message = privateUATIngressCommandHint()
		}
		return nil, 0, newLocalError("oauth_network_error", message)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 256*1024+1))
	if err != nil || len(body) > 256*1024 {
		return nil, response.StatusCode, newLocalError("oauth_response_invalid", "Task Manager OAuth response is invalid")
	}
	return body, response.StatusCode, nil
}

func parseOAuthError(status int, body []byte) error {
	var response struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &response)
	if response.Error == "" {
		response.Error = fmt.Sprintf("http_%d", status)
	}
	return &oauthProtocolError{Code: response.Error, Message: response.Description, Status: status}
}

func randomBase64URL(byteCount int) (string, error) {
	bytes := make([]byte, byteCount)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func validateOAuthOrigin(rawOrigin string) (string, error) {
	parsed, err := url.Parse(strings.TrimRight(rawOrigin, "/"))
	if err != nil || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" ||
		parsed.Path != "" {
		return "", newLocalError("invalid_origin", "Task Manager origin is invalid")
	}
	loopback := parsed.Scheme == "http" &&
		(parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "localhost" || parsed.Hostname() == "::1")
	if parsed.Scheme != "https" && !loopback {
		return "", newLocalError("invalid_origin", "Task Manager origin must use HTTPS")
	}
	return parsed.String(), nil
}

func validateTransportOrigin(rawTransportOrigin, logicalOrigin string) (string, error) {
	value := strings.TrimSpace(rawTransportOrigin)
	if value == "" {
		return logicalOrigin, nil
	}
	parsed, err := url.Parse(strings.TrimRight(value, "/"))
	if err != nil || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" ||
		parsed.Fragment != "" || parsed.Path != "" {
		return "", newLocalError("invalid_transport_origin", "Task Manager transport origin is invalid")
	}
	if parsed.Scheme == "https" && parsed.String() == logicalOrigin {
		return parsed.String(), nil
	}
	if parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" {
		return "", newLocalError(
			"invalid_transport_origin", "Task Manager transport origin must be the logical HTTPS origin or 127.0.0.1",
		)
	}
	if parsed.Port() == "" {
		return "", newLocalError("invalid_transport_origin", "Task Manager loopback transport requires an explicit port")
	}
	return parsed.String(), nil
}
