package ezy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandlerFuncT_ServeHTTP(t *testing.T) {
	var fn HandlerFunc[Void] = func(ctx Context, p *Void) error {
		return NewError(404, "not found")
	}

	t.Run("detect *Error should succeed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		fn.ServeHTTP(w, r)
		if w.Code != 404 {
			t.Errorf("detect *Error want code = 404, got = %d", w.Code)
		}
	})
}

func Benchmark_HandlerFuncT_ServeHTTP(b *testing.B) {
	var fn HandlerFunc[Void] = func(ctx Context, p *Void) error {
		return NewError(200, "not found")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	for i := 0; i < b.N; i++ {
		fn.ServeHTTP(w, r)
	}
}
