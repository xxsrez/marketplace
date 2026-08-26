package main

// canonicalRuntimeErrorCodes mirrors the accepted MindDiary LocalFileCompanion
// contract. The stdio adapter must never invent a second error vocabulary.
var canonicalRuntimeErrorCodes = map[string]struct{}{
	"invalid_path":                        {},
	"file_ingress_source_unavailable":     {},
	"file_ingress_source_unsupported":     {},
	"file_ingress_transport_unavailable":  {},
	"file_ingress_intent_expired":         {},
	"file_ingress_intent_conflict":        {},
	"bundle_file_size_limit_exceeded":     {},
	"bundle_file_size_mismatch":           {},
	"bundle_file_digest_mismatch":         {},
	"invalid_bundle_file_name":            {},
	"staging_quota_exceeded":              {},
	"capacity_soft_limit":                 {},
	"capacity_hard_limit":                 {},
	"capacity_fairness_limit":             {},
	"capacity_accounting_untrusted":       {},
	"local_companion_file_changed":        {},
	"local_companion_cancelled":           {},
	"local_companion_concurrency_limit":   {},
	"local_companion_invalid_source_kind": {},
	"local_companion_ref_not_found":       {},
	"local_companion_ref_expired":         {},
	"local_companion_ref_in_use":          {},
	"invalid_request":                     {},
	"invalid_upload_url":                  {},
}

func canonicalRuntimeErrorCode(code string) string {
	if _, exists := canonicalRuntimeErrorCodes[code]; exists {
		return code
	}
	return "file_ingress_transport_unavailable"
}
