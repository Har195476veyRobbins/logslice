// Package enrich provides log entry enrichment by attaching static or
// dynamic key-value fields to every entry that passes through the pipeline.
package enrich

import (
	"fmt"

	"github.com/yourorg/logslice/internal/parser"
)

// Enricher attaches additional fields to log entries.
type Enricher struct {
	fields map[string]string
}

// New creates an Enricher that will attach the given static fields to every
// log entry. At least one field must be provided.
func New(fields map[string]string) (*Enricher, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	for k := range fields {
		if k == "" {
			return nil, fmt.Errorf("%w: key must not be empty", ErrInvalidField)
		}
	}
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Enricher{fields: copy}, nil
}

// Apply attaches the configured fields to each entry, returning a new slice.
// Existing fields with the same key are overwritten.
func (e *Enricher) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, entry := range entries {
		merged := make(map[string]interface{}, len(entry.Fields)+len(e.fields))
		for k, v := range entry.Fields {
			merged[k] = v
		}
		for k, v := range e.fields {
			merged[k] = v
		}
		entry.Fields = merged
		result[i] = entry
	}
	return result
}
