package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

// Supplant the inner http.ResponseWriter so we can spy on
// the status code being sent out with the response.
func (l *loggingWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.statusCode = statusCode
}

// Logs the end of request action with tracing information, including
// the duration the request took.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := loggingWriter{ResponseWriter: w}

			h.ServeHTTP(&wrappedWriter, r)

			end := time.Since(start)

			logger.LogAttrs(
				r.Context(),
				slog.LevelInfo.Level(),
				"processed request",
				slog.Int("status", wrappedWriter.statusCode),
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("url", r.URL.Path),
				slog.String("agent", r.UserAgent()),
				slog.String("duration", end.String()),
			)
		})
	}
}
