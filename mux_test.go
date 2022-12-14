package pi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerMux_ServeHTTP(t *testing.T) {
	gen := func(b string) HandlerFunc {
		return func(ctx Context) error {
			return ctx.Text(b)
		}
	}

	sm := NewServerMux()
	sm.Route("/api/v1/users").Get(gen("OK"))
	sm.Route("/api/v1/system/status").Any(gen("OK"))
	sm.Route("/preflight").Options(gen("OK"))

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		args       args
		want       []byte
		wantStatus int
	}{
		{
			name: "should got OK",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/api/v1/users", nil),
			},
			want:       []byte("OK"),
			wantStatus: 200,
		},
		{
			name: "should respond with status 404",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/api", nil),
			},
			want:       nil,
			wantStatus: 404,
		},
		{
			name: "should reach status with POST method",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/api/v1/system/status", nil),
			},
			want:       []byte("OK"),
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm.ServeHTTP(tt.args.w, tt.args.r)
			if !bytes.Equal(tt.args.w.Body.Bytes(), tt.want) || tt.args.w.Code != tt.wantStatus {
				t.Errorf("%s want = %s, got = %s", tt.name, tt.want, tt.args.w.Body.Bytes())
			}
		})
	}
}

func TestServerMux_Group(t *testing.T) {
	sm := NewServerMux()
	sm.Group("/api/v1", func(sm ServerMux) {
		sm.Route("/users").Get(HandlerFunc(func(ctx Context) error {
			return ctx.Text("API")
		}))
	})
	sm.Route("/users").Get(HandlerFunc(func(ctx Context) error {
		return ctx.Text("OK")
	}))

	t.Run("request sub-route should succeed", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		w := httptest.NewRecorder()
		sm.ServeHTTP(w, r)
		if !bytes.Equal(w.Body.Bytes(), []byte("API")) {
			t.Errorf("request sub-route want = API, got = %s", w.Body.String())
		}
	})
	t.Run("request non-groupped route should succeed", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		sm.ServeHTTP(w, r)
		if !bytes.Equal(w.Body.Bytes(), []byte("OK")) {
			t.Errorf("request non-groupped route want = API, got = %s", w.Body.String())
		}
	})
}

func TestServerMux_PathParamCapture(t *testing.T) {
	gen := func(b string) HandlerFunc {
		return func(ctx Context) error {
			return ctx.Text(ctx.Param(b))
		}
	}

	sm := NewServerMux()
	sm.Route("/api/v1/users/:id").Get(gen("id"))
	sm.Route("/api/v1/users/:id/posts").Get(gen("id"))
	sm.Route("/api/v1/users/:id/posts/:po").Get(gen("po"))
	sm.Route("/*path").Get(gen("path"))
	sm.Route("*path").Get(gen("path"))
	sm.Route("/uploads/*path").Get(gen("path"))

	t.Run("capture path param should succeed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/users/100", nil)
		sm.ServeHTTP(w, r)
		if w.Body.String() != "100" {
			t.Fatalf("path param want = 100, got = %s", w.Body.String())
		}
	})

	t.Run("capture nested path param should succeed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/users/100/posts/101", nil)
		sm.ServeHTTP(w, r)
		if w.Body.String() != "101" {
			t.Fatalf("path param want = 101, got = %s", w.Body.String())
		}
	})

	t.Run("capture wildcard param should succeed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/helloworld", nil)
		sm.ServeHTTP(w, r)
		if w.Body.String() != "helloworld" {
			t.Fatalf("wildcard param want = helloworld, got = %s", w.Body.String())
		}
	})

	t.Run("capture nested wildcard param should succeed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/uploads/users/1.avatar.png", nil)
		sm.ServeHTTP(w, r)
		if w.Body.String() != "users/1.avatar.png" {
			t.Fatalf("path param want = users/1.avatar.png, got = %s", w.Body.String())
		}
	})
}

func TestServerMux_Use(t *testing.T) {
	sm := NewServerMux()
	sm.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			ctx.Header().Set("X-TEST-HEADER", "TEST")
			return next(ctx)
		}
	})
	sm.Group("/api/v1", func(sm ServerMux) {
		sm.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Header().Set("X-API-VERSION", "v1")
				return next(ctx)
			}
		})
		sm.Route("/users").Get(HandlerFunc(func(ctx Context) error {
			return ctx.Text("API")
		}))
	})
	sm.Route("/api/v1/users").Post(func(ctx Context) error {
		return ctx.Text("PUBLIC")
	})
	sm.Route("/users").Get(HandlerFunc(func(ctx Context) error {
		return ctx.Text("OK")
	}))

	t.Run("request /users should got right header value", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		sm.ServeHTTP(w, r)
		if v := w.Header().Get("X-TEST-HEADER"); v != "TEST" {
			t.Fatalf("response from /users should contains X-TEST-HEADER = TEST, got = %s", v)
		}
	})

	t.Run("request /api/v1/users should got right header value", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		sm.ServeHTTP(w, r)
		if v := w.Header().Get("X-TEST-HEADER"); v != "TEST" {
			t.Fatalf("response from /api/v1/users should contains X-TEST-HEADER = TEST, got = %s", v)
		}
	})

	t.Run("request /notfound should not contains header value", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		sm.ServeHTTP(w, r)
		if v := w.Header().Get("X-TEST-HEADER"); v != "" {
			t.Fatalf("response from /notfound should not contains X-TEST-HEADER, got = %s", v)
		}
	})

	t.Run("request POST /api/v1/users should not contains header X-API-VERSION", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
		sm.ServeHTTP(w, r)
		if v := w.Header().Get("X-API-VERSION"); v != "" {
			t.Fatalf("response from /notfound should not contains X-TEST-HEADER, got = %s", v)
		}
	})
}

func BenchmarkServerMux_ServeHTTP(b *testing.B) {
	gen := func(b string) HandlerFunc {
		return func(ctx Context) error {
			return ctx.Text(ctx.Param(b))
		}
	}
	sm := NewServerMux()
	sm.Route("/api/v1/users/:id/posts/:po").Get(gen("po"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/100/posts/101", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.ServeHTTP(w, r)
	}
}
