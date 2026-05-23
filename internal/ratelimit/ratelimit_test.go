package ratelimit_test

import (
	"testing"
	"time"

	"github.com/humanlogio/logslice/internal/parser"
	"github.com/humanlogio/logslice/internal/ratelimit"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Level:   "info",
			Message: "test message",
			Raw:     "test message",
		}
	}
	return entries
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(0)
	if err == nil {
		t.Fatal("expected error for zero rate, got nil")
	}
	_, err = ratelimit.New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := ratelimit.New(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestAllow_BurstUpToRate(t *testing.T) {
	const rate = 10.0
	l, err := ratelimit.New(rate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	allowed := 0
	for i := 0; i < 20; i++ {
		if l.Allow() {
			allowed++
		}
	}
	// Initial bucket is full so first `rate` calls should succeed.
	if allowed != int(rate) {
		t.Errorf("expected %d allowed, got %d", int(rate), allowed)
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l, err := ratelimit.New(1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Drain the bucket.
	for i := 0; i < 1000; i++ {
		l.Allow()
	}
	// Wait for refill.
	time.Sleep(10 * time.Millisecond)
	if !l.Allow() {
		t.Error("expected at least one token after sleep")
	}
}

func TestApply_LimitsOutput(t *testing.T) {
	const rate = 5.0
	l, err := ratelimit.New(rate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := makeEntries(20)
	result := l.Apply(entries)

	if len(result) != int(rate) {
		t.Errorf("expected %d entries, got %d", int(rate), len(result))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	l, _ := ratelimit.New(10)
	result := l.Apply([]parser.Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
