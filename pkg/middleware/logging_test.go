package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {

	// Setup test dummies
	buffer := bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(&buffer, nil))
	rq := httptest.NewRequest(http.MethodGet, "/", nil)
	rs := httptest.NewRecorder()

	// Run the middleware
	LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "failed")
		w.WriteHeader(400)
	})).ServeHTTP(rs, rq)

	// Parse buffer contents from logging operations. Test should only be producing one
	// log line, any more will be considered a parsing failure and thus an error failure.
	var val map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(buffer.String())).Decode(&val); err != nil {
		t.Errorf("failed decoding json from log buffer: %v", err)
	}

	// Verify
	if val == nil {
		t.Errorf("failed to write some log line")
	}

	if val["msg"] == "" {
		t.Errorf("failed to record some message")
	}

	if val["method"] != "GET" {
		t.Errorf("failed to record method, got %v", val["method"])
	}

	if val["status"] != float64(400) {
		t.Errorf("failed to record status code, got %v", val["status"])
	}

	if val["duration"] == nil || val["duration"] == "" {
		t.Errorf("failed to record duration")
	}
}
