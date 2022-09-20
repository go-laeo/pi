package pi

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_decode(t *testing.T) {
	type embeded struct {
		PageIndex int `query:"pi"`
		PageSize  int `query:"ps"`
	}

	type usr struct {
		Name     string `query:"name"`
		password string
		Date     string
		embeded
		Embeded embeded
		ID      int `query:"id"`
		Age     int `query:"-"`
		Untyped struct {
			PageIndex int `query:"pi"`
			PageSize  int `query:"ps"`
		}
		Untyped2 struct {
			Index int `query:"pi"`
			Size  int `query:"ps"`
		}
	}

	m := url.Values{
		"id":       []string{"1"},
		"name":     []string{"lebai"},
		"password": []string{"secured"},
		"Age":      []string{"18"},
		"Date":     []string{"2022"},
		"pi":       []string{"1"},
		"ps":       []string{"20"},
	}

	v := usr{}
	err := decode(m, reflect.ValueOf(&v))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if v.ID != 1 {
		t.Fatalf("id != 1, = %d", v.ID)
	}
	if v.Name != "lebai" {
		t.Fatalf("name != lebai = %s", v.Name)
	}
	if v.password != "" {
		t.Fatalf("password isnt empty")
	}
	if v.Age != 0 {
		t.Fatalf("age != 0 = %d", v.Age)
	}
	if v.Date != "2022" {
		t.Fatalf("date != 2022 = %s", v.Date)
	}
	if v.PageIndex != 1 {
		t.Fatalf("page index != 1 = %d", v.PageIndex)
	}
	if v.PageSize != 20 {
		t.Fatalf("page size != 20 = %d", v.PageSize)
	}
	if v.Embeded.PageIndex != 1 {
		t.Fatalf("embeded page index != 1 = %d", v.Embeded.PageIndex)
	}
	if v.Embeded.PageSize != 20 {
		t.Fatalf("embeded page size != 20 = %d", v.Embeded.PageSize)
	}
	if v.Untyped.PageIndex != 1 {
		t.Fatalf("untyped page index != 1 = %d", v.Untyped.PageIndex)
	}
	if v.Untyped.PageSize != 20 {
		t.Fatalf("untyped page size != 20 = %d", v.Untyped.PageSize)
	}
	if v.Untyped2.Index != 1 {
		t.Fatalf("untyped2 index != 1 = %d", v.Untyped2.Index)
	}
	if v.Untyped2.Size != 20 {
		t.Fatalf("untyped2 size != 20 = %d", v.Untyped2.Size)
	}
}

func BenchmarkDeocde(b *testing.B) {
	type usr struct {
		Name string `query:"name"`
		Date string
		ID   int `query:"id"`
		Age  int `query:"-"`
	}

	m := url.Values{
		"id":   []string{"1"},
		"name": []string{"lebai"},
		"Age":  []string{"18"},
		"Date": []string{"2022"},
	}
	v := usr{}
	p := reflect.ValueOf(&v)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decode(m, p)
	}
}
