package aggregator_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/aggregator"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(level, service string) parser.Entry {
	fields := map[string]string{}
	if service != "" {
		fields["service"] = service
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   "test message",
		Fields:    fields,
	}
}

func TestAggregate_Total(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", ""),
		makeEntry("ERROR", ""),
		makeEntry("INFO", ""),
	}
	summary := aggregator.Aggregate(entries)
	if summary.Total != 3 {
		t.Errorf("expected Total=3, got %d", summary.Total)
	}
}

func TestAggregate_ByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", ""),
		makeEntry("ERROR", ""),
		makeEntry("INFO", ""),
		makeEntry("WARN", ""),
	}
	summary := aggregator.Aggregate(entries)
	if len(summary.ByLevel) != 3 {
		t.Fatalf("expected 3 levels, got %d", len(summary.ByLevel))
	}
	// INFO should be first (highest count)
	if summary.ByLevel[0].Level != "INFO" || summary.ByLevel[0].Count != 2 {
		t.Errorf("expected INFO:2 first, got %s:%d", summary.ByLevel[0].Level, summary.ByLevel[0].Count)
	}
}

func TestAggregate_UnknownLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("", ""),
		makeEntry("", ""),
	}
	summary := aggregator.Aggregate(entries)
	if summary.ByLevel[0].Level != "UNKNOWN" {
		t.Errorf("expected UNKNOWN level, got %s", summary.ByLevel[0].Level)
	}
	if summary.ByLevel[0].Count != 2 {
		t.Errorf("expected count 2, got %d", summary.ByLevel[0].Count)
	}
}

func TestAggregate_ByField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "auth"),
		makeEntry("ERROR", "auth"),
		makeEntry("INFO", "gateway"),
	}
	summary := aggregator.Aggregate(entries)
	serviceCounts, ok := summary.ByField["service"]
	if !ok {
		t.Fatal("expected 'service' field in ByField")
	}
	if serviceCounts["auth"] != 2 {
		t.Errorf("expected auth=2, got %d", serviceCounts["auth"])
	}
	if serviceCounts["gateway"] != 1 {
		t.Errorf("expected gateway=1, got %d", serviceCounts["gateway"])
	}
}

func TestAggregate_Empty(t *testing.T) {
	summary := aggregator.Aggregate([]parser.Entry{})
	if summary.Total != 0 {
		t.Errorf("expected Total=0, got %d", summary.Total)
	}
	if len(summary.ByLevel) != 0 {
		t.Errorf("expected empty ByLevel")
	}
}
