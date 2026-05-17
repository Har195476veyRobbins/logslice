// Package pipeline wires together the parser, filter, aggregator,
// and output formatter into a single cohesive processing pipeline.
package pipeline

import (
	"bufio"
	"io"

	"github.com/logslice/logslice/internal/aggregator"
	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
)

// Config holds the configuration for a pipeline run.
type Config struct {
	// Filter options
	Filter filter.Filter

	// Output format: "text" or "json"
	Format string

	// When true the pipeline emits aggregation stats instead of log lines.
	AggregateOnly bool

	// Field to aggregate by (empty = aggregate by level only).
	AggregateField string
}

// Result is returned after the pipeline has finished processing.
type Result struct {
	// Lines written to the output writer.
	LinesWritten int

	// Aggregation stats (always populated).
	Stats aggregator.Stats
}

// Run reads log lines from r, applies the configured filter, formats
// matching entries and writes them to w.  It returns aggregation stats
// for every entry that passed the filter.
func Run(r io.Reader, w io.Writer, cfg Config) (Result, error) {
	p := parser.New()
	fmt := output.New(cfg.Format)

	var entries []parser.Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := p.Parse(line)
		if err != nil {
			// Skip unparseable lines silently.
			continue
		}
		if !filter.Apply(entry, cfg.Filter) {
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return Result{}, err
	}

	stats := aggregator.Aggregate(entries, cfg.AggregateField)

	var written int
	if !cfg.AggregateOnly {
		for _, e := range entries {
			if err := fmt.Write(w, e); err != nil {
				return Result{}, err
			}
			written++
		}
	}

	return Result{LinesWritten: written, Stats: stats}, nil
}
