// Package dedupe provides log-entry deduplication for logslice.
//
// It removes duplicate log entries within a configurable sliding window of
// recently seen messages.  Entries are compared by a hash of their Level,
// Message, and Timestamp fields, so two entries that share all three values
// are considered identical.
//
// Usage:
//
//	dd, err := dedupe.New(100)   // remember last 100 unique hashes
//	if err != nil {
//	    log.Fatal(err)
//	}
//	filtered := dd.Apply(entries)
//
// A windowSize of 0 is treated as 1, which deduplicates only back-to-back
// identical entries.
package dedupe
