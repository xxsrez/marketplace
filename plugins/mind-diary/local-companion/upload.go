package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	uploadIntentRoutePrefix = "/api/file-ingress/v1/upload-intents/"
	maxHostedResponseBytes  = int64(64 * 1024)
)

var localFileRefPattern = regexp.MustCompile(`^mdlocal_v1_[A-Za-z0-9_-]{16,256}$`)
var uploadCapabilityPattern = regexp.MustCompile(`^mdupload_v1_[A-Za-z0-9_-]+$`)

type hostedIntentStatus struct {
	Status     string
	ExpiresAt  string
	StagedFile *stagedFileReceipt
	Code       string
}

type localFileServiceImpl struct {
	store        *localFileStore
	httpClient   *http.Client
	publicOrigin string
}

func newLocalFileService(client *http.Client, origin string, workspaceRoots ...string) (*localFileServiceImpl, error) {
	if client == nil {
		return nil, errors.New("http client is required")
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
		parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("public origin is invalid")
	}
	ownedClient := *client
	ownedClient.Jar = nil
	ownedClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return errors.New("redirects are not accepted")
	}
	canonicalRoots, err := canonicalWorkspaceRoots(workspaceRoots)
	if err != nil {
		return nil, err
	}
	return &localFileServiceImpl{
		store:      newLocalFileStoreWithWorkspaceRoots(canonicalRoots),
		httpClient: &ownedClient, publicOrigin: parsed.Scheme + "://" + parsed.Host,
	}, nil
}

func (service *localFileServiceImpl) Close() error { return service.store.Close() }

func (service *localFileServiceImpl) PrepareLocalFile(
	_ context.Context,
	input prepareLocalFileInput,
) (prepareLocalFileResult, error) {
	return service.store.Prepare(input)
}

func (service *localFileServiceImpl) UploadPreparedFile(
	ctx context.Context,
	input uploadPreparedFileInput,
) (stagedFileReceipt, error) {
	prepared, err := service.store.acquire(input.LocalFileRef)
	if err != nil {
		return stagedFileReceipt{}, err
	}
	consume := false
	defer func() { service.store.release(input.LocalFileRef, consume) }()
	uploadURL, err := service.exactUploadURL(input.UploadURL)
	if err != nil {
		consume = true
		return stagedFileReceipt{}, err
	}

	status, err := service.readStatus(ctx, uploadURL)
	if err != nil {
		consume = !isRetryable(err)
		return stagedFileReceipt{}, err
	}
	if status.Status == "staged" {
		if err := verifyStagedReceipt(status.StagedFile, prepared.snapshot); err != nil {
			consume = true
			return stagedFileReceipt{}, err
		}
		consume = true
		return *status.StagedFile, nil
	}
	if status.Status == "rejected" {
		err := hostedRejection(status.Code)
		consume = !isRetryable(err)
		return stagedFileReceipt{}, err
	}

	before, err := prepared.file.Stat()
	if err != nil || !sameSnapshot(prepared.snapshot.info, before) {
		consume = true
		return stagedFileReceipt{}, newLocalError(
			"local_companion_file_changed", "the prepared local snapshot changed before upload",
		)
	}
	if _, err := prepared.file.Seek(0, io.SeekStart); err != nil {
		consume = true
		return stagedFileReceipt{}, newLocalError(
			"file_ingress_source_unavailable", "the prepared local snapshot cannot be streamed",
		)
	}
	stream := newStreamVerifier(io.LimitReader(prepared.file, prepared.snapshot.size+1))
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, stream)
	if err != nil {
		return stagedFileReceipt{}, newLocalError(
			"invalid_upload_url", "upload_url is invalid",
		)
	}
	request.ContentLength = prepared.snapshot.size
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("Accept", "application/json")
	response, requestErr := service.httpClient.Do(request)
	streamExact := stream.readSize == prepared.snapshot.size &&
		stream.digest() == prepared.snapshot.sha256
	after, statErr := prepared.file.Stat()
	fileStable := statErr == nil && sameSnapshot(prepared.snapshot.info, after)
	if requestErr != nil {
		reconciled, reconcileErr := service.reconcileUnknown(ctx, uploadURL, prepared.snapshot)
		if reconcileErr == nil && reconciled != nil {
			consume = true
			return *reconciled, nil
		}
		if !streamExact || !fileStable {
			consume = true
			return stagedFileReceipt{}, newLocalError(
				"local_companion_file_changed", "the prepared local snapshot changed during upload",
			)
		}
		return stagedFileReceipt{}, newLocalError(
			"file_ingress_transport_unavailable", "the hosted upload outcome is unknown", true,
		)
	}
	defer response.Body.Close()
	status, parseErr := decodeHostedResponse(response.Body)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseErr := hostedHTTPError(response.StatusCode, status, parseErr)
		if response.StatusCode >= 500 {
			reconciled, reconcileErr := service.reconcileUnknown(ctx, uploadURL, prepared.snapshot)
			if reconcileErr == nil && reconciled != nil {
				consume = true
				return *reconciled, nil
			}
		}
		consume = !isRetryable(responseErr)
		return stagedFileReceipt{}, responseErr
	}
	if parseErr != nil || status.Status != "staged" {
		reconciled, reconcileErr := service.reconcileUnknown(ctx, uploadURL, prepared.snapshot)
		if reconcileErr == nil && reconciled != nil {
			consume = true
			return *reconciled, nil
		}
		return stagedFileReceipt{}, newLocalError(
			"file_ingress_transport_unavailable", "the hosted upload response is invalid", true,
		)
	}
	if !streamExact || !fileStable {
		if err := verifyStagedReceipt(status.StagedFile, prepared.snapshot); err != nil {
			consume = true
			return stagedFileReceipt{}, newLocalError(
				"local_companion_file_changed", "the prepared local snapshot changed during upload",
			)
		}
	}
	if err := verifyStagedReceipt(status.StagedFile, prepared.snapshot); err != nil {
		consume = true
		return stagedFileReceipt{}, err
	}
	consume = true
	return *status.StagedFile, nil
}

func (service *localFileServiceImpl) exactUploadURL(value string) (string, error) {
	if value == "" || len(value) > 4_096 {
		return "", newLocalError("invalid_upload_url", "upload_url is invalid")
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme+"://"+parsed.Host != service.publicOrigin ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" ||
		!strings.HasPrefix(parsed.EscapedPath(), uploadIntentRoutePrefix) {
		return "", newLocalError("invalid_upload_url", "upload_url is not an exact trusted Mind Diary capability")
	}
	encodedCapability := strings.TrimPrefix(parsed.EscapedPath(), uploadIntentRoutePrefix)
	capability, err := url.PathUnescape(encodedCapability)
	if err != nil || encodedCapability != capability ||
		len(capability) < len("mdupload_v1_")+16 || len(capability) > 4096 ||
		strings.Contains(capability, "/") || !uploadCapabilityPattern.MatchString(capability) {
		return "", newLocalError("invalid_upload_url", "upload_url is not an exact trusted Mind Diary capability")
	}
	return parsed.String(), nil
}

func (service *localFileServiceImpl) readStatus(ctx context.Context, uploadURL string) (hostedIntentStatus, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uploadURL, nil)
	if err != nil {
		return hostedIntentStatus{}, newLocalError("invalid_upload_url", "upload_url is invalid")
	}
	request.Header.Set("Accept", "application/json")
	response, err := service.httpClient.Do(request)
	if err != nil {
		return hostedIntentStatus{}, newLocalError(
			"file_ingress_transport_unavailable", "the hosted upload status is unavailable", true,
		)
	}
	defer response.Body.Close()
	status, parseErr := decodeHostedResponse(response.Body)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return hostedIntentStatus{}, hostedHTTPError(response.StatusCode, status, parseErr)
	}
	if parseErr != nil || (status.Status != "pending" && status.Status != "staged" && status.Status != "rejected") {
		return hostedIntentStatus{}, newLocalError(
			"file_ingress_transport_unavailable", "the hosted upload status response is invalid", true,
		)
	}
	return status, nil
}

func (service *localFileServiceImpl) reconcileUnknown(
	ctx context.Context,
	uploadURL string,
	snapshot fileSnapshot,
) (*stagedFileReceipt, error) {
	status, err := service.readStatus(ctx, uploadURL)
	if err != nil {
		return nil, err
	}
	if status.Status == "rejected" {
		return nil, hostedRejection(status.Code)
	}
	if status.Status != "staged" {
		return nil, nil
	}
	if err := verifyStagedReceipt(status.StagedFile, snapshot); err != nil {
		return nil, err
	}
	return status.StagedFile, nil
}

func decodeHostedResponse(reader io.Reader) (hostedIntentStatus, error) {
	var envelope struct {
		OK   bool `json:"ok"`
		Data *struct {
			Status     string             `json:"status"`
			ExpiresAt  string             `json:"expires_at"`
			StagedFile *stagedFileReceipt `json:"staged_file"`
			Code       string             `json:"code"`
		} `json:"data"`
		Error *struct {
			Code      string `json:"code"`
			Retryable bool   `json:"retryable"`
		} `json:"error"`
	}
	var bounded bytes.Buffer
	read, copyErr := io.CopyN(&bounded, reader, maxHostedResponseBytes+1)
	if copyErr != nil && !errors.Is(copyErr, io.EOF) {
		return hostedIntentStatus{}, copyErr
	}
	if read > maxHostedResponseBytes {
		return hostedIntentStatus{}, errors.New("hosted response exceeds the bounded limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(bounded.Bytes()))
	if err := decoder.Decode(&envelope); err != nil {
		return hostedIntentStatus{}, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return hostedIntentStatus{}, errors.New("hosted response contains trailing data")
	}
	if envelope.OK && envelope.Data != nil {
		return hostedIntentStatus{
			Status: envelope.Data.Status, ExpiresAt: envelope.Data.ExpiresAt,
			StagedFile: envelope.Data.StagedFile, Code: envelope.Data.Code,
		}, nil
	}
	if envelope.Error != nil {
		return hostedIntentStatus{Code: envelope.Error.Code}, nil
	}
	return hostedIntentStatus{}, errors.New("invalid hosted response")
}

func hostedHTTPError(statusCode int, status hostedIntentStatus, parseErr error) error {
	code := status.Code
	if parseErr != nil {
		code = ""
	}
	switch {
	case statusCode == http.StatusNotFound || code == "file_ingress_source_unavailable":
		return newLocalError("file_ingress_source_unavailable", "the hosted upload intent is unavailable")
	case statusCode == http.StatusGone || code == "file_ingress_intent_expired":
		return newLocalError("file_ingress_intent_expired", "the hosted upload intent expired")
	case statusCode == http.StatusConflict || code == "file_ingress_intent_conflict":
		return newLocalError("file_ingress_intent_conflict", "the hosted upload intent cannot be reused")
	case code != "":
		return hostedRejection(code)
	default:
		return newLocalError(
			"file_ingress_transport_unavailable", "the hosted upload transport is unavailable",
			statusCode >= 500 || statusCode == http.StatusTooManyRequests,
		)
	}
}

func hostedRejection(code string) error {
	retryable := code == "capacity_soft_limit" || code == "capacity_accounting_untrusted" ||
		code == "file_ingress_transport_unavailable"
	for _, supported := range []string{
		"file_ingress_intent_conflict",
		"bundle_file_size_limit_exceeded",
		"bundle_file_size_mismatch",
		"bundle_file_digest_mismatch",
		"invalid_bundle_file_name",
		"staging_quota_exceeded",
		"capacity_soft_limit",
		"capacity_hard_limit",
		"capacity_fairness_limit",
		"capacity_accounting_untrusted",
		"file_ingress_transport_unavailable",
		"invalid_request",
	} {
		if code == supported {
			return newLocalError(code, "the hosted upload request was rejected", retryable)
		}
	}
	return newLocalError("file_ingress_intent_conflict", "the hosted upload request was rejected")
}

func isRetryable(err error) bool {
	var localErr *localError
	return errors.As(err, &localErr) && localErr.Retryable
}

func verifyStagedReceipt(receipt *stagedFileReceipt, snapshot fileSnapshot) error {
	_, expiryErr := time.Parse(time.RFC3339Nano, receiptExpiry(receipt))
	if receipt == nil || receipt.State != "verified" ||
		receipt.StagedFileRef == "" || len(receipt.StagedFileRef) > 512 || containsControl(receipt.StagedFileRef) ||
		receipt.SourceKind != snapshot.sourceKind ||
		receipt.DisplayFilename != snapshot.displayFilename ||
		!mimePattern.MatchString(receipt.MediaType) || len(receipt.MediaType) > 127 ||
		receipt.SHA256 != snapshot.sha256 || receipt.Size != snapshot.size ||
		receipt.ExpiresAt == "" || expiryErr != nil {
		return newLocalError(
			"file_ingress_transport_unavailable", "the hosted staged-file receipt does not match the prepared snapshot",
		)
	}
	return nil
}

func receiptExpiry(receipt *stagedFileReceipt) string {
	if receipt == nil {
		return ""
	}
	return receipt.ExpiresAt
}

func (status hostedIntentStatus) String() string {
	return fmt.Sprintf("hosted intent status %q", status.Status)
}
