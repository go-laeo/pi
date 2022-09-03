package main

import (
	"context"
	"embed"
	"net/http"

	"github.com/go-laeo/ezy"
)

//go:embed web
var stubs embed.FS

func main() {
	sm := ezy.NewServerMux(context.Background())
	sm.Get("/api/v1/users", ezy.HandlerFunc[ezy.Void](func(ctx ezy.Context, p *ezy.Void) error {
		return ctx.Text("users")
	}))
	sm.Any("/*fs", http.FileServer(http.FS(stubs)))
	http.ListenAndServe("localhost:8080", sm)
}
