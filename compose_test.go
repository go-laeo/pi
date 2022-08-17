package ezy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type greeting struct {
	Name string
}

func (g *greeting) Validate(ctx context.Context) error {
	if g.Name != "bbb" {
		return errors.New("validate failed")
	}
	return nil
}

func TestCompose(t *testing.T) {
	h := Compose(func(ctx Context, p *greeting) error {
		return json.NewEncoder(ctx).Encode(NewError(200, "Hi, "+p.Name))
	})

	tests := []struct {
		name string
		body io.Reader
		want int
	}{
		{
			name: "process_should_succeed",
			body: bytes.NewBufferString(`{"Name":"bbb"}`),
			want: 200,
		},
		{
			name: "process_should_failed_with_400",
			body: bytes.NewBufferString(`{"Name": "aaa"}`),
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/greeting", tt.body)
			h.ServeHTTP(w, r)
			if w.Result().StatusCode != tt.want {
				t.Errorf("%s want = %d, got = %d", tt.name, tt.want, w.Result().StatusCode)
			}
		})
	}
}
