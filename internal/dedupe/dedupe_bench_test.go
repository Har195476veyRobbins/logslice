package dedupe_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/dedupe"
	"github.com/logslice/logslice/internal/parser"
)

func BenchmarkApply_UniqueEntries(b *testing.B) {
	const n = 1000
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Message:   fmt.Sprintf("log line %d", i),
			Level:     "info",
			Timestamp: base.Add(time.Duration(i) * time.Millisecond),
		}
	}
	dd, _ := dedupe.New(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dd.Apply(entries)
	}
}

func BenchmarkApply_AllDuplicates(b *testing.B) {
	const n = 1000
	e := parser.Entry{Message: "repeated", Level: "error", Timestamp: base}
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = e
	}
	dd, _ := dedupe.New(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dd.Apply(entries)
	}
}
