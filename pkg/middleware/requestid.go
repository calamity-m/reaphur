package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/google/uuid"
)

// Header that will be queried and written to
var (
	RequestIDHeader = "X-Request-ID"
)

// Middleware that attaches the request's request id to the request context.
// Optionally if generate is true, the middleware will create a new UUID7
// when the request has not provided a request id.
//
// Works off of the X-Request-ID header.
func RequestIDMiddleware(logger *slog.Logger, writeHeader bool) func(http.Handler) http.Handler {

	gen := func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		return id.String(), err
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(RequestIDHeader)

			// Generate a UUID7 if id not populated
			if id == "" {
				newid, err := gen()
				if err != nil {
					logger.Error(fmt.Sprintf("failed to create a V7 request id due to: %v", err))
				}

				id = newid
				r.Header.Set(RequestIDHeader, id)
			}

			// Write the header if id populated
			if writeHeader {
				w.Header().Add(RequestIDHeader, id)
			}

			// Act like a wrapper, providing the request id to the context along the way
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), bindings.RequestIDKey{}, id)))
		})
	}
}
