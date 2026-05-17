package pipeline_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/pipeline"
)

const sampleLogs = `{"level":"info","msg":"server started","ts":"2024-01-10T10:00:00Z"}
{"level":"error","msg":"connection refused","ts":"2024-01-10T10:01:00Z"}
{"level":"debug","msg":"dialing host","ts":"2024-01-10T10:02:00Z"}
{"level":"info","msg":"request handled","ts":"2024-01-10T10:03:00Z"}
`

func TestRun_NoFilter(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Config{Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 4 {
		t.Errorf("expected 4 lines written, got %d", res.LinesWritten)
	}
	if res.Stats.Total != 4 {
		t.Errorf("expected stats.Total=4, got %d", res.Stats.Total)
	}
}

func TestRun_FilterByLevel(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Config{
		Format: "text",
		Filter: filter.Filter{Level: "error"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 1 {
		t.Errorf("expected 1 line written, got %d", res.LinesWritten)
	}
}

func TestRun_AggregateOnly(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Config{
		Format:        "text",
		AggregateOnly: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 0 {
		t.Errorf("expected 0 lines written in aggregate-only mode, got %d", res.LinesWritten)
	}
	if res.Stats.Total != 4 {
		t.Errorf("expected stats.Total=4, got %d", res.Stats.Total)
	}
}

func TestRun_FilterSince(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer

	since, _ := time.Parse(time.RFC3339, "2024-01-10T10:02:00Z")
	res, err := pipeline.Run(r, &w, pipeline.Config{
		Format: "text",
		Filter: filter.Filter{Since: since},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 2 {
		t.Errorf("expected 2 lines written, got %d", res.LinesWritten)
	}
}
