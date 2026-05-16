// Package filter provides log entry filtering capabilities based on
// level, time range, and keyword matching.
package filter

import (
	"strings"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options holds the filtering criteria applied to log entries.
type Options struct {
	// Level filters entries at or above the given severity (e.g. "warn").
	Level string
	// Since discards entries older than this time.
	Since time.Time
	// Until discards entries newer than this time.
	Until time.Time
	// Keywords requires all listed strings to appear in the message.
	Keywords []string
}

// levelRank maps canonical level names to a numeric rank for comparison.
var levelRank = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}

// Filter returns true when the entry satisfies all criteria in opts.
func Filter(entry parser.Entry, opts Options) bool {
	if opts.Level != "" {
		minRank, ok := levelRank[strings.ToLower(opts.Level)]
		if ok {
			entryRank, known := levelRank[strings.ToLower(entry.Level)]
			if !known || entryRank < minRank {
				return false
			}
		}
	}

	if !opts.Since.IsZero() && entry.Timestamp.Before(opts.Since) {
		return false
	}
	if !opts.Until.IsZero() && entry.Timestamp.After(opts.Until) {
		return false
	}

	msg := strings.ToLower(entry.Message)
	for _, kw := range opts.Keywords {
		if !strings.Contains(msg, strings.ToLower(kw)) {
			return false
		}
	}

	return true
}

// Apply filters a slice of entries and returns only those that match opts.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if Filter(e, opts) {
			out = append(out, e)
		}
	}
	return out
}
