// Package flatten provides utilities for flattening nested JSON log fields
// into a single-level map using dot-notation keys.
package flatten

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

const defaultSeparator = "."

// Flattener collapses nested map fields in log entries into dot-notation keys.
type Flattener struct {
	separator string
}

// New returns a Flattener using the given separator. If sep is empty,
// a dot (".") is used as the default separator.
func New(sep string) *Flattener {
	if sep == "" {
		sep = defaultSeparator
	}
	return &Flattener{separator: sep}
}

// Apply returns a copy of the entry with all nested map fields flattened.
// Existing top-level fields are preserved; nested keys are merged using
// dot-notation (e.g. "http.request.method").
func (f *Flattener) Apply(e parser.Entry) parser.Entry {
	out := parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Raw:       e.Raw,
		Fields:    make(map[string]any),
	}

	for k, v := range e.Fields {
		flattenValue(out.Fields, k, v, f.separator)
	}

	return out
}

// flattenValue recursively walks v. If v is a map[string]any it recurses;
// otherwise it assigns the value at the composed key.
func flattenValue(dst map[string]any, prefix string, v any, sep string) {
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			newKey := strings.Join([]string{prefix, k}, sep)
			flattenValue(dst, newKey, child, sep)
		}
	default:
		// Avoid overwriting an existing key by appending a suffix.
		key := prefix
		if _, exists := dst[key]; exists {
			key = fmt.Sprintf("%s%s_dup", prefix, sep)
		}
		dst[key] = val
	}
}
