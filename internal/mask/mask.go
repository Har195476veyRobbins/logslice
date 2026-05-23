// Package mask provides field-level value masking for log entries.
// It replaces the full value of specified fields with a fixed mask string,
// useful for hiding sensitive data such as tokens or passwords without
// removing the field key from the entry.
package mask

import (
	"errors"
	"strings"

	"logslice/internal/parser"
)

const defaultMask = "***"

// ErrNoFields is returned when New is called with an empty field list.
var ErrNoFields = errors.New("mask: at least one field name is required")

// Masker replaces the values of configured fields with a mask string.
type Masker struct {
	fields map[string]struct{}
	mask   string
}

// New creates a Masker that will replace values of the given fields.
// mask string may be empty, in which case defaultMask ("***") is used.
// Returns ErrNoFields if fields is empty.
func New(fields []string, maskStr string) (*Masker, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	if maskStr == "" {
		maskStr = defaultMask
	}
	m := &Masker{
		fields: make(map[string]struct{}, len(fields)),
		mask:   maskStr,
	}
	for _, f := range fields {
		m.fields[strings.ToLower(f)] = struct{}{}
	}
	return m, nil
}

// Apply returns a copy of entry with configured field values replaced by the
// mask string. The original entry is never modified.
func (m *Masker) Apply(e parser.Entry) parser.Entry {
	out := parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Raw:       e.Raw,
	}
	if len(e.Fields) == 0 {
		return out
	}
	out.Fields = make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		if _, masked := m.fields[strings.ToLower(k)]; masked {
			out.Fields[k] = m.mask
		} else {
			out.Fields[k] = v
		}
	}
	return out
}
