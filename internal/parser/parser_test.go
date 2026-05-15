package parser_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_JSONEntry(t *testing.T) {
	p := parser.New()
	line := `{"level":"info","msg":"server started","time":"2024-01-15T10:30:00Z"}`

	entry, err := p.Parse(line)
	require.NoError(t, err)

	assert.Equal(t, parser.FormatJSON, entry.Format)
	assert.Equal(t, "info", entry.Level)
	assert.Equal(t, "server started", entry.Message)
	assert.Equal(t, 2024, entry.Timestamp.Year())
	assert.Equal(t, time.January, entry.Timestamp.Month())
}

func TestParse_JSONWithSeverityField(t *testing.T) {
	p := parser.New()
	line := `{"severity":"error","message":"disk full","timestamp":"2024-03-01T08:00:00Z"}`

	entry, err := p.Parse(line)
	require.NoError(t, err)

	assert.Equal(t, "error", entry.Level)
	assert.Equal(t, "disk full", entry.Message)
}

func TestParse_PlainTextWithLevel(t *testing.T) {
	p := parser.New()
	line := "2024-01-15 10:30:00 ERROR something went wrong"

	entry, err := p.Parse(line)
	require.NoError(t, err)

	assert.Equal(t, parser.FormatPlainText, entry.Format)
	assert.Equal(t, "error", entry.Level)
	assert.Equal(t, line, entry.Message)
}

func TestParse_PlainTextNoLevel(t *testing.T) {
	p := parser.New()
	line := "just a plain log line with no level"

	entry, err := p.Parse(line)
	require.NoError(t, err)

	assert.Equal(t, parser.FormatPlainText, entry.Format)
	assert.Empty(t, entry.Level)
}

func TestParse_EmptyLine(t *testing.T) {
	p := parser.New()

	_, err := p.Parse("")
	assert.ErrorIs(t, err, parser.ErrEmptyLine)

	_, err = p.Parse("   ")
	assert.ErrorIs(t, err, parser.ErrEmptyLine)
}

func TestParse_InvalidJSONFallsBackToPlainText(t *testing.T) {
	p := parser.New()
	line := "{not valid json}"

	entry, err := p.Parse(line)
	require.NoError(t, err)
	assert.Equal(t, parser.FormatPlainText, entry.Format)
}

func TestParse_JSONPreservesAllFields(t *testing.T) {
	p := parser.New()
	line := `{"level":"debug","msg":"test","request_id":"abc-123","duration":42}`

	entry, err := p.Parse(line)
	require.NoError(t, err)

	assert.Equal(t, "abc-123", entry.Fields["request_id"])
	assert.Equal(t, float64(42), entry.Fields["duration"])
}
