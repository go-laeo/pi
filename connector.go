package pi

func (h HandlerFunc) Connect(cc ...func(next HandlerFunc) HandlerFunc) HandlerFunc {
	for i, j := 0, len(cc); i < j; i++ {
		h = cc[j-i-1](h)
	}
	return h
}
