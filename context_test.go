package ezy

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ctx_WithContext(t *testing.T) {
	t.Run("should not throw goroutine stack exceeds", func(_ *testing.T) {
		type a struct{}
		var _a *a
		type b struct{}
		var _b *b

		var c Context = &_ctx{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/", nil),
			b: bytes.NewBuffer(nil),
			c: 200,
		}
		c.SetContext(context.WithValue(c.Context(), _a, "a"))
		c.SetContext(context.WithValue(c.Context(), _b, "b"))
		c.Context().Value(_a)
	})
}
