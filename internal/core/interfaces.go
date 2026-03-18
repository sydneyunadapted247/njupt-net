package core

import "context"

// RequestOptions carries basic request inputs for session calls.
type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
	Form    map[string]string
}

// SessionResponse is a normalized transport response container.
type SessionResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	FinalURL   string
}

// SessionClient defines the transport contract consumed by protocol layers.
type SessionClient interface {
	Get(ctx context.Context, path string, opts RequestOptions) (*SessionResponse, error)
	PostForm(ctx context.Context, path string, opts RequestOptions) (*SessionResponse, error)
}
