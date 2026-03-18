package selfservice

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

var bindFields = []string{"FLDEXTRA1", "FLDEXTRA2", "FLDEXTRA3", "FLDEXTRA4"}

// BindOperator executes the SSOT-required write workflow for operator binding:
// PreState -> Submit -> ReadBack.
//
// It never treats HTTP 200 as success. Success is confirmed only when readback
// values match the intended target field values.
func (c *Client) BindOperator(ctx context.Context, fldextra map[string]string) error {
	if c == nil || c.session == nil {
		return fmt.Errorf("selfservice bind operator: session client is nil: %w", core.ErrAuth)
	}

	if len(fldextra) == 0 {
		return fmt.Errorf("selfservice bind operator: no target fields provided: %w", core.ErrBusinessFailed)
	}

	// 1) PreState: fetch fresh csrftoken and current FLDEXTRA values.
	preToken, preState, err := c.getOperatorState(ctx)
	if err != nil {
		return err
	}

	// 2) Submit: include csrftoken and submit form values.
	form := map[string]string{
		"csrftoken": preToken,
	}

	// Keep non-target fields stable by default and override only requested keys.
	for _, k := range bindFields {
		form[k] = preState[k]
	}
	for k, v := range fldextra {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		form[k] = v
	}

	if _, err := c.session.PostForm(ctx, "/Self/service/bind-operator", core.RequestOptions{Form: form}); err != nil {
		return fmt.Errorf("selfservice bind operator submit failed: %w", err)
	}

	// 3) ReadBack: fetch latest state and verify target-field convergence.
	_, postState, err := c.getOperatorState(ctx)
	if err != nil {
		return err
	}

	for k, expected := range fldextra {
		field := strings.TrimSpace(k)
		if field == "" {
			continue
		}
		got := postState[field]
		if got != expected {
			return fmt.Errorf(
				"selfservice bind operator readback mismatch for %s: expected=%q got=%q: %w",
				field,
				expected,
				got,
				core.ErrBusinessFailed,
			)
		}
	}

	return nil
}

func (c *Client) getOperatorState(ctx context.Context) (string, map[string]string, error) {
	resp, err := c.session.Get(ctx, "/Self/service/operatorId", core.RequestOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("selfservice operatorId request failed: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
	if err != nil {
		return "", nil, fmt.Errorf("selfservice operatorId parse failed: %w", err)
	}

	token, ok := doc.Find("input[name='csrftoken']").First().Attr("value")
	token = strings.TrimSpace(token)
	if !ok || token == "" {
		return "", nil, fmt.Errorf("selfservice operatorId missing csrftoken: %w", core.ErrToken)
	}

	state := map[string]string{}
	for _, field := range bindFields {
		state[field] = extractInputValue(doc, field)
	}

	return token, state, nil
}

func extractInputValue(doc *goquery.Document, name string) string {
	value, ok := doc.Find("input[name='" + name + "']").First().Attr("value")
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}
