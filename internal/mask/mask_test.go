package mask_test

import (
	"testing"
	"time"

	"logslice/internal/mask"
	"logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := mask.New(nil, "")
	if err == nil {
		t.Fatal("expected error for empty field list, got nil")
	}
}

func TestNew_ValidFields(t *testing.T) {
	m, err := mask.New([]string{"password"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Masker")
	}
}

func TestApply_NoFields(t *testing.T) {
	m, _ := mask.New([]string{"token"}, "")
	e := makeEntry(nil)
	out := m.Apply(e)
	if out.Message != e.Message {
		t.Errorf("message changed unexpectedly")
	}
	if out.Fields != nil && len(out.Fields) != 0 {
		t.Errorf("expected empty fields, got %v", out.Fields)
	}
}

func TestApply_MasksTargetField(t *testing.T) {
	m, _ := mask.New([]string{"password"}, "***")
	e := makeEntry(map[string]any{"password": "s3cr3t", "user": "alice"})
	out := m.Apply(e)

	if out.Fields["password"] != "***" {
		t.Errorf("expected password to be masked, got %v", out.Fields["password"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("expected user to be unchanged, got %v", out.Fields["user"])
	}
}

func TestApply_CaseInsensitiveFieldMatch(t *testing.T) {
	m, _ := mask.New([]string{"Authorization"}, "[REDACTED]")
	e := makeEntry(map[string]any{"authorization": "Bearer xyz"})
	out := m.Apply(e)

	if out.Fields["authorization"] != "[REDACTED]" {
		t.Errorf("expected field to be masked case-insensitively, got %v", out.Fields["authorization"])
	}
}

func TestApply_DefaultMaskString(t *testing.T) {
	m, _ := mask.New([]string{"token"}, "") // empty mask → default
	e := makeEntry(map[string]any{"token": "abc123"})
	out := m.Apply(e)

	if out.Fields["token"] != "***" {
		t.Errorf("expected default mask '***', got %v", out.Fields["token"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m, _ := mask.New([]string{"secret"}, "MASKED")
	origFields := map[string]any{"secret": "original", "id": 42}
	e := makeEntry(origFields)
	m.Apply(e)

	if origFields["secret"] != "original" {
		t.Error("original entry fields were mutated")
	}
}
