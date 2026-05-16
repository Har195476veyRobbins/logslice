package filter

import "errors"

// ErrInvalidLevel is returned when an unrecognised level string is provided
// in a context that requires strict validation (e.g. CLI flag parsing).
var ErrInvalidLevel = errors.New("filter: invalid log level")

// ValidateLevel returns ErrInvalidLevel when level is not a known severity name.
func ValidateLevel(level string) error {
	if level == "" {
		return nil
	}
	if _, ok := levelRank[level]; !ok {
		return ErrInvalidLevel
	}
	return nil
}
