package ezy

import "net/http"

func Cors[T any](next HandlerFunc[T]) HandlerFunc[T] {
	return func(ctx Context, p *T) error {
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
