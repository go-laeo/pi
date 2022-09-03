package ezy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnect(t *testing.T) {
	t.Run("Connect_connectors_should_run_as_middleware", func(t *testing.T) {
		var c Connector[Void] = func(next HandlerFunc[Void]) HandlerFunc[Void] {
			return func(ctx Context, p *Void) error {
				ctx.Header().Set("X-HTTP-SERVER", "net/http")
				return next(ctx, p)
			}
		}

		var h HandlerFunc[Void] = func(ctx Context, _ *Void) error {
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
	var c Connector[Void] = func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Header().Set("X-HTTP-SERVER", "net/http")
			return next(ctx, p)
		}
	}

	var h HandlerFunc[Void] = func(ctx Context, _ *Void) error {
		return ctx.Text("OK")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h = h.Connect(c)
	}
}
