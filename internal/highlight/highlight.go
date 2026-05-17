// Package highlight provides keyword highlighting for log output.
// It wraps matched keywords in ANSI escape sequences for terminal display.
package highlight

import (
	"strings"
)

const (
	// ANSI color codes
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Highlighter applies ANSI color codes to keywords within log lines.
type Highlighter struct {
	keywords []string
	color    string
	enabled  bool
}

// New creates a Highlighter for the given keywords.
// If keywords is empty or enabled is false, Apply returns input unchanged.
func New(keywords []string, enabled bool) *Highlighter {
	return &Highlighter{
		keywords: keywords,
		color:    colorYellow + colorBold,
		enabled:  enabled,
	}
}

// WithColor sets a custom ANSI color code for highlights.
func (h *Highlighter) WithColor(ansiCode string) *Highlighter {
	h.color = ansiCode
	return h
}

// Apply wraps each occurrence of any keyword in the line with ANSI color codes.
// Matching is case-insensitive. Returns the original line if highlighting is
// disabled or no keywords are configured.
func (h *Highlighter) Apply(line string) string {
	if !h.enabled || len(h.keywords) == 0 {
		return line
	}

	result := line
	for _, kw := range h.keywords {
		if kw == "" {
			continue
		}
		result = replaceInsensitive(result, kw, h.color+kw+colorReset)
	}
	return result
}

// replaceInsensitive replaces all case-insensitive occurrences of old in s
// with the literal string new, preserving the original casing in the output
// by wrapping the matched segment.
func replaceInsensitive(s, old, newVal string) string {
	if old == "" {
		return s
	}
	lower := strings.ToLower(s)
	target := strings.ToLower(old)

	var b strings.Builder
	offset := 0
	for {
		idx := strings.Index(lower[offset:], target)
		if idx < 0 {
			b.WriteString(s[offset:])
			break
		}
		abs := offset + idx
		b.WriteString(s[offset:abs])
		// wrap original casing with color
		match := s[abs : abs+len(old)]
		b.WriteString(strings.Replace(newVal, old, match, 1))
		offset = abs + len(old)
	}
	return b.String()
}
