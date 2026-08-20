package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const maxRemoteMCPResponseBytes = 4 * 1024 * 1024

type remoteMCPClient struct {
	httpClient *http.Client
	endpoint   string
	tokens     tokenSource
	mutex      sync.RWMutex
	sessionID  string
}

type remoteJSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCError   `json:"error"`
}

func newRemoteMCPClient(httpClient *http.Client, origin string, tokens tokenSource) *remoteMCPClient {
	return &remoteMCPClient{
		httpClient: httpClient,
		endpoint:   strings.TrimRight(origin, "/") + "/api/mcp",
		tokens:     tokens,
	}
}

func (client *remoteMCPClient) Forward(
	ctx context.Context,
	request jsonRPCRequest,
	requireAuth bool,
) (json.RawMessage, *jsonRPCError) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, remoteProtocolError("Remote MCP request could not be encoded")
	}
	for attempt := 0; attempt < 2; attempt++ {
		token := ""
		if requireAuth {
			token, err = client.tokens.Token(ctx)
			if err != nil {
				return nil, remoteProtocolError(boundedLocalMessage(err, "Task Manager authorization is unavailable"))
			}
		}
		response, body, contentType, requestErr := client.post(ctx, payload, token)
		if requestErr != nil {
			return nil, remoteProtocolError("Private UAT MCP is unavailable")
		}
		if response.StatusCode == http.StatusAccepted && len(body) == 0 {
			return nil, nil
		}
		if response.StatusCode == http.StatusUnauthorized && requireAuth && attempt == 0 &&
			isStructuredMCPContent(contentType) {
			client.tokens.Invalidate()
			continue
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			if !isStructuredMCPContent(contentType) {
				return nil, remoteProtocolError(fmt.Sprintf(
					"Private UAT access was rejected before Task Manager (HTTP %d)", response.StatusCode,
				))
			}
			return nil, remoteProtocolError(fmt.Sprintf("Remote MCP rejected the request (HTTP %d)", response.StatusCode))
		}
		decoded, decodeErr := decodeRemoteMCPResponse(body, contentType)
		if decodeErr != nil {
			return nil, remoteProtocolError("Remote MCP returned an invalid response")
		}
		if decoded.Error != nil {
			return nil, decoded.Error
		}
		return decoded.Result, nil
	}
	return nil, remoteProtocolError("Task Manager authorization is unavailable")
}

func (client *remoteMCPClient) post(
	ctx context.Context,
	payload []byte,
	token string,
) (*http.Response, []byte, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, nil, "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("MCP-Protocol-Version", mcpProtocolVersion)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	client.mutex.RLock()
	sessionID := client.sessionID
	client.mutex.RUnlock()
	if sessionID != "" {
		request.Header.Set("Mcp-Session-Id", sessionID)
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, nil, "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxRemoteMCPResponseBytes+1))
	if err != nil || len(body) > maxRemoteMCPResponseBytes {
		return response, nil, response.Header.Get("Content-Type"), errorsForRemoteBody()
	}
	if returnedSessionID := strings.TrimSpace(response.Header.Get("Mcp-Session-Id")); returnedSessionID != "" {
		client.mutex.Lock()
		client.sessionID = returnedSessionID
		client.mutex.Unlock()
	}
	return response, body, response.Header.Get("Content-Type"), nil
}

func decodeRemoteMCPResponse(body []byte, contentType string) (remoteJSONRPCResponse, error) {
	var response remoteJSONRPCResponse
	if strings.Contains(strings.ToLower(contentType), "text/event-stream") {
		scanner := bufio.NewScanner(bytes.NewReader(body))
		scanner.Buffer(make([]byte, 64*1024), maxRemoteMCPResponseBytes)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			candidate := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if candidate == "" || candidate == "[DONE]" {
				continue
			}
			if json.Unmarshal([]byte(candidate), &response) == nil && response.JSONRPC == "2.0" {
				return response, nil
			}
		}
		if err := scanner.Err(); err != nil {
			return remoteJSONRPCResponse{}, err
		}
		return remoteJSONRPCResponse{}, fmt.Errorf("SSE response contained no JSON-RPC message")
	}
	if err := json.Unmarshal(body, &response); err != nil || response.JSONRPC != "2.0" {
		return remoteJSONRPCResponse{}, fmt.Errorf("invalid JSON-RPC response")
	}
	return response, nil
}

func isStructuredMCPContent(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/event-stream")
}

func remoteProtocolError(message string) *jsonRPCError {
	return &jsonRPCError{Code: -32000, Message: message}
}

func boundedLocalMessage(err error, fallback string) string {
	var localErr *localError
	if errors.As(err, &localErr) {
		return localErr.Code + ": " + localErr.Message
	}
	return fallback
}

func errorsForRemoteBody() error {
	return fmt.Errorf("remote MCP response body exceeded the safe bound")
}
