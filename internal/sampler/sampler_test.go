package sampler

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "test message",
			Fields:    map[string]interface{}{},
		}
	}
	return entries
}

func TestNew_InvalidRate(t *testing.T) {
	for _, rate := range []uint64{0, 101, 200} {
		_, err := New(rate)
		if err == nil {
			t.Errorf("expected error for rate %d, got nil", rate)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	for _, rate := range []uint64{1, 50, 100} {
		s, err := New(rate)
		if err != nil {
			t.Fatalf("unexpected error for rate %d: %v", rate, err)
		}
		if s.Rate() != rate {
			t.Errorf("expected rate %d, got %d", rate, s.Rate())
		}
	}
}

func TestSample_FullRate(t *testing.T) {
	s, _ := New(100)
	entries := makeEntries(50)
	got := s.Sample(entries)
	if len(got) != 50 {
		t.Errorf("expected 50 entries at 100%% rate, got %d", len(got))
	}
}

func TestSample_ZeroEntries(t *testing.T) {
	s, _ := New(50)
	got := s.Sample([]parser.Entry{})
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}

func TestSample_ApproximateRate(t *testing.T) {
	s, _ := New(50)
	entries := makeEntries(1000)
	got := s.Sample(entries)
	// Allow ±5% tolerance around the expected 500 entries.
	if len(got) < 450 || len(got) > 550 {
		t.Errorf("expected ~500 entries at 50%% rate, got %d", len(got))
	}
}

func TestSample_LowRate(t *testing.T) {
	s, _ := New(10)
	entries := makeEntries(1000)
	got := s.Sample(entries)
	// Allow ±5% tolerance around the expected 100 entries.
	if len(got) < 50 || len(got) > 150 {
		t.Errorf("expected ~100 entries at 10%% rate, got %d", len(got))
	}
}
