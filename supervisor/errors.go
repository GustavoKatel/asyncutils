package supervisor

import "errors"

var (
	// ErrHandlerNil user passed a nil handler
	ErrHandlerNil = errors.New("Handler must not be nil")
)
