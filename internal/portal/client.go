package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

const (
	portal802LoginURL  = "https://10.10.244.11:802/eportal/portal/login"
	portal802LogoutURL = "https://10.10.244.11:802/eportal/portal/logout"
	jsonpCallback      = "dr1003"
)

// PortalResponse captures Portal 802 JSONP payload fields.
// Result and RetCode stay dynamic because upstream can return number or string.
type PortalResponse struct {
	Result  any    `json:"result"`
	Msg     string `json:"msg"`
	RetCode any    `json:"ret_code"`
}

// Client is the Portal 802 protocol client.
type Client struct {
	session core.SessionClient
}

func NewClient(session core.SessionClient) *Client {
	return &Client{session: session}
}

// Logout sends a best-effort portal logout call for gateway-side cleanup.
// It only requires transport success and intentionally ignores the JSONP business body.
func (c *Client) Logout(ctx context.Context, ip string) error {
	if c.session == nil {
		return fmt.Errorf("portal logout: session client is nil: %w", core.ErrPortal)
	}

	_, err := c.session.Get(ctx, portal802LogoutURL, core.RequestOptions{
		Query: map[string]string{
			"callback":     jsonpCallback,
			"login_method": "1",
			"wlan_user_ip": ip,
		},
	})
	if err != nil {
		return fmt.Errorf("portal logout request failed: %w", err)
	}

	return nil
}

// Login performs Portal 802 login with one retry after passive cleanup.
// On first business failure, it calls Logout once to clean possible stale gateway state,
// then retries exactly once to avoid infinite loops.
func (c *Client) Login(ctx context.Context, account, password, ip, isp string) error {
	return c.login(ctx, account, password, ip, isp, false)
}

func (c *Client) login(ctx context.Context, account, password, ip, isp string, isRetry bool) error {
	if c.session == nil {
		return fmt.Errorf("portal login: session client is nil: %w", core.ErrPortal)
	}

	userAccount := ",0," + strings.TrimSpace(account) + ispSuffix(isp)

	resp, err := c.session.Get(ctx, portal802LoginURL, core.RequestOptions{
		Query: map[string]string{
			"callback":      jsonpCallback,
			"login_method":  "1",
			"user_account":  userAccount,
			"user_password": password,
			"wlan_user_ip":  ip,
		},
	})
	if err != nil {
		return fmt.Errorf("portal login request failed: %w", err)
	}

	payload, err := parseJSONPPayload(string(resp.Body), jsonpCallback)
	if err != nil {
		return fmt.Errorf("portal login parse response failed: %w", err)
	}

	if isSuccessResult(payload.Result) {
		return nil
	}

	if !isRetry {
		_ = c.Logout(ctx, ip)
		return c.login(ctx, account, password, ip, isp, true)
	}

	return fmt.Errorf(
		"portal login failed ret_code=%v msg=%q: %w",
		payload.RetCode,
		strings.TrimSpace(payload.Msg),
		core.ErrPortal,
	)
}

func ispSuffix(isp string) string {
	switch strings.ToLower(strings.TrimSpace(isp)) {
	case "telecom":
		return "@dx"
	case "unicom":
		return "@lt"
	case "mobile":
		return "@cmcc"
	default:
		return ""
	}
}

func parseJSONPPayload(raw string, callback string) (*PortalResponse, error) {
	body := strings.TrimSpace(raw)
	prefix := strings.TrimSpace(callback) + "("

	if !strings.HasPrefix(body, prefix) {
		return nil, fmt.Errorf("invalid jsonp prefix")
	}

	body = strings.TrimPrefix(body, prefix)
	body = strings.TrimSpace(body)
	body = strings.TrimSuffix(body, ");")
	body = strings.TrimSuffix(body, ")")
	body = strings.TrimSpace(body)

	if body == "" {
		return nil, fmt.Errorf("empty jsonp payload")
	}

	var out PortalResponse
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		return nil, fmt.Errorf("unmarshal jsonp payload: %w", err)
	}

	return &out, nil
}

func isSuccessResult(v any) bool {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val) == "1"
	case float64:
		return val == 1
	case float32:
		return val == 1
	case int:
		return val == 1
	case int64:
		return val == 1
	case int32:
		return val == 1
	case uint:
		return val == 1
	case uint64:
		return val == 1
	case uint32:
		return val == 1
	default:
		return false
	}
}
