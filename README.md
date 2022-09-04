# pi

![build.yaml](https://github.com/go-laeo/pi/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/pi.svg)](https://pkg.go.dev/github.com/go-laeo/pi)

A `pi` is your good helper to build a clean JSON API server using Golang Generics.

# Quick Start

```Go
var h pi.HandlerFunc[any] = func (ctx pi.Context, _ *any) error {
    return ctx.Text("hello, world!")
}

sm := pi.NewServerMux(context.Background())
sm.Get("/api/v1/users", h)

http.ListenAndServe("localhost:8080", sm)
```

# Install

```shell
go get -u github.com/go-laeo/pi
```

# Features

- [x] Fast routing, routes group, route params and wildcard route
- [x] `net/http` compatible
- [x] Auto decode HTTP body using Generics
- [x] Middleware supports
- [x] No third-party depdencies
- [x] Unit tests and benchmarks

# Examples

See `_examples` folder.

# License

Apache 2.0 License
