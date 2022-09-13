package pi

type Connector func(next HandlerFunc) HandlerFunc

func (h HandlerFunc) Connect(cc ...Connector) HandlerFunc {
	for i, j := 0, len(cc); i < j; i++ {
		h = cc[j-i-1](h)
	}
	return h
}
