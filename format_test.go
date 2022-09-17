package pi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormat(t *testing.T) {
	type test struct {
		A string
		B int
	}
	raw := `{"A":"AA","B":1}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(raw))
	ctx := createContext(w, r, nil, nil)

	p := &test{}
	err := Format(ctx, p)
	if err != nil {
		t.Fatalf("TestFormat got error = %v", err)
	}
	if p.A != "AA" {
		t.Fatalf("TestFormat want A = AA, got = %s", p.A)
	}
	if p.B != 1 {
		t.Fatalf("TestFormat want B = 1, got = %d", p.B)
	}
}

func BenchmarkFormat(b *testing.B) {
	type test struct {
		A string
		B int
	}
	raw := `{"A":"AA","B":1}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(raw))
	ctx := createContext(w, r, nil, nil)

	p := &test{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(ctx, p)
	}
}
