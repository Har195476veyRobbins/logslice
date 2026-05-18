package tail_test

import (
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/tail"
)

func makeReader(lines ...string) *strings.Reader {
	return strings.NewReader(strings.Join(lines, "\n"))
}

func TestLines_InvalidCount(t *testing.T) {
	_, err := tail.Lines(makeReader("a", "b"), 0)
	if err == nil {
		t.Fatal("expected error for count=0, got nil")
	}
}

func TestLines_FewerLinesThanN(t *testing.T) {
	got, err := tail.Lines(makeReader("a", "b"), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestLines_ExactN(t *testing.T) {
	got, err := tail.Lines(makeReader("a", "b", "c"), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestLines_LastN(t *testing.T) {
	got, err := tail.Lines(makeReader("a", "b", "c", "d", "e"), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "c" || got[1] != "d" || got[2] != "e" {
		t.Fatalf("expected last 3 lines [c d e], got %v", got)
	}
}

func TestHead_InvalidCount(t *testing.T) {
	_, err := tail.Head(makeReader("a", "b"), 0)
	if err == nil {
		t.Fatal("expected error for count=0, got nil")
	}
}

func TestHead_FirstN(t *testing.T) {
	got, err := tail.Head(makeReader("a", "b", "c", "d"), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected [a b], got %v", got)
	}
}

func TestHead_FewerLinesThanN(t *testing.T) {
	got, err := tail.Head(makeReader("x"), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("expected [x], got %v", got)
	}
}
