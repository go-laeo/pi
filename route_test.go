package ezy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteInsert(t *testing.T) {
	type fields struct {
		sub     map[string]*Route
		pattern string
	}
	type args struct {
		h      HandlerFunc[any]
		route  string
		method string
	}
	tests := []struct {
		want   *Route
		args   args
		fields fields
		name   string
	}{
		{
			name: "insert should succeed",
			fields: fields{
				sub: make(map[string]*Route),
			},
			args: args{
				route:  "/api/v1/users",
				method: "GET",
				h: func(ctx Context, p *any) error {
					return nil
				},
			},
			want: &Route{
				pattern: "users",
				sub:     make(map[string]*Route),
			},
		},
		{
			name: "insert same path should succeed",
			fields: fields{
				sub: make(map[string]*Route),
			},
			args: args{
				route:  "/api/v1/users",
				method: "GET",
				h: func(ctx Context, p *any) error {
					return nil
				},
			},
			want: &Route{
				pattern: "users",
				sub:     make(map[string]*Route),
			},
		},
		{
			name: "insert at same part should succeed",
			fields: fields{
				sub: make(map[string]*Route),
			},
			args: args{
				route:  "/api",
				method: "GET",
				h: func(ctx Context, p *any) error {
					return nil
				},
			},
			want: &Route{
				pattern: "api",
				sub:     make(map[string]*Route),
			},
		},
		{
			name: "empty string should not panic",
			fields: fields{
				sub: make(map[string]*Route),
			},
			args: args{
				route:  "/",
				method: "GET",
			},
			want: &Route{
				pattern: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Route{
				pattern: tt.fields.pattern,
				sub:     tt.fields.sub,
			}
			if got := p.Insert(tt.args.route, tt.args.method, tt.args.h); got.pattern != tt.want.pattern || len(got.sub) > 0 {
				t.Fatalf("node.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteSpecialInsert(t *testing.T) {
	t.Run("insert dynamic route should succeed", func(t *testing.T) {
		ro := &Route{}
		node := ro.Insert("/api/v1/users/:id", "GET", nil)
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
		ro := &Route{}
		node := ro.Insert("/api/v1/users/:id/posts/:po", "GET", nil)
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
		ro := &Route{}
		node := ro.Insert("/api/v1/users/:id/:po", "GET", nil)
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
		ro := &Route{}
		node := ro.Insert("/api/v1/users/*id", "GET", nil)
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
	root := &Route{
		sub: make(map[string]*Route),
	}
	for i := 0; i < b.N; i++ {
		root.Insert("/api/v1/users", "GET", nil)
	}
}

func TestRouteSearch(t *testing.T) {
	gen := func(b string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(b))
		}
	}

	root := &Route{}
	root.Insert("/api", "GET", gen("/api"))
	root.Insert("/api/", "GET", gen("/api/"))
	root.Insert("/api/v1/users", "GET", gen("/api/v1/users"))
	root.Insert("/api/v1/users/:id", "GET", gen("/api/v1/users/:id"))
	root.Insert("/api/v1/users/:id/posts", "GET", gen("/api/v1/users/:id/posts"))
	root.Insert("/api/v1/users/:id/posts/:po", "GET", gen("/api/v1/users/:id/posts/:po"))
	root.Insert("/*path", "GET", gen("/*path"))
	root.Insert("*path", "GET", gen("*path"))
	root.Insert("/uploads/*path", "GET", gen("/uploads/*path"))

	type args struct {
		route string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "search static route should succeed",
			args: args{route: "/api/v1/users"},
			want: "/api/v1/users",
		},
		{
			name: "search dynamic route should succeed",
			args: args{route: "/api/v1/users/102"},
			want: "/api/v1/users/:id",
		},
		{
			name: "search static route next to dynamic one should succeed",
			args: args{route: "/api/v1/users/101/posts"},
			want: "/api/v1/users/:id/posts",
		},
		{
			name: "search dynamic route nested in another one should succeed",
			args: args{route: "/api/v1/users/101/posts/120"},
			want: "/api/v1/users/:id/posts/:po",
		},
		{
			name: "search wildcard route should succeed",
			args: args{route: "/pathtowildcard"},
			want: "/*path",
		},
		{
			name: "search wildcard route nested in static one should succeed",
			args: args{route: "/uploads/users/1.avatar.png"},
			want: "/uploads/*path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := root.Search(tt.args.route)
			if got == nil {
				t.Fatalf("%s does not find corresponding route", tt.name)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.args.route, nil)
			h, ok := got.hmap[http.MethodGet]
			if !ok {
				t.Fatalf("[%s] does not define http.Handler", tt.name)
			}

			h.ServeHTTP(w, r)
			if tt.want != w.Body.String() {
				t.Fatalf("[%s] want = %s, got = %s", tt.name, tt.want, w.Body.String())
			}
		})
	}
}

func BenchmarkRouteSearch(b *testing.B) {
	gen := func(b string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(b))
		}
	}

	root := &Route{}
	root.Insert("/api", "GET", gen("/api"))
	root.Insert("/api/", "GET", gen("/api/"))
	root.Insert("/api/v1/users", "GET", gen("/api/v1/users"))
	root.Insert("/api/v1/users/:id", "GET", gen("/api/v1/users/:id"))
	root.Insert("/api/v1/users/:id/posts", "GET", gen("/api/v1/users/:id/posts"))
	root.Insert("/api/v1/users/:id/posts/:po", "GET", gen("/api/v1/users/:id/posts/:po"))
	root.Insert("/*path", "GET", gen("/*path"))
	root.Insert("*path", "GET", gen("*path"))
	root.Insert("/uploads/*path", "GET", gen("/uploads/*path"))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root.Search("/api/v1/users/:id/posts/:po")
	}
}
