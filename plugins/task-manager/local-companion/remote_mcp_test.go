package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestRemoteMCPForwardsSessionAuthAndSSEWithoutLosingToolMetadata(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	methods := []string{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var rpc jsonRPCRequest
		if err := json.NewDecoder(request.Body).Decode(&rpc); err != nil {
			t.Errorf("request: %v", err)
			return
		}
		mutex.Lock()
		methods = append(methods, rpc.Method)
		mutex.Unlock()
		switch rpc.Method {
		case "initialize":
			writer.Header().Set("Content-Type", "application/json")
			writer.Header().Set("Mcp-Session-Id", "uat-session-1")
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"jsonrpc": "2.0", "id": json.RawMessage(rpc.ID),
				"result": map[string]any{"protocolVersion": mcpProtocolVersion, "capabilities": map[string]any{"tools": map[string]any{}}},
			})
		case "notifications/initialized":
			if request.Header.Get("Mcp-Session-Id") != "uat-session-1" {
				t.Errorf("session header = %q", request.Header.Get("Mcp-Session-Id"))
			}
			writer.WriteHeader(http.StatusAccepted)
		case "tools/list":
			if request.Header.Get("Mcp-Session-Id") != "uat-session-1" {
				t.Errorf("session header = %q", request.Header.Get("Mcp-Session-Id"))
			}
			writer.Header().Set("Content-Type", "text/event-stream")
			_, _ = writer.Write([]byte("event: message\ndata: {\"jsonrpc\":\"2.0\",\"id\":2,\"result\":{\"tools\":[{\"name\":\"upload_file\",\"_meta\":{\"openai/fileParams\":[\"file\"]}}]}}\n\n"))
		case "tools/call":
			if request.Header.Get("Authorization") != "Bearer tm_oat_test" {
				t.Errorf("authorization = %q", request.Header.Get("Authorization"))
			}
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"jsonrpc": "2.0", "id": json.RawMessage(rpc.ID),
				"result": map[string]any{"content": []any{}, "structuredContent": map[string]any{"ok": true}},
			})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client := newRemoteMCPClient(server.Client(), server.URL, &staticTokenSource{token: "tm_oat_test"})
	initialize, rpcErr := client.Forward(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "initialize", Params: json.RawMessage(`{"protocolVersion":"2025-11-25"}`),
	}, false)
	if rpcErr != nil || !strings.Contains(string(initialize), mcpProtocolVersion) {
		t.Fatalf("initialize = %s, error = %#v", initialize, rpcErr)
	}
	if _, rpcErr := client.Forward(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", Method: "notifications/initialized",
	}, false); rpcErr != nil {
		t.Fatal(rpcErr)
	}
	remoteList, rpcErr := client.Forward(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", ID: json.RawMessage("2"), Method: "tools/list", Params: json.RawMessage(`{}`),
	}, false)
	if rpcErr != nil {
		t.Fatal(rpcErr)
	}
	merged, rpcErr := mergeRemoteTools(remoteList)
	if rpcErr != nil {
		t.Fatal(rpcErr)
	}
	mergedJSON, _ := json.Marshal(merged)
	if !strings.Contains(string(mergedJSON), `"openai/fileParams":["file"]`) ||
		!strings.Contains(string(mergedJSON), `"name":"upload_local_file"`) {
		t.Fatalf("merged tools = %s", mergedJSON)
	}
	call, rpcErr := client.Forward(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", ID: json.RawMessage("3"), Method: "tools/call",
		Params: json.RawMessage(`{"name":"get_workspace","arguments":{}}`),
	}, true)
	if rpcErr != nil || !strings.Contains(string(call), `"ok":true`) {
		t.Fatalf("call = %s, error = %#v", call, rpcErr)
	}
	mutex.Lock()
	defer mutex.Unlock()
	if strings.Join(methods, ",") != "initialize,notifications/initialized,tools/list,tools/call" {
		t.Fatalf("methods = %#v", methods)
	}
}

func TestRemoteMCPRejectsHTMLAccessErrorWithoutEchoingBody(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte("<html>private edge detail</html>"))
	}))
	defer server.Close()
	client := newRemoteMCPClient(server.Client(), server.URL, &staticTokenSource{token: "tm_oat_test"})
	_, rpcErr := client.Forward(context.Background(), jsonRPCRequest{
		JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/call",
		Params: json.RawMessage(`{"name":"get_workspace","arguments":{}}`),
	}, true)
	if rpcErr == nil || !strings.Contains(rpcErr.Message, "before Task Manager") {
		t.Fatalf("error = %#v", rpcErr)
	}
	if strings.Contains(rpcErr.Message, "private edge detail") || strings.Contains(rpcErr.Message, "<html>") {
		t.Fatalf("error disclosed body: %q", rpcErr.Message)
	}
}
