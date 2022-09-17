package main

import (
	"context"
	"net"
	"net/http"
)

func main() {
	ctx := context.Background() // custom context

	srv := http.Server{
		Addr: ":8080",
		BaseContext: func(l net.Listener) context.Context {
			return ctx // inject to *http.Request
		},
	}

	srv.ListenAndServe()
}
