// Package parser provides log line parsing for both JSON and plain-text formats.
package parser

import (
	"encoding/json"
	"strings"
	"time"
)

// Format represents the detected log format.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatPlainText
)

// LogEntry represents a parsed log line with normalized fields.
type LogEntry struct {
	Raw       string
	Format    Format
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]interface{}
}

// Parser parses individual log lines into LogEntry structs.
type Parser struct {
	timeFormats []string
}

// New returns a new Parser with default time formats.
func New() *Parser {
	return &Parser{
		timeFormats: []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
		},
	}
}

// Parse attempts to parse a raw log line into a LogEntry.
func (p *Parser) Parse(line string) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrEmptyLine
	}

	entry := &LogEntry{Raw: line}

	if strings.HasPrefix(line, "{") {
		if err := p.parseJSON(line, entry); err == nil {
			return entry, nil
		}
	}

	p.parsePlainText(line, entry)
	return entry, nil
}

func (p *Parser) parseJSON(line string, entry *LogEntry) error {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return err
	}

	entry.Format = FormatJSON
	entry.Fields = fields
	entry.Level = extractString(fields, "level", "severity", "lvl")
	entry.Message = extractString(fields, "message", "msg", "text")

	if ts := extractString(fields, "time", "timestamp", "ts", "@timestamp"); ts != "" {
		entry.Timestamp, _ = p.parseTime(ts)
	}
	return nil
}

func (p *Parser) parsePlainText(line string, entry *LogEntry) {
	entry.Format = FormatPlainText
	entry.Fields = map[string]interface{}{}
	entry.Message = line
	entry.Level = detectLevelFromText(line)
}

func (p *Parser) parseTime(s string) (time.Time, error) {
	for _, fmt := range p.timeFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, ErrUnparsableTime
}

func extractString(fields map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := fields[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func detectLevelFromText(line string) string {
	upper := strings.ToUpper(line)
	for _, lvl := range []string{"ERROR", "WARN", "INFO", "DEBUG", "FATAL", "TRACE"} {
		if strings.Contains(upper, lvl) {
			return strings.ToLower(lvl)
		}
	}
	return ""
}
