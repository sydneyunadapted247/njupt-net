package core

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrAuth           = errors.New("auth failed")
	ErrToken          = errors.New("token missing or invalid")
	ErrBusinessFailed = errors.New("business verification failed")
	ErrPortal         = errors.New("portal request failed")
)

// AuthError represents an authentication stage failure with context.
type AuthError struct {
	Op  string
	Msg string
	Err error
}

func (e *AuthError) Error() string {
	if e == nil {
		return "auth error"
	}
	if e.Msg != "" {
		return e.Op + ": " + e.Msg
	}
	if e.Err != nil {
		return e.Op + ": " + e.Err.Error()
	}
	return e.Op + ": auth error"
}

func (e *AuthError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
