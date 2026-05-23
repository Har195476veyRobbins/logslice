package truncate_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/logslice/internal/parser"
	"github.com/your-org/logslice/internal/truncate"
)

func makeEntry(msg string, fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Raw:       msg,
		Fields:    fields,
	}
}

func TestNew_TooSmallMaxLen(t *testing.T) {
	_, err := truncate.New(2)
	if err == nil {
		t.Fatal("expected error for maxLen < 4")
	}
}

func TestNew_ValidMaxLen(t *testing.T) {
	tr, err := truncate.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.MaxLen != 10 {
		t.Fatalf("expected MaxLen 10, got %d", tr.MaxLen)
	}
}

func TestApply_ShortMessageUnchanged(t *testing.T) {
	tr, _ := truncate.New(20)
	e := makeEntry("hello", nil)
	out := tr.Apply(e)
	if out.Message != "hello" {
		t.Errorf("expected 'hello', got %q", out.Message)
	}
}

func TestApply_LongMessageTruncated(t *testing.T) {
	tr, _ := truncate.New(10)
	long := strings.Repeat("a", 20)
	e := makeEntry(long, nil)
	out := tr.Apply(e)
	if len(out.Message) != 10 {
		t.Errorf("expected length 10, got %d", len(out.Message))
	}
	if !strings.HasSuffix(out.Message, "...") {
		t.Errorf("expected ellipsis suffix, got %q", out.Message)
	}
}

func TestApply_StringFieldTruncated(t *testing.T) {
	tr, _ := truncate.New(8)
	fields := map[string]any{"key": strings.Repeat("x", 20)}
	e := makeEntry("msg", fields)
	out := tr.Apply(e)
	v, ok := out.Fields["key"].(string)
	if !ok {
		t.Fatal("expected string field")
	}
	if len(v) != 8 {
		t.Errorf("expected length 8, got %d", len(v))
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	tr, _ := truncate.New(8)
	fields := map[string]any{"count": 42}
	e := makeEntry("msg", fields)
	out := tr.Apply(e)
	if out.Fields["count"] != 42 {
		t.Errorf("expected 42, got %v", out.Fields["count"])
	}
}

func TestApplyAll_TruncatesAll(t *testing.T) {
	tr, _ := truncate.New(6)
	entries := []parser.Entry{
		makeEntry(strings.Repeat("b", 10), nil),
		makeEntry("short", nil),
	}
	out := tr.ApplyAll(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if len(out[0].Message) != 6 {
		t.Errorf("expected length 6, got %d", len(out[0].Message))
	}
	if out[1].Message != "short" {
		t.Errorf("expected 'short', got %q", out[1].Message)
	}
}
