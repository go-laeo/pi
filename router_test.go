package ezy

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {
	ro := NewRouter()
	ro.Get("/api/v1/users", Compose(func(ctx Context, p *any) error {
		ctx.Write([]byte("OK"))
		return nil
	}))
	ro.Any("/api/v1/system/status", Compose(func(ctx Context, p *Void) error {
		ctx.Write([]byte("OK"))
		return nil
	}))
	ro.Options("/preflight", Compose(func(ctx Context, p *Void) error {
		return errors.New("failed")
	}))

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
			ro.ServeHTTP(tt.args.w, tt.args.r)
			if !bytes.Equal(tt.args.w.Body.Bytes(), tt.want) || tt.args.w.Code != tt.wantStatus {
				t.Errorf("%s want = %s, got = %s", tt.name, tt.want, tt.args.w.Body.Bytes())
			}
		})
	}
}
