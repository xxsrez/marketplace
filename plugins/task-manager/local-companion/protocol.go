package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
)

const mcpProtocolVersion = "2025-11-25"

type localUploader interface {
	UploadLocalFile(context.Context, localFileInput) (uploadResult, error)
}

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func serveMCP(ctx context.Context, input io.Reader, output io.Writer, uploader localUploader) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	var outputMutex sync.Mutex
	var writeError error
	writeResponse := func(response jsonRPCResponse) {
		outputMutex.Lock()
		defer outputMutex.Unlock()
		if writeError == nil {
			writeError = encoder.Encode(response)
		}
	}
	activeCalls := map[string]context.CancelFunc{}
	var activeMutex sync.Mutex
	var calls sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var request jsonRPCRequest
		if err := json.Unmarshal(line, &request); err != nil || request.JSONRPC != "2.0" || request.Method == "" {
			writeResponse(jsonRPCResponse{
				JSONRPC: "2.0", ID: json.RawMessage("null"),
				Error: &jsonRPCError{Code: -32700, Message: "Parse error"},
			})
			continue
		}
		if request.Method == "notifications/cancelled" {
			var params struct {
				RequestID json.RawMessage `json:"requestId"`
			}
			if json.Unmarshal(request.Params, &params) == nil && len(params.RequestID) > 0 {
				activeMutex.Lock()
				cancel := activeCalls[string(params.RequestID)]
				activeMutex.Unlock()
				if cancel != nil {
					cancel()
				}
			}
			continue
		}
		if request.Method == "tools/call" && len(request.ID) > 0 {
			callCtx, cancel := context.WithCancel(ctx)
			key := string(request.ID)
			activeMutex.Lock()
			activeCalls[key] = cancel
			activeMutex.Unlock()
			calls.Add(1)
			go func(request jsonRPCRequest) {
				defer calls.Done()
				defer cancel()
				defer func() {
					activeMutex.Lock()
					delete(activeCalls, key)
					activeMutex.Unlock()
				}()
				result, protocolErr := handleMCPRequest(callCtx, request, uploader)
				response := jsonRPCResponse{JSONRPC: "2.0", ID: request.ID, Result: result}
				if protocolErr != nil {
					response.Result = nil
					response.Error = protocolErr
				}
				writeResponse(response)
			}(request)
			continue
		}
		result, protocolErr := handleMCPRequest(ctx, request, uploader)
		if len(request.ID) == 0 {
			continue
		}
		response := jsonRPCResponse{JSONRPC: "2.0", ID: request.ID, Result: result}
		if protocolErr != nil {
			response.Result = nil
			response.Error = protocolErr
		}
		writeResponse(response)
	}
	scanError := scanner.Err()
	if scanError != nil {
		activeMutex.Lock()
		for _, cancel := range activeCalls {
			cancel()
		}
		activeMutex.Unlock()
	}
	calls.Wait()
	if scanError != nil {
		return scanError
	}
	return writeError
}

func handleMCPRequest(
	ctx context.Context,
	request jsonRPCRequest,
	uploader localUploader,
) (any, *jsonRPCError) {
	switch request.Method {
	case "initialize":
		var params struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(request.Params, &params)
		protocolVersion := mcpProtocolVersion
		if params.ProtocolVersion == "2025-11-25" || params.ProtocolVersion == "2024-11-05" {
			protocolVersion = params.ProtocolVersion
		}
		return map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
			"serverInfo": map[string]any{
				"name": "task-manager-local", "title": "Task Manager Local Files", "version": "0.1.0",
			},
		}, nil
	case "notifications/initialized", "notifications/cancelled":
		return nil, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": []any{uploadLocalFileTool()}}, nil
	case "tools/call":
		return handleToolCall(ctx, request.Params, uploader)
	default:
		return nil, &jsonRPCError{Code: -32601, Message: "Method not found"}
	}
}

func uploadLocalFileTool() map[string]any {
	return map[string]any{
		"name":        "upload_local_file",
		"title":       "Upload a local file",
		"description": "Read one exact user-authorized absolute filesystem path through the Codex host, verify a stable regular-file snapshot, and upload it into Task Manager as an unbound StoredFile. Returns fileRef for attach_file_to_task. Never accepts directories, globs, final-component symlinks, devices, inline base64, or remote URLs; never sends the full path to Task Manager.",
		"inputSchema": map[string]any{
			"type": "object", "additionalProperties": false,
			"required": []string{"path", "idempotencyKey"},
			"properties": map[string]any{
				"path": map[string]any{
					"type": "string", "minLength": 1,
					"description": "Exact absolute path to one user-authorized regular file. Host sandbox and approval policy still apply.",
				},
				"idempotencyKey": map[string]any{
					"type": "string", "minLength": 1, "maxLength": 200,
					"description": "Stable retry key; reuse only for identical bytes, filename, and MIME metadata.",
				},
				"expectedByteSize": map[string]any{
					"type": "integer", "minimum": 0,
					"description": "Optional expected byte size checked before any network request.",
				},
				"expectedSha256": map[string]any{
					"type": "string", "pattern": "^[0-9a-fA-F]{64}$",
					"description": "Optional expected SHA-256 checked before any network request.",
				},
				"displayFilename": map[string]any{
					"type": "string", "minLength": 1, "maxLength": 255,
					"description": "Optional basename sent to Task Manager instead of the local basename; path separators are rejected.",
				},
			},
		},
		"annotations": map[string]any{
			"readOnlyHint": false, "destructiveHint": false,
			"idempotentHint": true, "openWorldHint": true,
		},
	}
}

func handleToolCall(ctx context.Context, rawParams json.RawMessage, uploader localUploader) (any, *jsonRPCError) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(rawParams, &params); err != nil || params.Name == "" {
		return nil, &jsonRPCError{Code: -32602, Message: "Invalid params"}
	}
	if params.Name != "upload_local_file" {
		return nil, &jsonRPCError{Code: -32602, Message: "Unknown tool"}
	}
	var input localFileInput
	decoderErr := json.Unmarshal(params.Arguments, &input)
	if decoderErr != nil {
		return toolError("invalid_arguments", "tool arguments are invalid"), nil
	}
	result, err := uploader.UploadLocalFile(ctx, input)
	if err != nil {
		var localErr *localError
		if errors.As(err, &localErr) {
			return toolError(localErr.Code, localErr.Message), nil
		}
		return toolError("internal_error", "local upload companion could not complete the request"), nil
	}
	textBytes, _ := json.Marshal(result)
	return map[string]any{
		"content":           []any{map[string]any{"type": "text", "text": string(textBytes)}},
		"structuredContent": result,
		"isError":           false,
	}, nil
}

func toolError(code, message string) map[string]any {
	return map[string]any{
		"content": []any{map[string]any{
			"type": "text", "text": fmt.Sprintf("%s: %s", code, message),
		}},
		"isError": true,
	}
}
