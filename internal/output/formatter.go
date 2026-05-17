// Package output provides formatting utilities for rendering log entries
// to various output formats such as plain text and JSON.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format represents an output format type.
type Format string

const (
	// FormatText renders log entries as human-readable plain text.
	FormatText Format = "text"
	// FormatJSON renders log entries as JSON objects.
	FormatJSON Format = "json"
)

// Entry represents a parsed log entry to be formatted.
type Entry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]interface{}
}

// Formatter writes formatted log entries to an io.Writer.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter that writes to w using the given format.
// If format is unrecognized, FormatText is used as a fallback.
func New(w io.Writer, format Format) *Formatter {
	if format != FormatText && format != FormatJSON {
		format = FormatText
	}
	return &Formatter{format: format, w: w}
}

// Write formats and writes a single Entry to the underlying writer.
func (f *Formatter) Write(e Entry) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(e)
	default:
		return f.writeText(e)
	}
}

func (f *Formatter) writeText(e Entry) error {
	var sb strings.Builder
	if !e.Timestamp.IsZero() {
		sb.WriteString(e.Timestamp.Format(time.RFC3339))
		sb.WriteByte(' ')
	}
	if e.Level != "" {
		sb.WriteString("[")
		sb.WriteString(strings.ToUpper(e.Level))
		sb.WriteString("] ")
	}
	sb.WriteString(e.Message)
	for k, v := range e.Fields {
		sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}
	sb.WriteByte('\n')
	_, err := io.WriteString(f.w, sb.String())
	return err
}

func (f *Formatter) writeJSON(e Entry) error {
	m := make(map[string]interface{}, len(e.Fields)+3)
	for k, v := range e.Fields {
		m[k] = v
	}
	if !e.Timestamp.IsZero() {
		m["timestamp"] = e.Timestamp.Format(time.RFC3339)
	}
	if e.Level != "" {
		m["level"] = e.Level
	}
	m["message"] = e.Message
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = f.w.Write(b)
	return err
}
