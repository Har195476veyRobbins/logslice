package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

var base = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func entry(level, msg string, t time.Time) parser.Entry {
	return parser.Entry{Level: level, Message: msg, Timestamp: t}
}

func TestFilter_ByLevel(t *testing.T) {
	e := entry("debug", "hello", base)
	if filter.Filter(e, filter.Options{Level: "warn"}) {
		t.Error("debug entry should not pass warn filter")
	}
	e2 := entry("error", "boom", base)
	if !filter.Filter(e2, filter.Options{Level: "warn"}) {
		t.Error("error entry should pass warn filter")
	}
}

func TestFilter_Since(t *testing.T) {
	old := entry("info", "old", base.Add(-time.Hour))
	new_ := entry("info", "new", base.Add(time.Hour))
	opts := filter.Options{Since: base}
	if filter.Filter(old, opts) {
		t.Error("old entry should be excluded by Since")
	}
	if !filter.Filter(new_, opts) {
		t.Error("new entry should pass Since filter")
	}
}

func TestFilter_Until(t *testing.T) {
	old := entry("info", "old", base.Add(-time.Hour))
	new_ := entry("info", "new", base.Add(time.Hour))
	opts := filter.Options{Until: base}
	if !filter.Filter(old, opts) {
		t.Error("old entry should pass Until filter")
	}
	if filter.Filter(new_, opts) {
		t.Error("new entry should be excluded by Until")
	}
}

func TestFilter_Keywords(t *testing.T) {
	e := entry("info", "database connection timeout", base)
	if !filter.Filter(e, filter.Options{Keywords: []string{"timeout"}}) {
		t.Error("entry should match keyword 'timeout'")
	}
	if filter.Filter(e, filter.Options{Keywords: []string{"timeout", "auth"}}) {
		t.Error("entry should not match both keywords")
	}
}

func TestApply(t *testing.T) {
	entries := []parser.Entry{
		entry("debug", "verbose", base),
		entry("info", "started", base),
		entry("error", "failed", base),
	}
	result := filter.Apply(entries, filter.Options{Level: "info"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}
