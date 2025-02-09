package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/calamity-m/reaphur/pkg/bindings"
)

func TestRequestIDMiddleware(t *testing.T) {

	t.Run("extract and maintain existing request id", func(t *testing.T) {
		want := "faked-request-id"

		// Setup test dummies
		buf := bytes.Buffer{}
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add(RequestIDHeader, want)

		// Create it
		middleware := RequestIDMiddleware(logger, true)

		// Run our middleware and verify
		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Context().Value(bindings.RequestIDKey{})

			if got != want {
				t.Errorf("got %q but want %q for request id ctx value", got, want)
			}
		})).ServeHTTP(response, request)

		if response.Header().Get(RequestIDHeader) != want {
			t.Errorf("got %q but want %q for request id header", response.Header().Get("X-Request-ID"), want)
		}
	})

	t.Run("generate new request id", func(t *testing.T) {
		// Setup test dummies
		buf := bytes.Buffer{}
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// Create it
		middleware := RequestIDMiddleware(logger, true)

		// Run our middleware and verify
		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Context().Value(bindings.RequestIDKey{})

			if got == "" {
				t.Error("got nothing, but want some generated id")
			}
		})).ServeHTTP(response, request)

		if response.Header().Get(RequestIDHeader) == "" {
			t.Error("there should be some request header content")
		}
	})

	t.Run("omit header", func(t *testing.T) {
		want := "faked-request-id"

		// Setup test dummies
		buf := bytes.Buffer{}
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add(RequestIDHeader, want)

		// Create it
		middleware := RequestIDMiddleware(logger, false)

		// Run our middleware and verify
		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Context().Value(bindings.RequestIDKey{})

			if got != want {
				t.Errorf("got %q but want %q for request id ctx value", got, want)
			}
		})).ServeHTTP(response, request)

		if response.Header().Get(RequestIDHeader) == want {
			t.Errorf("got %q but there should be no request header", response.Header().Get("X-Request-ID"))
		}
	})

}
