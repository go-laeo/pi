package ezy

type Composer[T any] []func(next HandlerFunc[T]) HandlerFunc[T]

func Compose[T any](middlewares ...func(next HandlerFunc[T]) HandlerFunc[T]) Composer[T] {
	return Composer[T](middlewares)
}

func (c Composer[T]) With(middlewares ...func(next HandlerFunc[T]) HandlerFunc[T]) Composer[T] {
	return Composer[T](append(c, middlewares...))
}

func (c Composer[T]) For(h HandlerFunc[T]) HandlerFunc[T] {
	for i, j := 0, len(c); i < j; i++ {
		h = c[j-i-1](h)
	}
	return h
}
