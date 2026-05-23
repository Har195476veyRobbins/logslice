package flatten_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/flatten"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	f := flatten.New("")
	if f == nil {
		t.Fatal("expected non-nil Flattener")
	}
}

func TestApply_NoNestedFields(t *testing.T) {
	f := flatten.New(".")
	e := makeEntry(map[string]any{"host": "localhost", "port": 8080})
	out := f.Apply(e)

	if out.Fields["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", out.Fields["host"])
	}
	if out.Fields["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", out.Fields["port"])
	}
}

func TestApply_SingleLevelNested(t *testing.T) {
	f := flatten.New(".")
	e := makeEntry(map[string]any{
		"http": map[string]any{
			"method": "GET",
			"status": 200,
		},
	})
	out := f.Apply(e)

	if out.Fields["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", out.Fields["http.method"])
	}
	if out.Fields["http.status"] != 200 {
		t.Errorf("expected http.status=200, got %v", out.Fields["http.status"])
	}
	if _, ok := out.Fields["http"]; ok {
		t.Error("expected nested 'http' key to be removed")
	}
}

func TestApply_DeeplyNested(t *testing.T) {
	f := flatten.New(".")
	e := makeEntry(map[string]any{
		"request": map[string]any{
			"headers": map[string]any{
				"content-type": "application/json",
			},
		},
	})
	out := f.Apply(e)

	key := "request.headers.content-type"
	if out.Fields[key] != "application/json" {
		t.Errorf("expected %s=application/json, got %v", key, out.Fields[key])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := flatten.New("_")
	e := makeEntry(map[string]any{
		"db": map[string]any{"name": "postgres"},
	})
	out := f.Apply(e)

	if out.Fields["db_name"] != "postgres" {
		t.Errorf("expected db_name=postgres, got %v", out.Fields["db_name"])
	}
}

func TestApply_PreservesEntryMetadata(t *testing.T) {
	f := flatten.New(".")
	e := makeEntry(map[string]any{"k": "v"})
	out := f.Apply(e)

	if out.Level != e.Level {
		t.Errorf("level mismatch: want %s got %s", e.Level, out.Level)
	}
	if out.Message != e.Message {
		t.Errorf("message mismatch: want %s got %s", e.Message, out.Message)
	}
}
