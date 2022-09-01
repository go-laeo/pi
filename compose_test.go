package ezy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompose(t *testing.T) {
	one := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("1"))
			return next(ctx, p)
		}
	}

	two := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("2"))
			return next(ctx, p)
		}
	}

	three := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("3"))
			return next(ctx, p)
		}
	}

	var h HandlerFunc[Void] = func(ctx Context, p *Void) error {
		_, err := ctx.Write([]byte("h"))
		return err
	}

	t.Run("(Composer).With should returns new slice", func(t *testing.T) {
		com := Compose(one, two)
		com2 := com.With(three)
		if len(com2) != 3 {
			t.Errorf("com2 should contains 5 elements, got = %d", len(com2))
		}
		if &com[0] == &com2[0] {
			t.Errorf(".With should returns new slice")
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		com2.For(h).ServeHTTP(w, r)

		if w.Body.String() != "123h" {
			t.Errorf(".With should append middlewares to the end, want = 123h, got = %s", w.Body.String())
		}
	})

	t.Run("empty composer should works normally", func(t *testing.T) {
		com := Compose[Void]()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		com.For(h).ServeHTTP(w, r)

		if w.Body.String() != "h" {
			t.Errorf("test empty composer want = h, got = %s", w.Body.String())
		}
	})

	t.Run("one composer shhould works", func(t *testing.T) {
		com := Compose(one)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		com.For(h).ServeHTTP(w, r)

		if w.Body.String() != "1h" {
			t.Errorf("test one composer want = 1h, got = %s", w.Body.String())
		}
	})
}

func BenchmarkComposerWith(b *testing.B) {
	one := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("1"))
			return next(ctx, p)
		}
	}

	two := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("2"))
			return next(ctx, p)
		}
	}

	three := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("3"))
			return next(ctx, p)
		}
	}

	com := Compose(one, two)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		com.With(three)
	}
}

func BenchmarkComposerFor(b *testing.B) {
	one := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("1"))
			return next(ctx, p)
		}
	}

	two := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("2"))
			return next(ctx, p)
		}
	}

	three := func(next HandlerFunc[Void]) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			ctx.Write([]byte("3"))
			return next(ctx, p)
		}
	}

	var h HandlerFunc[Void] = func(ctx Context, p *Void) error {
		_, err := ctx.Write([]byte("h"))
		return err
	}

	com := Compose(one, two, three)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		com.For(h)
	}
}
