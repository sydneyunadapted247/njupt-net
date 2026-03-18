package selfservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

var (
	// ErrCheckcodeTokenMissing marks the preflight token extraction failure.
	ErrCheckcodeTokenMissing = errors.New("checkcode token missing")
)

// Login executes the strict Self login chain defined by FINAL-SSOT.
//
// Sequence:
// 1) Preflight: GET login page and extract checkcode token from HTML.
// 2) Bypass: GET randomCode endpoint to trigger server-side captcha state machine.
// 3) Verify: POST login form with honeypot fields and judge success by redirect target.
func (c *Client) Login(ctx context.Context, account, password string) error {
	if c == nil || c.session == nil {
		return fmt.Errorf("selfservice login: session client is nil: %w", core.ErrAuth)
	}

	checkcode, err := c.preflightCheckcode(ctx)
	if err != nil {
		return err
	}

	if err := c.bypassRandomCode(ctx); err != nil {
		// Network/transport errors are returned as-is with context.
		return err
	}

	return c.verifyLogin(ctx, account, password, checkcode)
}

func (c *Client) preflightCheckcode(ctx context.Context) (string, error) {
	resp, err := c.session.Get(ctx, "/Self/login/?302=LI", core.RequestOptions{})
	if err != nil {
		return "", fmt.Errorf("selfservice login preflight request failed: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
	if err != nil {
		return "", fmt.Errorf("selfservice login preflight parse failed: %w", err)
	}

	checkcode, ok := doc.Find("input[name='checkcode']").First().Attr("value")
	checkcode = strings.TrimSpace(checkcode)
	if !ok || checkcode == "" {
		return "", &core.AuthError{
			Op:  "selfservice.login.preflight",
			Msg: "failed to extract checkcode token from login page",
			Err: fmt.Errorf("%w: %w", core.ErrAuth, ErrCheckcodeTokenMissing),
		}
	}

	return checkcode, nil
}

func (c *Client) bypassRandomCode(ctx context.Context) error {
	_, err := c.session.Get(ctx, "/Self/login/randomCode", core.RequestOptions{
		Query: map[string]string{
			"t": strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	})
	if err != nil {
		return fmt.Errorf("selfservice login randomCode request failed: %w", err)
	}
	return nil
}

func (c *Client) verifyLogin(ctx context.Context, account, password, checkcode string) error {
	resp, err := c.session.PostForm(ctx, "/Self/login/verify", core.RequestOptions{
		Form: map[string]string{
			"account":   account,
			"password":  password,
			"checkcode": checkcode,
			"code":      "",
			"foo":       account,
			"bar":       password,
		},
	})
	if err != nil {
		return fmt.Errorf("selfservice login verify request failed: %w", err)
	}

	if strings.Contains(resp.FinalURL, "/Self/dashboard") {
		return nil
	}

	if strings.Contains(resp.FinalURL, "/Self/login/") {
		msg := extractLoginErrorMessage(resp.Body)
		return &core.AuthError{
			Op:  "selfservice.login.verify",
			Msg: msg,
			Err: core.ErrAuth,
		}
	}

	return &core.AuthError{
		Op:  "selfservice.login.verify",
		Msg: fmt.Sprintf("unexpected login redirect target: %s", strings.TrimSpace(resp.FinalURL)),
		Err: core.ErrAuth,
	}
}

func extractLoginErrorMessage(body []byte) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "login failed"
	}

	candidates := []string{
		".alert-danger",
		".alert",
		"div.error",
		"span.error",
		"#error",
		".swal2-content",
		".swal-content",
	}

	for _, selector := range candidates {
		text := normalizeText(doc.Find(selector).First().Text())
		if text != "" {
			return text
		}
	}

	// Fallback: scan common message-bearing nodes and return the first meaningful line.
	fallback := ""
	doc.Find("div,span,p,li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		_ = i
		text := normalizeText(s.Text())
		if text == "" {
			return true
		}
		if looksLikeErrorMessage(text) {
			fallback = text
			return false
		}
		return true
	})

	if fallback != "" {
		return fallback
	}

	return "login failed"
}

func normalizeText(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func looksLikeErrorMessage(s string) bool {
	t := strings.ToLower(strings.TrimSpace(s))
	if t == "" {
		return false
	}
	keywords := []string{"失败", "错误", "invalid", "error", "fail", "not valid", "登录"}
	for _, kw := range keywords {
		if strings.Contains(t, kw) {
			return true
		}
	}
	return false
}
