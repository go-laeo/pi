package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-laeo/ezy"
)

func main() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	process().Connect(logging, cors).ServeHTTP(w, r)
	if w.Body.String() != "hello, world!" {
		panic("unexpected response body!")
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		panic("unexpected response header!")
	}
	if w.Header().Get("Access-Control-Methods") != "" {
		panic("unexpected response header!!")
	}

	println("test ok!")
}

func process() ezy.HandlerFunc[ezy.Void] {
	return func(ctx ezy.Context, p *ezy.Void) error {
		return ctx.Text("hello, world!")
	}
}

func logging(next ezy.HandlerFunc[ezy.Void]) ezy.HandlerFunc[ezy.Void] {
	return func(ctx ezy.Context, p *ezy.Void) error {
		log.SetPrefix(fmt.Sprintf("[%s] ", ctx.IP()))
		log.Println(ctx.Method(), ctx.URL().Path)

		return next(ctx, p)
	}
}

func cors(next ezy.HandlerFunc[ezy.Void]) ezy.HandlerFunc[ezy.Void] {
	return func(ctx ezy.Context, p *ezy.Void) error {
		ctx.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Is(http.MethodOptions) {
			ctx.Header().Set("Access-Control-Methods", "POST, PUT, PATCH, DELETE")
			ctx.Header().Set("Access-Control-Allow-Headers", "*")
			ctx.Header().Set("Access-Control-Max-Age", "86400")
			ctx.WriteHeader(http.StatusNoContent)
			return nil
		}

		return next(ctx, p)
	}
}
