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
	sm.Get("/api/v1/users", pi.HandlerFunc(func(ctx pi.Context) error {
		return ctx.Text("users")
	}))
	sm.Any("/*fs", pi.FileServer(http.FS(stubs), "web/index.html"))
	http.ListenAndServe("localhost:8080", sm)
}
