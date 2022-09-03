package ezy

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {
	gen := func(b string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(b))
		}
	}

	sm := NewServerMux(context.Background())
	sm.Get("/api/v1/users", gen("OK"))
	sm.Any("/api/v1/system/status", gen("OK"))
	sm.Options("/preflight", gen("OK"))

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
			want:       []byte("not found"),
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
	sm := NewServerMux(context.Background())
	sm.Group("/api/v1", func(sm *ServerMux) {
		sm.Get("/users", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("API"))
		}))
	})
	sm.Get("/users", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
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
	gen := func(b string) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			return ctx.Text(ctx.Param(b))
		}
	}

	sm := NewServerMux(context.Background())
	sm.Get("/api/v1/users/:id", gen("id"))
	sm.Get("/api/v1/users/:id/posts", gen("id"))
	sm.Get("/api/v1/users/:id/posts/:po", gen("po"))
	sm.Get("/*path", gen("path"))
	sm.Get("*path", gen("path"))
	sm.Get("/uploads/*path", gen("path"))

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

func BenchmarkServerMux_ServeHTTP(b *testing.B) {
	gen := func(b string) HandlerFunc[Void] {
		return func(ctx Context, p *Void) error {
			return ctx.Text(ctx.Param(b))
		}
	}
	sm := NewServerMux(context.Background())
	sm.Get("/api/v1/users/:id/posts/:po", gen("po"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/100/posts/101", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.ServeHTTP(w, r)
	}
}
