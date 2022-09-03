package ezy

type Connector[T any] func(next HandlerFunc[T]) HandlerFunc[T]

func (h HandlerFunc[T]) Connect(cc ...Connector[T]) HandlerFunc[T] {
	for i, j := 0, len(cc); i < j; i++ {
		h = cc[j-i-1](h)
	}
	return h
}
