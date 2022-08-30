package ezy

import "net/http"

const (
	customANY = "ANY"
)

var defaultNotFoundHandler = Compose(func(ctx Context, p *any) error {
	ctx.WriteHeader(404)
	ctx.Write([]byte("not found"))
	return nil
})

type Router struct {
	root     *Route
	notfound http.Handler
}

func NewRouter() *Router {
	return &Router{
		root:     &Route{},
		notfound: defaultNotFoundHandler,
	}
}

var _ http.Handler = (*Router)(nil)

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := ro.root.Search(r.URL.Path)
	if n == nil {
		ro.notfound.ServeHTTP(w, r)
		return
	}

	fn, ok := n.hmap[r.Method]
	if ok {
		fn.ServeHTTP(w, r)
		return
	}

	fn, ok = n.hmap[customANY]
	if ok {
		fn.ServeHTTP(w, r)
		return
	}

	ro.notfound.ServeHTTP(w, r)
}

func (ro *Router) Get(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodGet, h)
}

func (ro *Router) Post(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodPost, h)
}

func (ro *Router) Put(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodPut, h)
}

func (ro *Router) Delete(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodDelete, h)
}

func (ro *Router) Patch(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodPatch, h)
}

func (ro *Router) Options(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodOptions, h)
}

func (ro *Router) Head(path string, h http.Handler) {
	ro.root.Insert(path, http.MethodHead, h)
}

func (ro *Router) Any(path string, h http.Handler) {
	ro.root.Insert(path, customANY, h)
}
