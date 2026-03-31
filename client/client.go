package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/galimru/idealista-mcp/auth"
	"github.com/galimru/idealista-mcp/internal/config"
	"github.com/galimru/idealista-mcp/internal/debug"
	"github.com/galimru/idealista-mcp/internal/signing"
)

// APIClient is the interface that tools use to interact with the Idealista API.
type APIClient interface {
	Get(ctx context.Context, rawURL string, queryParams url.Values) ([]byte, error)
	Post(ctx context.Context, rawURL string, bodyParams url.Values) ([]byte, error)
}

// APIError is returned when the server responds with a non-2xx status.
type APIError struct {
	URL    string
	Status int
	Body   []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API %s returned %d: %s", e.URL, e.Status, e.Body)
}

// IsNotFound reports whether err is an API 404 response.
func IsNotFound(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Status == http.StatusNotFound
}

// Client makes signed, authenticated requests to the Idealista API.
type Client struct {
	authClient       *auth.Client
	signer           *signing.Signer
	httpClient       *http.Client
	deviceIdentifier string
}

// New creates an API client backed by the given auth.Client and signing.Signer.
func New(a *auth.Client, s *signing.Signer, deviceIdentifier string) *Client {
	var transport http.RoundTripper = http.DefaultTransport
	if os.Getenv("IDEALISTA_DEBUG") != "" {
		transport = debug.NewTransport(transport)
	}
	return &Client{
		authClient:       a,
		signer:           s,
		httpClient:       &http.Client{Timeout: 30 * time.Second, Transport: transport},
		deviceIdentifier: deviceIdentifier,
	}
}

// Get performs a signed, authenticated GET request and returns the response body.
func (c *Client) Get(ctx context.Context, rawURL string, queryParams url.Values) ([]byte, error) {
	body, status, err := c.doGet(ctx, rawURL, queryParams)
	if err != nil {
		return nil, err
	}
	if status == http.StatusUnauthorized {
		c.authClient.Invalidate()
		body, status, err = c.doGet(ctx, rawURL, queryParams)
		if err != nil {
			return nil, err
		}
	}
	if status < 200 || status >= 300 {
		return nil, &APIError{URL: rawURL, Status: status, Body: body}
	}
	return body, nil
}

// Post performs a signed, authenticated POST with a form body and returns the response body.
func (c *Client) Post(ctx context.Context, rawURL string, bodyParams url.Values) ([]byte, error) {
	body, status, err := c.doPost(ctx, rawURL, bodyParams)
	if err != nil {
		return nil, err
	}
	if status == http.StatusUnauthorized {
		c.authClient.Invalidate()
		body, status, err = c.doPost(ctx, rawURL, bodyParams)
		if err != nil {
			return nil, err
		}
	}
	if status < 200 || status >= 300 {
		return nil, &APIError{URL: rawURL, Status: status, Body: body}
	}
	return body, nil
}

func (c *Client) doGet(ctx context.Context, rawURL string, queryParams url.Values) ([]byte, int, error) {
	token, err := c.authClient.Token(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get token: %w", err)
	}

	seed, signature, err := c.signer.Sign(http.MethodGet, queryParams, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("sign request: %w", err)
	}

	fullURL := rawURL
	if len(queryParams) > 0 {
		fullURL = rawURL + "?" + queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, 0, err
	}
	c.setCommonHeaders(req, token, seed, signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}

func (c *Client) doPost(ctx context.Context, rawURL string, bodyParams url.Values) ([]byte, int, error) {
	token, err := c.authClient.Token(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get token: %w", err)
	}

	seed, signature, err := c.signer.Sign(http.MethodPost, nil, bodyParams)
	if err != nil {
		return nil, 0, fmt.Errorf("sign request: %w", err)
	}

	encoded := bodyParams.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewBufferString(encoded))
	if err != nil {
		return nil, 0, err
	}
	c.setCommonHeaders(req, token, seed, signature)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}

func (c *Client) setCommonHeaders(req *http.Request, token, seed, signature string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("app_version", config.AppVersion)
	req.Header.Set("device_identifier", c.deviceIdentifier)
	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("seed", seed)
	req.Header.Set("Signature", signature)
}
