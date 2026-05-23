// Package redact provides log entry field redaction to mask sensitive values
// such as passwords, tokens, and PII before output.
package redact

import (
	"regexp"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

const defaultMask = "[REDACTED]"

// Redactor replaces sensitive patterns in log message and field values.
type Redactor struct {
	patterns []*regexp.Regexp
	fields   map[string]struct{}
	mask     string
}

// New creates a Redactor that masks values matching any of the given regex
// patterns and any fields whose keys appear in sensitiveFields.
// Returns an error if any pattern fails to compile.
func New(patterns []string, sensitiveFields []string) (*Redactor, error) {
	if len(patterns) == 0 && len(sensitiveFields) == 0 {
		return nil, ErrNoRules
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, &ErrInvalidPattern{Pattern: p, Err: err}
		}
		compiled = append(compiled, re)
	}
	fields := make(map[string]struct{}, len(sensitiveFields))
	for _, f := range sensitiveFields {
		fields[strings.ToLower(f)] = struct{}{}
	}
	return &Redactor{patterns: compiled, fields: fields, mask: defaultMask}, nil
}

// Apply returns a copy of the entry with sensitive data masked.
func (r *Redactor) Apply(e parser.Entry) parser.Entry {
	out := parser.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   r.redactString(e.Message),
		Raw:       e.Raw,
		Fields:    make(map[string]string, len(e.Fields)),
	}
	for k, v := range e.Fields {
		if _, sensitive := r.fields[strings.ToLower(k)]; sensitive {
			out.Fields[k] = r.mask
		} else {
			out.Fields[k] = r.redactString(v)
		}
	}
	return out
}

func (r *Redactor) redactString(s string) string {
	for _, re := range r.patterns {
		s = re.ReplaceAllString(s, r.mask)
	}
	return s
}
