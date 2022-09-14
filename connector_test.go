package pi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnect(t *testing.T) {
	t.Run("Connect_connectors_should_run_as_middleware", func(t *testing.T) {
		var c = func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Header().Set("X-HTTP-SERVER", "net/http")
				return next(ctx)
			}
		}

		var h HandlerFunc = func(ctx Context) error {
			return ctx.Text("OK")
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		h.Connect(c).ServeHTTP(w, r)

		if w.Body.String() != "OK" {
			t.Fatalf("test want = OK, got = %s", w.Body.String())
		}
		if w.Header().Get("X-HTTP-SERVER") != "net/http" {
			t.Fatalf("header want = net/http, got = %s", w.Header().Get("X-HTTP-SERVER"))
		}
	})
}

func BenchmarkConnect(b *testing.B) {
	var c = func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			ctx.Header().Set("X-HTTP-SERVER", "net/http")
			return next(ctx)
		}
	}

	var h HandlerFunc = func(ctx Context) error {
		return ctx.Text("OK")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h = h.Connect(c)
	}
}
