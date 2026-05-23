// Package ratelimit provides a token-bucket rate limiter for log entry processing.
// It allows callers to cap the number of log entries emitted per second,
// preventing downstream systems from being overwhelmed during high-volume bursts.
package ratelimit

import (
	"errors"
	"sync"
	"time"

	"github.com/humanlogio/logslice/internal/parser"
)

// ErrInvalidRate is returned when a non-positive rate is supplied to New.
var ErrInvalidRate = errors.New("ratelimit: rate must be greater than zero")

// Limiter enforces a maximum number of log entries per second using a
// token-bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens added per nanosecond
	lastTick time.Time
}

// New creates a Limiter that allows at most ratePerSec entries per second.
// It returns ErrInvalidRate if ratePerSec is less than or equal to zero.
func New(ratePerSec float64) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, ErrInvalidRate
	}
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec / float64(time.Second),
		lastTick: time.Now(),
	}, nil
}

// Allow reports whether the next log entry should be allowed through.
// It refills tokens based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := float64(now.Sub(l.lastTick))
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// Apply filters entries, returning only those that pass the rate limit.
func (l *Limiter) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if l.Allow() {
			out = append(out, e)
		}
	}
	return out
}
