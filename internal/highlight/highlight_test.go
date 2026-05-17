package highlight_test

import (
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/highlight"
)

func TestNew_DisabledReturnsUnchanged(t *testing.T) {
	h := highlight.New([]string{"error"}, false)
	line := "this is an error message"
	got := h.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestNew_NoKeywordsReturnsUnchanged(t *testing.T) {
	h := highlight.New([]string{}, true)
	line := "nothing to highlight here"
	got := h.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_SingleKeyword(t *testing.T) {
	h := highlight.New([]string{"error"}, true)
	line := "an error occurred"
	got := h.Apply(line)
	if !strings.Contains(got, "error") {
		t.Errorf("expected keyword present in output, got %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI codes in output, got %q", got)
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	h := highlight.New([]string{"error"}, true)
	line := "An ERROR occurred"
	got := h.Apply(line)
	// original casing preserved inside the highlight
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected original casing preserved, got %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI codes in output, got %q", got)
	}
}

func TestApply_MultipleKeywords(t *testing.T) {
	h := highlight.New([]string{"warn", "timeout"}, true)
	line := "warn: connection timeout"
	got := h.Apply(line)
	if !strings.Contains(got, "warn") {
		t.Errorf("expected 'warn' in output, got %q", got)
	}
	if !strings.Contains(got, "timeout") {
		t.Errorf("expected 'timeout' in output, got %q", got)
	}
	if strings.Count(got, "\033[") < 2 {
		t.Errorf("expected at least 2 ANSI sequences, got %q", got)
	}
}

func TestApply_KeywordNotPresent(t *testing.T) {
	h := highlight.New([]string{"critical"}, true)
	line := "everything is fine"
	got := h.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line when keyword absent, got %q", got)
	}
}

func TestWithColor_CustomCode(t *testing.T) {
	const cyan = "\033[36m"
	h := highlight.New([]string{"info"}, true).WithColor(cyan + "info" + "\033[0m")
	line := "info level log"
	got := h.Apply(line)
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI codes with custom color, got %q", got)
	}
}
