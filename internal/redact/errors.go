package redact

import (
	"errors"
	"fmt"
)

// ErrNoRules is returned when neither patterns nor sensitive fields are provided.
var ErrNoRules = errors.New("redact: at least one pattern or sensitive field is required")

// ErrInvalidPattern is returned when a regex pattern fails to compile.
type ErrInvalidPattern struct {
	Pattern string
	Err     error
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprintf("redact: invalid pattern %q: %v", e.Pattern, e.Err)
}

func (e *ErrInvalidPattern) Unwrap() error { return e.Err }
