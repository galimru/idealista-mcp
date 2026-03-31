package debug

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// sensitiveHeaders are redacted from debug output to prevent credential leakage.
var sensitiveHeaders = map[string]bool{
	"Authorization": true,
	"Signature":     true,
}

// NewTransport wraps rt with request/response logging to stderr.
// Authorization and Signature headers are redacted.
func NewTransport(rt http.RoundTripper) http.RoundTripper {
	return &loggingTransport{rt}
}

type loggingTransport struct{ rt http.RoundTripper }

func (d *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(strings.NewReader(string(reqBody)))
	}

	fmt.Fprintf(os.Stderr, "\n> %s %s\n", req.Method, req.URL)
	for k, vs := range req.Header {
		if sensitiveHeaders[k] {
			fmt.Fprintf(os.Stderr, "> %s: [REDACTED]\n", k)
			continue
		}
		fmt.Fprintf(os.Stderr, "> %s: %s\n", k, strings.Join(vs, ", "))
	}
	if len(reqBody) > 0 {
		fmt.Fprintf(os.Stderr, ">\n> %s\n", reqBody)
	}

	resp, err := d.rt.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "* error: %v\n", err)
		return resp, err
	}

	respBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(strings.NewReader(string(respBody)))

	fmt.Fprintf(os.Stderr, "\n< %s\n", resp.Status)
	for k, vs := range resp.Header {
		fmt.Fprintf(os.Stderr, "< %s: %s\n", k, strings.Join(vs, ", "))
	}
	if len(respBody) > 0 {
		fmt.Fprintf(os.Stderr, "<\n< %s\n", respBody)
	}
	return resp, nil
}
