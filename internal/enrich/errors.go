package enrich

import "errors"

// ErrNoFields is returned when New is called with an empty fields map.
var ErrNoFields = errors.New("enrich: at least one field must be provided")

// ErrInvalidField is returned when a field key is invalid (e.g. empty string).
var ErrInvalidField = errors.New("enrich: invalid field")
