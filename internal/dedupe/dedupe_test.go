package dedupe_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/dedupe"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(msg, level string, ts time.Time) parser.Entry {
	return parser.Entry{Message: msg, Level: level, Timestamp: ts}
}

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNew_InvalidWindowSize(t *testing.T) {
	_, err := dedupe.New(-1)
	if err == nil {
		t.Fatal("expected error for negative window size")
	}
}

func TestNew_ValidWindowSize(t *testing.T) {
	_, err := dedupe.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_NoDuplicates(t *testing.T) {
	dd, _ := dedupe.New(10)
	entries := []parser.Entry{
		makeEntry("a", "info", base),
		makeEntry("b", "info", base.Add(time.Second)),
		makeEntry("c", "warn", base.Add(2*time.Second)),
	}
	out := dd.Apply(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestApply_RemovesDuplicates(t *testing.T) {
	dd, _ := dedupe.New(10)
	e := makeEntry("hello", "info", base)
	entries := []parser.Entry{e, e, e}
	out := dd.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_WindowEviction(t *testing.T) {
	// window of 2: after two unique entries the first hash is evicted,
	// so a re-appearance of the first entry should pass through.
	dd, _ := dedupe.New(2)
	a := makeEntry("a", "info", base)
	b := makeEntry("b", "info", base.Add(time.Second))
	c := makeEntry("c", "info", base.Add(2*time.Second))

	entries := []parser.Entry{a, b, c, a}
	out := dd.Apply(entries)
	// a, b, c are unique; after c the window holds b,c so a is evicted -> a passes again
	if len(out) != 4 {
		t.Fatalf("expected 4 entries after eviction, got %d", len(out))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	dd, _ := dedupe.New(5)
	out := dd.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}
