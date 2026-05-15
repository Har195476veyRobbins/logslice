// Package parser implements log line parsing for logslice.
//
// It supports two log formats:
//
//   - JSON: structured log lines beginning with '{'. Common field names such as
//     "level", "severity", "msg", "message", "time", "timestamp" and "ts" are
//     automatically promoted to the top-level LogEntry fields. All other fields
//     are available via LogEntry.Fields.
//
//   - Plain-text: unstructured log lines. The parser attempts to detect a log
//     level by scanning the line for well-known level keywords (ERROR, WARN,
//     INFO, DEBUG, FATAL, TRACE) and stores the full line as the message.
//
// Usage:
//
//	p := parser.New()
//	entry, err := p.Parse(line)
//	if err != nil {
//		// handle ErrEmptyLine or other errors
//	}
//	fmt.Println(entry.Level, entry.Message)
package parser
