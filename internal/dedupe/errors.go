package dedupe

import "errors"

// ErrInvalidWindowSize is returned when a negative window size is provided.
var ErrInvalidWindowSize = errors.New("dedupe: window size must be zero or positive")
