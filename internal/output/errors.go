package output

import (
	"errors"
	"fmt"
)

// ErrUnknownFormat is returned when an unrecognised output format is requested.
var ErrUnknownFormat = errors.New("output: unknown format")

// ErrNilEntry is returned when a nil log entry is passed to Format.
var ErrNilEntry = errors.New("output: entry must not be nil")

// ValidateFormat returns an error if f is not a supported Format value.
func ValidateFormat(f Format) error {
	switch f {
	case FormatText, FormatJSON:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrUnknownFormat, f)
	}
}
