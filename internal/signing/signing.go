package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/google/uuid"
)

// Signer generates HMAC-SHA256 request signatures matching the Idealista mobile
// app's signing scheme: seed + METHOD + sortedEncode(queryParams) + sortedEncode(bodyParams).
type Signer struct {
	secret string
}

// New creates a Signer using the provided raw signing secret. The raw secret is
// base64-encoded before use as the HMAC key.
func New(raw string) (*Signer, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("signing secret is required")
	}
	return &Signer{secret: base64.StdEncoding.EncodeToString([]byte(raw))}, nil
}

// Sign generates a seed UUID and computes an HMAC-SHA256 signature for the given
// HTTP method and parameters. Returns the seed and hex-encoded signature.
//
// The signature input is: seed + METHOD + sortedEncode(queryParams) + sortedEncode(bodyParams)
// This matches the JavaScript signing logic used in the Idealista mobile app.
func (s *Signer) Sign(method string, queryParams, bodyParams url.Values) (seed, signature string, err error) {
	seed = uuid.New().String()
	input := seed + strings.ToUpper(method) + sortedEncode(queryParams) + sortedEncode(bodyParams)
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(input))
	signature = hex.EncodeToString(mac.Sum(nil))
	return seed, signature, nil
}

// sortedEncode encodes url.Values into a sorted query string using JS encodeURIComponent
// semantics: keys are sorted lexicographically, spaces encoded as %20 (not +).
func sortedEncode(params url.Values) string {
	if len(params) == 0 {
		return ""
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(params))
	for _, k := range keys {
		for _, v := range params[k] {
			parts = append(parts, jsEncodeURIComponent(k)+"="+jsEncodeURIComponent(v))
		}
	}
	return strings.Join(parts, "&")
}

// jsEncodeURIComponent encodes a string matching JavaScript's encodeURIComponent:
// unreserved chars A-Z a-z 0-9 - _ . ! ~ * ' ( ) are not encoded; spaces → %20.
func jsEncodeURIComponent(s string) string {
	encoded := url.QueryEscape(s)
	// url.QueryEscape uses + for spaces; JS encodeURIComponent uses %20.
	return strings.ReplaceAll(encoded, "+", "%20")
}
