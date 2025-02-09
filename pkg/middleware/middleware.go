package middleware

import (
	"net/http"
)

// Wraps a HTTP handler in the provided middlewares like a stack, with the last middleware
// added being the closest to the base handler, and the first added being the first/last
// handler wrapping called.
//
//	Usage:
//	handler = Wrap(RequestIDMiddleware(), LoggingMiddleware())(handler)
func Wrap(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {

	// Return the expected outer signature, which is the function that can be applied to a handler
	return func(h http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			// Wrap our middleware like a stack
			h = middlewares[i](h)
		}

		// Finally return our wrapped handler
		return h
	}
}
