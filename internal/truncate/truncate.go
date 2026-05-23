// Package truncate provides utilities for truncating log entry fields
// to a maximum byte length, preserving structure and readability.
package truncate

import (
	"fmt"

	"github.com/your-org/logslice/internal/parser"
)

const ellipsis = "..."

// Truncator trims log entry message and field values that exceed MaxLen bytes.
type Truncator struct {
	MaxLen int
}

// New returns a Truncator with the given maximum byte length.
// Returns an error if maxLen is less than len(ellipsis)+1.
func New(maxLen int) (*Truncator, error) {
	if maxLen < len(ellipsis)+1 {
		return nil, fmt.Errorf("truncate: maxLen must be at least %d, got %d", len(ellipsis)+1, maxLen)
	}
	return &Truncator{MaxLen: maxLen}, nil
}

// Apply returns a copy of the entry with message and string field values
// truncated to MaxLen bytes. Non-string field values are left unchanged.
func (t *Truncator) Apply(e parser.Entry) parser.Entry {
	out := parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   t.trim(e.Message),
		Raw:       e.Raw,
		Fields:    make(map[string]any, len(e.Fields)),
	}
	for k, v := range e.Fields {
		if s, ok := v.(string); ok {
			out.Fields[k] = t.trim(s)
		} else {
			out.Fields[k] = v
		}
	}
	return out
}

// ApplyAll applies Apply to each entry in the slice and returns the results.
func (t *Truncator) ApplyAll(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = t.Apply(e)
	}
	return out
}

func (t *Truncator) trim(s string) string {
	if len(s) <= t.MaxLen {
		return s
	}
	return s[:t.MaxLen-len(ellipsis)] + ellipsis
}
