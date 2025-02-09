package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"strings"
	"testing"

	"github.com/calamity-m/reaphur/pkg/bindings"
)

func assertStructuredLog(t testing.TB, buffer *bytes.Buffer, want map[string]interface{}) {
	t.Helper()

	var decoded map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(buffer.String())).Decode(&decoded); err != nil {
		t.Errorf("failed decoding json from log buffer: %v", err)
	}

	// We don't want to compare the time, I'm not insane
	delete(decoded, "time")

	if !maps.Equal(decoded, want) {
		t.Errorf("got %v but want %v", decoded, want)
	}
}

func TestCustomizedHandler(t *testing.T) {

	t.Run("handler supports text and structured logging", func(t *testing.T) {
		// First test text base
		wantText := map[string]string{
			"msg":   "echo",
			"level": "INFO",
		}
		buffer := bytes.Buffer{}

		handler := NewCustomizedHandler(&buffer, nil)
		logger := slog.New(handler)
		logger.Info("echo")

		// We aren't doing enough tests on the text handler right now, so don't
		// bother with an assertion function. Though that may change in the future
		// if more than just the stdlib text handler is used the underlying handler.
		line := buffer.String()
		for key, val := range wantText {
			check := fmt.Sprintf("%s=%s", key, val)
			if !strings.Contains(line, check) {
				t.Errorf("did not find %q in %q", check, line)
			}
		}

		// Now test structured. Remainder of tests will work off of the structured
		// handler.
		wantStructured := map[string]any{
			"msg":   "echo",
			"level": "INFO",
		}
		buffer = bytes.Buffer{}

		handler = NewCustomizedHandler(&buffer, &CustomHandlerCfg{Structed: true})
		logger = slog.New(handler)
		logger.Info("echo")

		assertStructuredLog(t, &buffer, wantStructured)
	})

	t.Run("handler logs static attributes successfully", func(t *testing.T) {
		want := map[string]any{
			"msg":    "echo",
			"level":  "INFO",
			"string": "one",
			"int":    float64(1), // slog casts ints to float64 by default it seems
			"bool":   true,
		}
		buffer := bytes.Buffer{}

		handler := NewCustomizedHandler(&buffer, &CustomHandlerCfg{
			Structed: true,
			StaticAttributes: []slog.Attr{
				slog.String("string", "one"),
				slog.Int("int", 1),
				slog.Bool("bool", true),
			},
		})
		logger := slog.New(handler)
		logger.Info("echo")

		assertStructuredLog(t, &buffer, want)
	})

}

func TestCustomziedHandlerRequestId(t *testing.T) {

	t.Run("handler successfully logs request id", func(t *testing.T) {
		want := map[string]any{
			"msg":        "echo",
			"level":      "INFO",
			"request-id": "req-id",
		}
		buffer := bytes.Buffer{}

		handler := NewCustomizedHandler(&buffer, &CustomHandlerCfg{
			Structed:        true,
			RecordRequestId: true,
		})
		logger := slog.New(handler)

		logger.InfoContext(context.WithValue(context.Background(), bindings.RequestIDKey{}, "req-id"), "echo")

		assertStructuredLog(t, &buffer, want)
	})

	t.Run("handler successfully ignores missing request id", func(t *testing.T) {
		want := map[string]any{
			"msg":        "echo",
			"level":      "INFO",
			"request-id": "unknown",
		}
		buffer := bytes.Buffer{}

		handler := NewCustomizedHandler(&buffer, &CustomHandlerCfg{
			Structed:        true,
			RecordRequestId: true,
		})
		logger := slog.New(handler)

		logger.InfoContext(context.Background(), "echo")

		assertStructuredLog(t, &buffer, want)
	})

	t.Run("handler successfully ignores request id", func(t *testing.T) {
		want := map[string]any{
			"msg":   "echo",
			"level": "INFO",
		}
		buffer := bytes.Buffer{}

		handler := NewCustomizedHandler(&buffer, &CustomHandlerCfg{
			Structed:        true,
			RecordRequestId: false,
		})
		logger := slog.New(handler)

		logger.InfoContext(context.WithValue(context.Background(), bindings.RequestIDKey{}, "req-id"), "echo")

		assertStructuredLog(t, &buffer, want)
	})
}
