package dedupe

import (
	"hash/fnv"
	"strconv"

	"github.com/logslice/logslice/internal/parser"
)

// Deduplicator removes consecutive or windowed duplicate log entries.
type Deduplicator struct {
	windowSize int
	seen       map[uint64]int
	order      []uint64
}

// New creates a new Deduplicator. windowSize controls how many unique
// message hashes are remembered; 0 means deduplicate only consecutive
// identical entries.
func New(windowSize int) (*Deduplicator, error) {
	if windowSize < 0 {
		return nil, ErrInvalidWindowSize
	}
	if windowSize == 0 {
		windowSize = 1
	}
	return &Deduplicator{
		windowSize: windowSize,
		seen:       make(map[uint64]int, windowSize),
	}, nil
}

// Apply filters entries, dropping duplicates within the configured window.
func (d *Deduplicator) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		h := hashEntry(e)
		if _, exists := d.seen[h]; exists {
			continue
		}
		d.record(h)
		out = append(out, e)
	}
	return out
}

func (d *Deduplicator) record(h uint64) {
	if len(d.order) >= d.windowSize {
		evict := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, evict)
	}
	d.seen[h] = 1
	d.order = append(d.order, h)
}

func hashEntry(e parser.Entry) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(e.Message))
	_, _ = h.Write([]byte(e.Level))
	_, _ = h.Write([]byte(strconv.FormatInt(e.Timestamp.UnixNano(), 10)))
	return h.Sum64()
}
