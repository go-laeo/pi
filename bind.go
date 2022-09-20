package pi

import (
	"net/url"
	"reflect"
)

// Bind populate values from v and fills to p.
//
//	type Paging struct{
//		Index int `query:"i"`
//		Size int `query:"s"`
//	}
//
//	p := Paging{}
//	err := pi.Bind(ctx.URL().Query(), &p)
//
// The above code shows us how we can mapping
// queries to a struct by simply call `pi.Bind()`.
func Bind[T any](v url.Values, p *T) error {
	return decode(v, reflect.ValueOf(p))
}
