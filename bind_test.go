package pi

import (
	"net/url"
	"testing"
)

func TestBind(t *testing.T) {
	type usr struct {
		Name string `query:"name"`
		ID   int
	}

	v := url.Values{"ID": []string{"1"}, "name": []string{"Go"}}
	p := usr{}
	err := Bind(v, &p)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if p.ID != 1 {
		t.Fatalf("ID != 1 = %d", p.ID)
	}
	if p.Name != "Go" {
		t.Fatalf("Name != Go = %s", p.Name)
	}
}

func BenchmarkBind(b *testing.B) {
	type usr struct {
		Name string `query:"name"`
		ID   int
	}

	v := url.Values{"ID": []string{"1"}, "name": []string{"Go"}}
	p := usr{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bind(v, &p)
	}
}
