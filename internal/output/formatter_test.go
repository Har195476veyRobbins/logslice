package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

var testTime = time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

func TestFormatter_TextNoFields(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	e := output.Entry{Timestamp: testTime, Level: "info", Message: "hello world"}
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %q", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected message in output, got: %q", got)
	}
	if !strings.Contains(got, "2024-03-15T12:00:00Z") {
		t.Errorf("expected timestamp in output, got: %q", got)
	}
}

func TestFormatter_TextWithFields(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	e := output.Entry{
		Level:   "warn",
		Message: "disk usage high",
		Fields:  map[string]interface{}{"host": "srv1"},
	}
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "host=srv1") {
		t.Errorf("expected field in output, got: %q", got)
	}
}

func TestFormatter_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	e := output.Entry{
		Timestamp: testTime,
		Level:     "error",
		Message:   "connection refused",
		Fields:    map[string]interface{}{"port": 5432},
	}
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["level"] != "error" {
		t.Errorf("expected level=error, got %v", m["level"])
	}
	if m["message"] != "connection refused" {
		t.Errorf("expected message, got %v", m["message"])
	}
	if m["timestamp"] != "2024-03-15T12:00:00Z" {
		t.Errorf("expected timestamp, got %v", m["timestamp"])
	}
}

func TestFormatter_FallbackToText(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.Format("invalid"))
	e := output.Entry{Message: "fallback test"}
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "fallback test") {
		t.Errorf("expected message in fallback output")
	}
}
