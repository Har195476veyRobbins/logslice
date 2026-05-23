// Package truncate provides field-level truncation for log entries.
//
// Truncation is useful when log messages or structured fields contain
// excessively long values that reduce readability or exceed downstream
// size limits (e.g., line-length caps in log shippers).
//
// Usage:
//
//	tr, err := truncate.New(120)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Truncate a single entry
//	out := tr.Apply(entry)
//
//	// Truncate a batch of entries
//	outs := tr.ApplyAll(entries)
//
// Both the Message field and any string-typed Fields values are subject
// to truncation. Numeric or other non-string field values are passed
// through unchanged. Truncated strings are suffixed with "...".
package truncate
