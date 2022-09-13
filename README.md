# pi

![build.yaml](https://github.com/go-laeo/pi/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/pi.svg)](https://pkg.go.dev/github.com/go-laeo/pi)

A `pi` is your good helper to build a clean JSON API server using Golang Generics.

# Code Sample

```Go
type UserData struct {
    Name string
    Password string
}

var h pi.HandlerFunc = func (ctx pi.Context) error {
    data := &UserData{}
    err := pi.Format(ctx, data)
    if err != nil {
        return pi.NewError(400, err.Error())
    }

    // do sth. actions...

    return ctx.Text("hello, world!")
}

sm := pi.NewServerMux(context.Background())
sm.Post("/api/v1/users", h)

http.ListenAndServe("localhost:8080", sm)
```

# Install

```shell
go get -u github.com/go-laeo/pi
```

# Features

- [x] Fast routing, routes group, route params and wildcard route
- [x] `net/http` compatible (`pi.HandlerFunc` is a `http.Handler`)
- [x] ~~Auto~~ Manually decode HTTP body using function `pi.Format[T any]`
- [x] Middleware supports by `pi.Connector`
- [x] Built-in `pi.FileServer` for SPA
- [x] No third-party depdencies
- [x] Unit tests and benchmarks

# Examples

See `_examples` folder.

# License

Apache 2.0 License
