package selfservice

import (
	"context"
	"fmt"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

// Client is the Self protocol facade placeholder.
type Client struct {
	session core.SessionClient
}

func NewClient(session core.SessionClient) *Client {
	return &Client{session: session}
}

func (c *Client) Health(ctx context.Context) (*core.OperationResult[map[string]string], error) {
	_ = ctx
	if c.session == nil {
		return nil, &core.AuthError{Op: "selfservice.health", Msg: "session client is nil", Err: core.ErrAuth}
	}
	return nil, fmt.Errorf("selfservice health: %w", core.ErrNotImplemented)
}
