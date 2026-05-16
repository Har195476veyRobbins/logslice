// Package filter implements log entry filtering for logslice.
//
// Filtering is performed against [parser.Entry] values produced by the
// internal/parser package.  Three independent criteria are supported and
// combined with logical AND semantics:
//
//   - Level: minimum severity threshold (trace < debug < info < warn < error < fatal).
//   - Time range: Since and Until boundaries (inclusive).
//   - Keywords: all provided strings must appear (case-insensitive) in the
//     entry message.
//
// Example:
//
//	opts := filter.Options{
//		Level:    "warn",
//		Since:    time.Now().Add(-1 * time.Hour),
//		Keywords: []string{"database"},
//	}
//	matched := filter.Apply(entries, opts)
package filter
