package pi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type HandlerFunc func(ctx Context) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p url.Values
	if v := r.Context().Value(_routePathParam); v != nil {
		p = v.(url.Values)
	}

	c := &_ctx{w: w, r: r, b: &bytes.Buffer{}, c: 200, p: p}
	if err := h(c); err != nil {
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
