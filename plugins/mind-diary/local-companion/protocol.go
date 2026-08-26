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

const mcpProtocolVersion = "2026-07-28"

type localFileService interface {
	PrepareLocalFile(context.Context, prepareLocalFileInput) (prepareLocalFileResult, error)
	UploadPreparedFile(context.Context, uploadPreparedFileInput) (stagedFileReceipt, error)
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

func serveMCP(ctx context.Context, input io.Reader, output io.Writer, service localFileService) error {
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
				result, protocolErr := handleMCPRequest(callCtx, request, service)
				response := jsonRPCResponse{JSONRPC: "2.0", ID: request.ID, Result: result}
				if protocolErr != nil {
					response.Result = nil
					response.Error = protocolErr
				}
				writeResponse(response)
			}(request)
			continue
		}
		if request.Method == "tools/call" {
			continue
		}
		result, protocolErr := handleMCPRequest(ctx, request, service)
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
	if scanError := scanner.Err(); scanError != nil {
		activeMutex.Lock()
		for _, cancel := range activeCalls {
			cancel()
		}
		activeMutex.Unlock()
		calls.Wait()
		return scanError
	}
	calls.Wait()
	return writeError
}

func handleMCPRequest(
	ctx context.Context,
	request jsonRPCRequest,
	service localFileService,
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
				"name": "mind-diary-local", "title": "Mind Diary Local Files", "version": "0.1.0",
			},
		}, nil
	case "notifications/initialized", "notifications/cancelled":
		return nil, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": []any{prepareLocalFileTool(), uploadPreparedFileTool()}}, nil
	case "tools/call":
		return handleToolCall(ctx, request, service)
	default:
		return nil, &jsonRPCError{Code: -32601, Message: "Method not found"}
	}
}

func prepareLocalFileTool() map[string]any {
	return map[string]any{
		"name":        "prepare_local_file",
		"title":       "Prepare one local file for Mind Diary",
		"description": "Open one exact user-authorized absolute path as a stable regular-file snapshot, stream-hash at most 256 MiB, and return only a short-lived process-local reference plus path-free metadata for create_file_upload_intent. Directories, globs, traversal, final symlinks, special files and inline bytes are rejected.",
		"inputSchema": map[string]any{
			"type": "object", "additionalProperties": false,
			"required": []string{"path"},
			"properties": map[string]any{
				"path": map[string]any{
					"type": "string", "minLength": 1, "maxLength": 16384,
					"description": "One exact absolute path authorized by the Codex host. It is retained only inside this process and never returned or sent to Mind Diary.",
				},
				"source_kind": map[string]any{
					"type": "string", "enum": []string{"local_path", "workspace/generated_artifact"},
					"description": "Explicit provenance class; defaults to local_path. workspace/generated_artifact requires a canonical path inside trusted process-configured workspace roots.",
				},
				"display_filename": map[string]any{
					"type": "string", "minLength": 1, "maxLength": 255,
				},
				"claimed_media_type": map[string]any{
					"type": "string", "maxLength": 256,
					"description": "Optional advisory MIME. Invalid or unsafe evidence becomes application/octet-stream.",
				},
				"expected_size": map[string]any{
					"type": "integer", "minimum": 0, "maximum": maxLocalFileBytes,
				},
				"expected_sha256": map[string]any{
					"type": "string", "pattern": "^sha256:[0-9a-f]{64}$",
				},
			},
		},
		"outputSchema": prepareLocalFileOutputSchema(),
		"annotations": map[string]any{
			"readOnlyHint": true, "destructiveHint": false,
			"idempotentHint": false, "openWorldHint": false,
		},
	}
}

func uploadPreparedFileTool() map[string]any {
	return map[string]any{
		"name":        "upload_prepared_file",
		"title":       "Upload one prepared file to Mind Diary",
		"description": "Stream the exact process-local snapshot to one same-origin one-use upload_url returned by hosted create_file_upload_intent. Performs credentialless GET-before-PUT and reconciles unknown outcomes. Never accepts a path, bearer, cookie, arbitrary URL, base64 or directory.",
		"inputSchema": map[string]any{
			"type": "object", "additionalProperties": false,
			"required": []string{"local_file_ref", "upload_url"},
			"properties": map[string]any{
				"local_file_ref": map[string]any{
					"type": "string", "pattern": "^mdlocal_v1_[A-Za-z0-9_-]{16,256}$",
				},
				"upload_url": map[string]any{
					"type": "string", "format": "uri", "maxLength": 4096,
					"description": "Exact UAT capability URL returned by hosted create_file_upload_intent for the prepared metadata.",
				},
			},
		},
		"outputSchema": stagedFileOutputSchema(),
		"annotations": map[string]any{
			"readOnlyHint": false, "destructiveHint": false,
			"idempotentHint": true, "openWorldHint": true,
		},
	}
}

func prepareLocalFileOutputSchema() map[string]any {
	return map[string]any{
		"type": "object", "additionalProperties": false,
		"required": []string{
			"local_file_ref", "source_kind", "display_filename", "claimed_media_type",
			"expected_size", "expected_sha256", "expires_at",
		},
		"properties": map[string]any{
			"local_file_ref":     map[string]any{"type": "string", "pattern": "^mdlocal_v1_[A-Za-z0-9_-]{16,256}$"},
			"source_kind":        map[string]any{"type": "string", "enum": []string{"local_path", "workspace/generated_artifact"}},
			"display_filename":   map[string]any{"type": "string", "minLength": 1, "maxLength": 255},
			"claimed_media_type": map[string]any{"type": "string", "minLength": 3, "maxLength": 127},
			"expected_size":      map[string]any{"type": "integer", "minimum": 0, "maximum": maxLocalFileBytes},
			"expected_sha256":    map[string]any{"type": "string", "pattern": "^sha256:[0-9a-f]{64}$"},
			"expires_at":         map[string]any{"type": "string", "format": "date-time"},
		},
	}
}

func stagedFileOutputSchema() map[string]any {
	return map[string]any{
		"type": "object", "additionalProperties": false,
		"required": []string{
			"staged_file_ref", "state", "source_kind", "display_filename",
			"media_type", "sha256", "size", "expires_at", "replayed",
		},
		"properties": map[string]any{
			"staged_file_ref":  map[string]any{"type": "string", "minLength": 1, "maxLength": 512},
			"state":            map[string]any{"const": "verified"},
			"source_kind":      map[string]any{"type": "string", "enum": []string{"local_path", "workspace/generated_artifact"}},
			"display_filename": map[string]any{"type": "string", "minLength": 1, "maxLength": 255},
			"media_type":       map[string]any{"type": "string", "minLength": 3, "maxLength": 127},
			"sha256":           map[string]any{"type": "string", "pattern": "^sha256:[0-9a-f]{64}$"},
			"size":             map[string]any{"type": "integer", "minimum": 0, "maximum": maxLocalFileBytes},
			"expires_at":       map[string]any{"type": "string", "format": "date-time"},
			"replayed":         map[string]any{"type": "boolean"},
		},
	}
}

func handleToolCall(
	ctx context.Context,
	request jsonRPCRequest,
	service localFileService,
) (any, *jsonRPCError) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(request.Params, &params); err != nil || params.Name == "" {
		return nil, &jsonRPCError{Code: -32602, Message: "Invalid params"}
	}
	var result any
	var err error
	switch params.Name {
	case "prepare_local_file":
		var input prepareLocalFileInput
		if err := decodeExactArguments(params.Arguments, &input, []string{
			"path", "source_kind", "display_filename", "claimed_media_type", "expected_size", "expected_sha256",
		}); err != nil {
			return toolError("invalid_request", "tool arguments are invalid", false), nil
		}
		result, err = service.PrepareLocalFile(ctx, input)
	case "upload_prepared_file":
		var input uploadPreparedFileInput
		if err := decodeExactArguments(params.Arguments, &input, []string{
			"local_file_ref", "upload_url",
		}); err != nil {
			return toolError("invalid_request", "tool arguments are invalid", false), nil
		}
		result, err = service.UploadPreparedFile(ctx, input)
	default:
		return nil, &jsonRPCError{Code: -32602, Message: "Unknown tool"}
	}
	if err != nil {
		var localErr *localError
		if errors.As(err, &localErr) {
			return toolError(localErr.Code, localErr.Message, localErr.Retryable), nil
		}
		return toolError("file_ingress_transport_unavailable", "the local file companion could not complete the request", false), nil
	}
	textBytes, _ := json.Marshal(result)
	return map[string]any{
		"content":           []any{map[string]any{"type": "text", "text": string(textBytes)}},
		"structuredContent": result,
		"isError":           false,
	}, nil
}

func decodeExactArguments(raw json.RawMessage, destination any, allowed []string) error {
	if len(raw) == 0 {
		raw = json.RawMessage("{}")
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil || keys == nil {
		return errors.New("arguments must be an object")
	}
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, key := range allowed {
		allowedSet[key] = struct{}{}
	}
	for key := range keys {
		if _, ok := allowedSet[key]; !ok {
			return fmt.Errorf("unexpected argument %q", key)
		}
	}
	return json.Unmarshal(raw, destination)
}

func toolError(code, message string, retryable bool) map[string]any {
	code = canonicalRuntimeErrorCode(code)
	return map[string]any{
		"content": []any{map[string]any{
			"type": "text", "text": fmt.Sprintf("%s: %s", code, message),
		}},
		"structuredContent": map[string]any{
			"error": map[string]any{"code": code, "retryable": retryable},
		},
		"isError": true,
	}
}
