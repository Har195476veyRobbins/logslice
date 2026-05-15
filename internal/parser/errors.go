package parser

import "errors"

// Sentinel errors returned by the parser.
var (
	// ErrEmptyLine is returned when an empty or whitespace-only line is parsed.
	ErrEmptyLine = errors.New("parser: empty log line")

	// ErrUnparsableTime is returned when no known time format matches.
	ErrUnparsableTime = errors.New("parser: unable to parse timestamp")
)
