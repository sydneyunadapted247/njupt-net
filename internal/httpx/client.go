package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

// DefaultSessionClient is the baseline implementation for core.SessionClient.
//
// Design notes aligned with FINAL-SSOT:
// 1) Cookie jar is always enabled to keep session cookies like JSESSIONID automatically.
// 2) TLS verification is disabled by default to tolerate Portal 802 certificate issues.
// 3) Redirects are blocked by default so callers can inspect original 302 status and Location.
type DefaultSessionClient struct {
	baseURL string
	http    *http.Client
}

// NewDefaultSessionClient creates a client with SSOT-safe defaults.
func NewDefaultSessionClient(baseURL string) (*DefaultSessionClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			_ = req
			_ = via
			// Keep the original redirect response for protocol-level business judgment.
			return http.ErrUseLastResponse
		},
	}

	return &DefaultSessionClient{
		baseURL: strings.TrimSpace(baseURL),
		http:    httpClient,
	}, nil
}

// Get sends a GET request and returns a normalized core.SessionResponse.
func (c *DefaultSessionClient) Get(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
	reqURL, err := c.buildURL(path, opts.Query)
	if err != nil {
		return nil, fmt.Errorf("build get url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create get request: %w", err)
	}
	applyHeaders(req, opts.Headers)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do get request: %w", err)
	}
	defer resp.Body.Close()

	return adaptResponse(resp)
}

// PostForm sends an x-www-form-urlencoded POST request and returns a normalized response.
func (c *DefaultSessionClient) PostForm(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
	reqURL, err := c.buildURL(path, opts.Query)
	if err != nil {
		return nil, fmt.Errorf("build post url: %w", err)
	}

	form := url.Values{}
	for k, v := range opts.Form {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create post request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	applyHeaders(req, opts.Headers)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do post request: %w", err)
	}
	defer resp.Body.Close()

	return adaptResponse(resp)
}

func (c *DefaultSessionClient) buildURL(path string, query map[string]string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	var u *url.URL
	var err error

	if parsed, parseErr := url.Parse(path); parseErr == nil && parsed.Scheme != "" && parsed.Host != "" {
		u = parsed
	} else {
		if strings.TrimSpace(c.baseURL) == "" {
			return "", fmt.Errorf("relative path requires non-empty baseURL")
		}
		base, baseErr := url.Parse(c.baseURL)
		if baseErr != nil {
			return "", fmt.Errorf("parse baseURL: %w", baseErr)
		}
		rel, relErr := url.Parse(path)
		if relErr != nil {
			return "", fmt.Errorf("parse path: %w", relErr)
		}
		u = base.ResolveReference(rel)
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), err
}

func applyHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		if strings.TrimSpace(k) == "" {
			continue
		}
		req.Header.Set(k, v)
	}
}

// adaptResponse converts stdlib response objects into core.SessionResponse.
// FinalURL selection rule:
// - Prefer resolved Location header when present (critical for intercepted redirects).
// - Fallback to response request URL when Location does not exist.
func adaptResponse(resp *http.Response) (*core.SessionResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	finalURL := ""
	if loc := resp.Header.Get("Location"); strings.TrimSpace(loc) != "" {
		if parsed, parseErr := url.Parse(loc); parseErr == nil {
			if resp.Request != nil && resp.Request.URL != nil {
				finalURL = resp.Request.URL.ResolveReference(parsed).String()
			} else {
				finalURL = parsed.String()
			}
		} else {
			finalURL = loc
		}
	} else if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}

	headers := make(map[string][]string, len(resp.Header))
	for k, v := range resp.Header {
		copied := make([]string, len(v))
		copy(copied, v)
		headers[k] = copied
	}

	return &core.SessionResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
		FinalURL:   finalURL,
	}, nil
}
