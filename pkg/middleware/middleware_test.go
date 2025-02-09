package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrap(t *testing.T) {
	middleware := func(name string) func(http.Handler) http.Handler {
		// Return wrapped
		return func(h http.Handler) http.Handler {
			// Wrap a handler function so that it looks like a handler
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, name)
				h.ServeHTTP(w, r)
				fmt.Fprint(w, name)
			})
		}
	}

	rs := httptest.NewRecorder()
	rq := httptest.NewRequest(http.MethodGet, "/", nil)
	middlewares := [](func(http.Handler) http.Handler){
		middleware("|one|"),
		middleware("|two|"),
		middleware("|three|"),
		middleware("|four|"),
		middleware("|five|"),
	}
	want := "|one||two||three||four||five|base|five||four||three||two||one|"

	wrapper := Wrap(middlewares...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "base")
	}))

	wrapper.ServeHTTP(rs, rq)

	if rs.Body.String() != want {
		t.Errorf("middlewares not applied in correct order, got %q but want %q", rs.Body.String(), want)
	}

}
