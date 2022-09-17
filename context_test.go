package pi

import (
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

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		c := createContext(w, r, nil, nil)
		c.SetContext(context.WithValue(c.Context(), _a, "a"))
		c.SetContext(context.WithValue(c.Context(), _b, "b"))
		c.Context().Value(_a)
	})
}
