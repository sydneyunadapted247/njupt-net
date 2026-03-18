package workflow

import (
	"context"
	"fmt"

	"github.com/hicancan/njupt-net-cli/internal/core"
	"github.com/hicancan/njupt-net-cli/internal/httpx"
	"github.com/hicancan/njupt-net-cli/internal/selfservice"
)

const defaultSelfBaseURL = "http://10.10.244.240:8080"

// MigrateBroadband performs a guarded multi-step migration workflow:
// 1) login source account
// 2) clear source binding with empty FLDEXTRA values
// 3) create a fresh isolated session client
// 4) login target account
// 5) bind target account with desired fields
//
// The function fails fast and returns immediately when any step errors.
func MigrateBroadband(
	ctx context.Context,
	session core.SessionClient,
	fromUser,
	fromPwd,
	toUser,
	toPwd string,
	targetFields map[string]string,
) error {
	if session == nil {
		return fmt.Errorf("workflow migrate: source session is nil: %w", core.ErrAuth)
	}
	if fromUser == "" || fromPwd == "" || toUser == "" || toPwd == "" {
		return fmt.Errorf("workflow migrate: from/to credentials are required")
	}
	if len(targetFields) == 0 {
		return fmt.Errorf("workflow migrate: target fields are required: %w", core.ErrBusinessFailed)
	}

	// Step 1: login source account in the provided session.
	srcClient := selfservice.NewClient(session)
	if err := srcClient.Login(ctx, fromUser, fromPwd); err != nil {
		return fmt.Errorf("workflow migrate: source login failed: %w", err)
	}

	// Step 2: clear source-side bindings with explicit empty values.
	clearFields := map[string]string{
		"FLDEXTRA1": "",
		"FLDEXTRA2": "",
		"FLDEXTRA3": "",
		"FLDEXTRA4": "",
	}
	if err := srcClient.BindOperator(ctx, clearFields); err != nil {
		return fmt.Errorf("workflow migrate: source unbind failed: %w", err)
	}

	// Step 3: isolate JSESSIONID by creating a fresh transport session.
	freshSession, err := httpx.NewDefaultSessionClient(defaultSelfBaseURL)
	if err != nil {
		return fmt.Errorf("workflow migrate: create fresh session failed: %w", err)
	}

	// Step 4: login target account in the fresh session.
	dstClient := selfservice.NewClient(freshSession)
	if err := dstClient.Login(ctx, toUser, toPwd); err != nil {
		return fmt.Errorf("workflow migrate: target login failed: %w", err)
	}

	// Step 5: bind target fields to destination account.
	if err := dstClient.BindOperator(ctx, targetFields); err != nil {
		return fmt.Errorf("workflow migrate: target bind failed: %w", err)
	}

	return nil
}
