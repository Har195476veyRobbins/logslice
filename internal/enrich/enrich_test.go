package enrich_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/enrich"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    fields,
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := enrich.New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty fields map")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := enrich.New(map[string]string{"": "value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_ValidFields(t *testing.T) {
	e, err := enrich.New(map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Enricher")
	}
}

func TestApply_AttachesFields(t *testing.T) {
	e, _ := enrich.New(map[string]string{"env": "prod", "region": "us-east-1"})
	entries := []parser.Entry{
		makeEntry("hello", nil),
		makeEntry("world", map[string]interface{}{"existing": "yes"}),
	}

	out := e.Apply(entries)

	for _, entry := range out {
		if entry.Fields["env"] != "prod" {
			t.Errorf("expected env=prod, got %v", entry.Fields["env"])
		}
		if entry.Fields["region"] != "us-east-1" {
			t.Errorf("expected region=us-east-1, got %v", entry.Fields["region"])
		}
	}
	if out[1].Fields["existing"] != "yes" {
		t.Error("existing field should be preserved")
	}
}

func TestApply_OverwritesExistingKey(t *testing.T) {
	e, _ := enrich.New(map[string]string{"env": "staging"})
	entries := []parser.Entry{
		makeEntry("msg", map[string]interface{}{"env": "prod"}),
	}

	out := e.Apply(entries)

	if out[0].Fields["env"] != "staging" {
		t.Errorf("expected env to be overwritten to staging, got %v", out[0].Fields["env"])
	}
}

func TestApply_EmptyInput(t *testing.T) {
	e, _ := enrich.New(map[string]string{"env": "prod"})
	out := e.Apply([]parser.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
