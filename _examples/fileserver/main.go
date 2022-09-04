package main

import (
	"context"
	"embed"
	"net/http"

	"github.com/go-laeo/pi"
)

//go:embed web
var stubs embed.FS

func main() {
	sm := pi.NewServerMux(context.Background())
	sm.Get("/api/v1/users", pi.HandlerFunc[pi.Void](func(ctx pi.Context, p *pi.Void) error {
		return ctx.Text("users")
	}))
	sm.Any("/*fs", http.FileServer(http.FS(stubs)))
	http.ListenAndServe("localhost:8080", sm)
}
