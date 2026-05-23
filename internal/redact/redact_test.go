package redact

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    fields,
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil, nil)
	if err != ErrNoRules {
		t.Fatalf("expected ErrNoRules, got %v", err)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]string{"[invalid"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
	var pe *ErrInvalidPattern
	if !isErrInvalidPattern(err, &pe) {
		t.Fatalf("expected *ErrInvalidPattern, got %T", err)
	}
}

func TestApply_RedactsMessagePattern(t *testing.T) {
	r, err := New([]string{`password=\S+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	e := makeEntry("login password=secret123 ok", nil)
	out := r.Apply(e)
	if out.Message != "login [REDACTED] ok" {
		t.Errorf("unexpected message: %q", out.Message)
	}
}

func TestApply_RedactsSensitiveField(t *testing.T) {
	r, err := New(nil, []string{"token", "password"})
	if err != nil {
		t.Fatal(err)
	}
	e := makeEntry("auth", map[string]string{"token": "abc123", "user": "alice"})
	out := r.Apply(e)
	if out.Fields["token"] != "[REDACTED]" {
		t.Errorf("token should be redacted, got %q", out.Fields["token"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("user should be unchanged, got %q", out.Fields["user"])
	}
}

func TestApply_PreservesOriginalEntry(t *testing.T) {
	r, err := New([]string{`secret`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	e := makeEntry("secret value", nil)
	_ = r.Apply(e)
	if e.Message != "secret value" {
		t.Error("original entry message should not be mutated")
	}
}

func TestApply_FieldPatternRedaction(t *testing.T) {
	r, err := New([]string{`\d{4}-\d{4}-\d{4}-\d{4}`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	e := makeEntry("payment", map[string]string{"card": "1234-5678-9012-3456"})
	out := r.Apply(e)
	if out.Fields["card"] != "[REDACTED]" {
		t.Errorf("card number should be redacted, got %q", out.Fields["card"])
	}
}

func isErrInvalidPattern(err error, target **ErrInvalidPattern) bool {
	if pe, ok := err.(*ErrInvalidPattern); ok {
		*target = pe
		return true
	}
	return false
}
