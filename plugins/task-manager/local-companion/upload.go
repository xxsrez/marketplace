package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type tokenSource interface {
	Token(context.Context) (string, error)
	Invalidate()
}

type uploadClient struct {
	httpClient *http.Client
	origin     string
	tokens     tokenSource
}

type uploadResult struct {
	FileRef        string `json:"fileRef"`
	Filename       string `json:"filename"`
	MediaType      string `json:"mediaType"`
	ByteSize       int64  `json:"byteSize"`
	ChecksumSHA256 string `json:"checksumSha256"`
	Kind           string `json:"kind"`
	State          string `json:"state"`
	ImageWidth     *int64 `json:"imageWidth,omitempty"`
	ImageHeight    *int64 `json:"imageHeight,omitempty"`
	ReadyExpiresAt string `json:"readyExpiresAt,omitempty"`
	Version        int64  `json:"version"`
	SourceVerified bool   `json:"sourceVerified"`
}

type agentStoredFile struct {
	Ref            string `json:"ref"`
	Filename       string `json:"filename"`
	MediaType      string `json:"mediaType"`
	ByteSize       int64  `json:"byteSize"`
	ChecksumSHA256 string `json:"checksumSha256"`
	Kind           string `json:"kind"`
	State          string `json:"state"`
	ImageWidth     *int64 `json:"imageWidth,omitempty"`
	ImageHeight    *int64 `json:"imageHeight,omitempty"`
	ReadyExpiresAt string `json:"readyExpiresAt,omitempty"`
	Version        int64  `json:"version"`
}

func newUploadClient(httpClient *http.Client, origin string, tokens tokenSource) *uploadClient {
	return &uploadClient{httpClient: httpClient, origin: strings.TrimRight(origin, "/"), tokens: tokens}
}

func (client *uploadClient) UploadLocalFile(ctx context.Context, input localFileInput) (uploadResult, error) {
	prepared, err := prepareLocalFile(input, fileReadHooks{})
	if err != nil {
		return uploadResult{}, err
	}
	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	for attempt := 0; attempt < 2; attempt++ {
		token, err := client.tokens.Token(ctx)
		if err != nil {
			return uploadResult{}, err
		}
		result, status, err := client.upload(ctx, token, idempotencyKey, prepared)
		if status == http.StatusUnauthorized && attempt == 0 {
			client.tokens.Invalidate()
			continue
		}
		return result, err
	}
	return uploadResult{}, newLocalError("remote_unauthenticated", "Task Manager authorization is unavailable")
}

func (client *uploadClient) upload(
	ctx context.Context,
	token string,
	idempotencyKey string,
	prepared preparedLocalFile,
) (uploadResult, int, error) {
	request, err := http.NewRequestWithContext(
		ctx, http.MethodPost, client.origin+"/api/agent/v1/files", bytes.NewReader(prepared.Bytes),
	)
	if err != nil {
		return uploadResult{}, 0, newLocalError("remote_request_failed", "upload request could not be created")
	}
	request.ContentLength = prepared.ByteSize
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", prepared.MediaType)
	request.Header.Set("Idempotency-Key", idempotencyKey)
	request.Header.Set("X-File-Filename", url.PathEscape(prepared.Filename))
	request.Header.Set("Accept", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return uploadResult{}, 0, newLocalError(
			"network_error", "upload did not complete; retry the identical file with the same idempotencyKey",
		)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024+1))
	if err != nil || len(body) > 1024*1024 {
		return uploadResult{}, response.StatusCode, newLocalError(
			"remote_response_invalid", "Task Manager returned an unreadable response",
		)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return uploadResult{}, response.StatusCode, parseAgentAPIError(response.StatusCode, body)
	}
	var envelope struct {
		Data agentStoredFile `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return uploadResult{}, response.StatusCode, newLocalError(
			"remote_response_invalid", "Task Manager returned invalid file metadata",
		)
	}
	data := envelope.Data
	if data.Ref == "" || data.State != "ready" || data.Filename != prepared.Filename ||
		data.MediaType != prepared.MediaType || data.ByteSize != prepared.ByteSize ||
		!strings.EqualFold(data.ChecksumSHA256, prepared.ChecksumSHA256) {
		return uploadResult{}, response.StatusCode, newLocalError(
			"server_metadata_mismatch", "Task Manager file metadata does not match the stable local snapshot",
		)
	}
	result := uploadResult{
		FileRef: data.Ref, Filename: data.Filename, MediaType: data.MediaType,
		ByteSize: data.ByteSize, ChecksumSHA256: strings.ToLower(data.ChecksumSHA256),
		Kind: data.Kind, State: data.State, ImageWidth: data.ImageWidth,
		ImageHeight: data.ImageHeight, ReadyExpiresAt: data.ReadyExpiresAt,
		Version: data.Version, SourceVerified: true,
	}
	return result, response.StatusCode, nil
}

func parseAgentAPIError(status int, body []byte) error {
	var envelope struct {
		Error struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &envelope)
	code := envelope.Error.Code
	if code == "" {
		code = fmt.Sprintf("http_%d", status)
	}
	message := envelope.Error.Message
	if message == "" {
		message = "Task Manager rejected the upload"
	}
	if envelope.Error.RequestID != "" {
		message += " (requestId " + envelope.Error.RequestID + ")"
	}
	return newLocalError("remote_"+code, message)
}
