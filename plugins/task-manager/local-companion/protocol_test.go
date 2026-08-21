package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type fakeLocalUploader struct {
	result       uploadResult
	err          error
	input        localFileInput
	calls        int
	attachResult attachResult
	attachErr    error
	attachInput  attachLocalFileInput
	attachCalls  int
}

func (uploader *fakeLocalUploader) AttachLocalFileToTask(
	_ context.Context,
	input attachLocalFileInput,
) (attachResult, error) {
	uploader.attachCalls++
	uploader.attachInput = input
	return uploader.attachResult, uploader.attachErr
}

func (uploader *fakeLocalUploader) UploadLocalFile(_ context.Context, input localFileInput) (uploadResult, error) {
	uploader.calls++
	uploader.input = input
	return uploader.result, uploader.err
}

func TestMCPInitializeAndToolsListHaveNoUploadSideEffects(t *testing.T) {
	t.Parallel()
	uploader := &fakeLocalUploader{}
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"ping","params":{}}`,
	}, "\n") + "\n"
	var output bytes.Buffer
	if err := serveMCP(context.Background(), strings.NewReader(input), &output, uploader); err != nil {
		t.Fatal(err)
	}
	if uploader.calls != 0 || uploader.attachCalls != 0 {
		t.Fatalf("discovery invoked service: uploads=%d binds=%d", uploader.calls, uploader.attachCalls)
	}
	if strings.Contains(output.String(), "authorization") || strings.Contains(output.String(), "Keychain") {
		t.Fatalf("discovery returned an authentication side effect: %s", output.String())
	}
}

type cancellingUploader struct {
	canceled chan struct{}
}

func (uploader *cancellingUploader) AttachLocalFileToTask(
	context.Context,
	attachLocalFileInput,
) (attachResult, error) {
	return attachResult{}, newLocalError("not_available", "bind is not available")
}

func (uploader *cancellingUploader) UploadLocalFile(ctx context.Context, _ localFileInput) (uploadResult, error) {
	<-ctx.Done()
	close(uploader.canceled)
	return uploadResult{}, newLocalError(
		"network_error", "upload outcome is unknown; retry identical bytes with the same idempotencyKey",
	)
}

func TestMCPInitializeListAndUploadCall(t *testing.T) {
	t.Parallel()
	uploader := &fakeLocalUploader{result: uploadResult{
		FileRef: "file_protocol_123", Filename: "report.pdf", MediaType: "application/pdf",
		ByteSize: 123, ChecksumSHA256: strings.Repeat("a", 64), Kind: "file",
		State: "ready", Version: 1, SourceVerified: true,
	}, attachResult: attachResult{
		TaskRef: "TM-123", AttachmentRef: "attachment_protocol_123",
		Filename: "report.pdf", MediaType: "application/pdf", ByteSize: 123,
		ChecksumSHA256: strings.Repeat("a", 64), Kind: "file", State: "ready", Version: 1,
	}}
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"upload_local_file","arguments":{"path":"/tmp/report.pdf","idempotencyKey":"stable-key","expectedByteSize":123,"expectedSha256":"` + strings.Repeat("a", 64) + `"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"attach_local_file_to_task","arguments":{"taskRef":"TM-123","fileRef":"file_protocol_123","idempotencyKey":"stable-bind-key"}}}`,
	}, "\n") + "\n"
	var output bytes.Buffer
	if err := serveMCP(context.Background(), strings.NewReader(input), &output, uploader); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("responses = %d: %s", len(lines), output.String())
	}
	var initialize map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &initialize); err != nil {
		t.Fatal(err)
	}
	result := initialize["result"].(map[string]any)
	if result["protocolVersion"] != "2025-11-25" {
		t.Fatalf("initialize = %#v", result)
	}
	var list map[string]any
	_ = json.Unmarshal([]byte(lines[1]), &list)
	tools := list["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 2 || tools[0].(map[string]any)["name"] != "upload_local_file" ||
		tools[1].(map[string]any)["name"] != "attach_local_file_to_task" {
		t.Fatalf("tools = %#v", tools)
	}
	if uploader.input.Path != "/tmp/report.pdf" || uploader.input.IdempotencyKey != "stable-key" {
		t.Fatalf("call input = %#v", uploader.input)
	}
	callsByID := map[int]map[string]any{}
	for _, line := range lines[2:] {
		var call map[string]any
		_ = json.Unmarshal([]byte(line), &call)
		callsByID[int(call["id"].(float64))] = call
	}
	call := callsByID[3]
	callResult := call["result"].(map[string]any)
	if callResult["isError"] != false {
		t.Fatalf("call result = %#v", callResult)
	}
	structured := callResult["structuredContent"].(map[string]any)
	if structured["fileRef"] != "file_protocol_123" || structured["sourceVerified"] != true {
		t.Fatalf("structured result = %#v", structured)
	}
	if uploader.attachInput.TaskRef != "TM-123" || uploader.attachInput.FileRef != "file_protocol_123" ||
		uploader.attachInput.IdempotencyKey != "stable-bind-key" {
		t.Fatalf("bind input = %#v", uploader.attachInput)
	}
	bindCall := callsByID[4]
	bindStructured := bindCall["result"].(map[string]any)["structuredContent"].(map[string]any)
	if bindStructured["attachmentRef"] != "attachment_protocol_123" || bindStructured["taskRef"] != "TM-123" {
		t.Fatalf("bind result = %#v", bindStructured)
	}
}

func TestMCPToolErrorIsBoundedAndDoesNotEchoPath(t *testing.T) {
	t.Parallel()
	uploader := &fakeLocalUploader{err: newLocalError("local_path_denied", "host access to the authorized path was denied")}
	input := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"upload_local_file","arguments":{"path":"/Users/example/Secret/report.pdf","idempotencyKey":"stable"}}}` + "\n"
	var output bytes.Buffer
	if err := serveMCP(context.Background(), strings.NewReader(input), &output, uploader); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "/Users/example") || strings.Contains(output.String(), "Secret") {
		t.Fatalf("response echoed local path: %s", output.String())
	}
	var response map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &response); err != nil {
		t.Fatal(err)
	}
	result := response["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatalf("response = %#v", response)
	}
}

func TestMCPCancelNotificationCancelsUploadWithoutChangingRetryContract(t *testing.T) {
	t.Parallel()
	uploader := &cancellingUploader{canceled: make(chan struct{})}
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":"upload-1","method":"tools/call","params":{"name":"upload_local_file","arguments":{"path":"/tmp/report.pdf","idempotencyKey":"stable"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":"upload-1","reason":"test"}}`,
	}, "\n") + "\n"
	var output bytes.Buffer
	if err := serveMCP(context.Background(), strings.NewReader(input), &output, uploader); err != nil {
		t.Fatal(err)
	}
	<-uploader.canceled
	if !strings.Contains(output.String(), "same idempotencyKey") {
		t.Fatalf("cancel result lost retry contract: %s", output.String())
	}
}
