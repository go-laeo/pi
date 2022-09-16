package pi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRouteInsert(t *testing.T) {
	type fields struct {
		sub     map[string]*_route
		pattern string
	}
	type args struct {
		h     HandlerFunc
		route string
	}
	tests := []struct {
		want   *_route
		args   args
		fields fields
		name   string
	}{
		{
			name: "insert should succeed",
			fields: fields{
				sub: make(map[string]*_route),
			},
			args: args{
				route: "/api/v1/users",
				h: func(ctx Context) error {
					return nil
				},
			},
			want: &_route{
				pattern: "users",
				sub:     make(map[string]*_route),
			},
		},
		{
			name: "insert same path should succeed",
			fields: fields{
				sub: make(map[string]*_route),
			},
			args: args{
				route: "/api/v1/users",
				h: func(ctx Context) error {
					return nil
				},
			},
			want: &_route{
				pattern: "users",
				sub:     make(map[string]*_route),
			},
		},
		{
			name: "insert at same part should succeed",
			fields: fields{
				sub: make(map[string]*_route),
			},
			args: args{
				route: "/api",
				h: func(ctx Context) error {
					return nil
				},
			},
			want: &_route{
				pattern: "api",
				sub:     make(map[string]*_route),
			},
		},
		{
			name: "empty string should not panic",
			fields: fields{
				sub: make(map[string]*_route),
			},
			args: args{
				route: "/",
			},
			want: &_route{
				pattern: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &_route{
				pattern: tt.fields.pattern,
				sub:     tt.fields.sub,
			}
			if got := p.Insert(tt.args.route).Get(tt.args.h).(*_route); got.pattern != tt.want.pattern || len(got.sub) > 0 {
				t.Fatalf("node.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteSpecialInsert(t *testing.T) {
	t.Run("insert dynamic route should succeed", func(t *testing.T) {
		ro := &_route{}
		node := ro.Insert("/api/v1/users/:id").Get(nil).(*_route)
		if node.pattern != ":id" {
			t.Fatalf("dynamic route pattern want = :id, got = %s", node.pattern)
		}
		if node.placeholder != "id" {
			t.Fatalf("dynamic route placeholder want = id, got = %s", node.placeholder)
		}
		if !node.parent.hasDynamicChild {
			t.Fatalf("dynamic route should set hasDynamicChild on parent, got = %v", node.parent.hasDynamicChild)
		}
	})

	t.Run("insert nested dynamic route should succeed", func(t *testing.T) {
		ro := &_route{}
		node := ro.Insert("/api/v1/users/:id/posts/:po").Get(nil).(*_route)
		if node.pattern != ":po" {
			t.Fatalf("dynamic route pattern want = :po, got = %s", node.pattern)
		}
		if node.placeholder != "po" {
			t.Fatalf("dynamic route placeholder want = po, got = %s", node.placeholder)
		}
		if !node.parent.hasDynamicChild {
			t.Fatalf("dynamic route should set hasDynamicChild on parent, got = %v", node.parent.hasDynamicChild)
		}
	})

	t.Run("insert dynamic route next to dynamic one should succeed", func(t *testing.T) {
		ro := &_route{}
		node := ro.Insert("/api/v1/users/:id/:po").Get(nil).(*_route)
		if node.pattern != ":po" {
			t.Fatalf("dynamic route pattern want = :po, got = %s", node.pattern)
		}
		if node.placeholder != "po" {
			t.Fatalf("dynamic route placeholder want = po, got = %s", node.placeholder)
		}
		if !node.parent.hasDynamicChild {
			t.Fatalf("dynamic route should set hasDynamicChild on parent, got = %v", node.parent.hasDynamicChild)
		}
	})

	t.Run("insert wildcard route should succeed", func(t *testing.T) {
		ro := &_route{}
		node := ro.Insert("/api/v1/users/*id").Get(nil).(*_route)
		if node.pattern != "*id" {
			t.Fatalf("wildcard route pattern want = :id, got = %s", node.pattern)
		}
		if node.placeholder != "id" {
			t.Fatalf("wildcard route placeholder want = id, got = %s", node.placeholder)
		}
		if !node.parent.hasWildcardChild {
			t.Fatalf("wildcard route should set hasWildcardChild on parent, got = %v", node.parent.hasWildcardChild)
		}
	})
}

func BenchmarkNodeInsert(b *testing.B) {
	root := &_route{
		sub: make(map[string]*_route),
	}
	for i := 0; i < b.N; i++ {
		root.Insert("/api/v1/users").Get(nil).Post(nil)
	}
}

func TestRouteSearch(t *testing.T) {
	gen := func(b string) HandlerFunc {
		return func(ctx Context) error {
			return ctx.Text(b)
		}
	}

	root := &_route{}
	root.Insert("/api").Get(gen("/api"))
	root.Insert("/api/").Get(gen("/api/"))
	root.Insert("/api/v1/users").Get(gen("/api/v1/users"))
	root.Insert("/api/v1/users/admin/share").Get(gen("/api/v1/users/admin/share"))
	root.Insert("/api/v1/users/:id").Get(gen("/api/v1/users/:id"))
	root.Insert("/api/v1/users/:id/posts").Get(gen("/api/v1/users/:id/posts"))
	root.Insert("/api/v1/users/:id/posts/:po").Get(gen("/api/v1/users/:id/posts/:po"))
	root.Insert("/*path").Get(gen("/*path"))
	root.Insert("*path").Get(gen("*path"))
	root.Insert("/uploads/*path").Get(gen("/uploads/*path"))

	type args struct {
		route string
	}
	tests := []struct {
		captured map[string]string
		name     string
		args     args
		want     string
	}{
		{
			name: "search static route should succeed",
			args: args{route: "/api/v1/users"},
			want: "/api/v1/users",
		},
		{
			name:     "search dynamic route should succeed",
			args:     args{route: "/api/v1/users/102"},
			want:     "/api/v1/users/:id",
			captured: map[string]string{"id": "102"},
		},
		{
			name:     "search static route next to dynamic one should succeed",
			args:     args{route: "/api/v1/users/101/posts"},
			want:     "/api/v1/users/:id/posts",
			captured: map[string]string{"id": "101"},
		},
		{
			name:     "search dynamic route nested in another one should succeed",
			args:     args{route: "/api/v1/users/101/posts/120"},
			want:     "/api/v1/users/:id/posts/:po",
			captured: map[string]string{"id": "101", "po": "120"},
		},
		{
			name:     "search wildcard route should succeed",
			args:     args{route: "/pathtowildcard"},
			want:     "/*path",
			captured: map[string]string{"path": "pathtowildcard"},
		},
		{
			name:     "search wildcard route nested in static one should succeed",
			args:     args{route: "/uploads/users/1.avatar.png"},
			want:     "/uploads/*path",
			captured: map[string]string{"path": "users/1.avatar.png"},
		},
		{
			name: "search static first",
			args: args{route: "/api/v1/users/admin/share"},
			want: "/api/v1/users/admin/share",
		},
		{
			name:     "search static then dynamic",
			args:     args{route: "/api/v1/users/admin/posts"},
			want:     "/api/v1/users/:id/posts",
			captured: map[string]string{"id": "admin"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captured := make(url.Values)
			got := root.Search(tt.args.route, captured)
			if got == nil {
				t.Fatalf("%s does not find corresponding route", tt.name)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.args.route, nil)
			ok := got.Invoke(createContext(w, r, captured))
			if !ok {
				t.Fatalf("[%s] does not define http.Handler", tt.name)
			}

			if tt.want != w.Body.String() {
				t.Fatalf("[%s] want = %s, got = %s", tt.name, tt.want, w.Body.String())
			}
			if tt.captured != nil {
				for k, v := range tt.captured {
					if captured.Get(k) != v {
						t.Fatalf("[%s] want captured %s = %s, got %s", tt.name, k, v, captured.Get(k))
					}
				}
			}
		})
	}
}

func BenchmarkRouteSearch(b *testing.B) {
	gen := func(b string) HandlerFunc {
		return func(ctx Context) error {
			return ctx.Text(b)
		}
	}

	root := &_route{}
	root.Insert("/api").Get(gen("/api"))
	root.Insert("/api/").Get(gen("/api/"))
	root.Insert("/api/v1/users").Get(gen("/api/v1/users"))
	root.Insert("/api/v1/users/:id").Get(gen("/api/v1/users/:id"))
	root.Insert("/api/v1/users/:id/posts").Get(gen("/api/v1/users/:id/posts"))
	root.Insert("/api/v1/users/:id/posts/:po").Get(gen("/api/v1/users/:id/posts/:po"))
	root.Insert("/*path").Get(gen("/*path"))
	root.Insert("*path").Get(gen("*path"))
	root.Insert("/uploads/*path").Get(gen("/uploads/*path"))

	cap := make(url.Values)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root.Search("/api/v1/users/100/posts/101", cap)
	}
}
