package ezy

import (
	"testing"
)

func Test_node_Insert(t *testing.T) {
	type fields struct {
		lit string
		sub map[string]*Route
	}
	type args struct {
		route  string
		method string
		h      HandlerFunc[any]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Route
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
				lit: "users",
				sub: make(map[string]*Route),
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
				lit: "users",
				sub: make(map[string]*Route),
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
				lit: "api",
				sub: make(map[string]*Route),
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
				lit: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Route{
				lit: tt.fields.lit,
				sub: tt.fields.sub,
			}
			if got := p.Insert(tt.args.route, tt.args.method, tt.args.h); got.lit != tt.want.lit || len(got.sub) > 0 {
				t.Errorf("node.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNodeInsert(b *testing.B) {
	root := &Route{
		sub: make(map[string]*Route),
	}
	for i := 0; i < b.N; i++ {
		root.Insert("/api/v1/users", "GET", nil)
	}
}

func Test_node_Search(t *testing.T) {
	root := &Route{}
	root.Insert("/api", "GET", nil)
	root.Insert("/api/", "GET", nil)
	root.Insert("/api/v1/users", "GET", nil)
	root.Insert("/api/v1/users/:id", "GET", nil)

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
			want: "users",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := root.Search(tt.args.route); got.lit != tt.want {
				t.Errorf("node.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_node_Search(b *testing.B) {
	root := &Route{}
	root.Insert("/api", "GET", nil)
	root.Insert("/api/", "GET", nil)
	root.Insert("/api/v1/users", "GET", nil)
	root.Insert("/api/v1/users/:id", "GET", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root.Search("/api/v1/users")
	}
}
