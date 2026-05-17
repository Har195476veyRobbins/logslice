// Package sampler provides log entry sampling functionality for logslice.
// It supports rate-based sampling to reduce the volume of log entries
// processed by the pipeline while preserving statistical representation.
package sampler

import (
	"errors"
	"sync/atomic"

	"github.com/logslice/logslice/internal/parser"
)

// ErrInvalidRate is returned when the sample rate is out of the valid range.
var ErrInvalidRate = errors.New("sample rate must be between 1 and 100")

// Sampler holds the configuration for log entry sampling.
type Sampler struct {
	rate uint64 // percentage of entries to keep (1–100)
	counter atomic.Uint64
}

// New creates a new Sampler with the given rate.
// rate must be between 1 and 100 (inclusive), where 100 means keep all entries.
func New(rate uint64) (*Sampler, error) {
	if rate < 1 || rate > 100 {
		return nil, ErrInvalidRate
	}
	return &Sampler{rate: rate}, nil
}

// Sample returns a filtered slice of entries based on the configured rate.
// It uses a deterministic counter so that results are reproducible within
// a single process run.
func (s *Sampler) Sample(entries []parser.Entry) []parser.Entry {
	if s.rate == 100 {
		return entries
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		n := s.counter.Add(1)
		// Keep entry if (n mod 100) falls within the desired rate window.
		if (n%100)+1 <= s.rate {
			result = append(result, e)
		}
	}
	return result
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}
