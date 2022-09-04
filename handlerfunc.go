package pi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Void struct{}

type HandlerFunc[T any] func(ctx Context, p *T) error

func (h HandlerFunc[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := new(T)
	var err error

	if r.ContentLength != 0 {
		if _, ok := any(input).(*Void); !ok { // skip decode of type *Void
			err = json.NewDecoder(r.Body).Decode(input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(NewError(http.StatusBadRequest, err.Error()))
				return
			}

			if v, ok := any(input).(Validator); ok {
				if err = v.Validate(r.Context()); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(NewError(http.StatusBadRequest, err.Error()))
					return
				}
			}
		}
	}

	var p url.Values
	if v := r.Context().Value(&routePathParam{}); v != nil {
		p = v.(url.Values)
	}

	c := &_ctx{w: w, r: r, b: &bytes.Buffer{}, c: 200, p: p}
	if err = h(c, input); err != nil {
		switch v := err.(type) {
		case *Error:
			c.WriteHeader(v.Code)
			c.Json(v)
		default:
			c.WriteHeader(http.StatusInternalServerError)
			c.Json(NewError(http.StatusInternalServerError, err.Error()))
		}
	}

	w.WriteHeader(c.c)
	if _, err := io.Copy(w, c.b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewError(http.StatusInternalServerError, err.Error()))
		return
	}
}
