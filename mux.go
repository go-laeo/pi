package pi

import (
	"context"
	"net/http"
	"net/url"
)

const (
	customANY = "ANY"
)

type routePathParam struct{}

var defaultOnNotFound http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(404)
})

type ServerMux interface {
	http.Handler

	Get(path string, h http.Handler)
	Post(path string, h http.Handler)
	Put(path string, h http.Handler)
	Delete(path string, h http.Handler)
	Patch(path string, h http.Handler)
	Options(path string, h http.Handler)
	Head(path string, h http.Handler)

	// Any insert a handler for path without check http method.
	Any(path string, h http.Handler)

	// Group insert routes with same prefix.
	Group(prefix string, fn func(sm ServerMux))

	// OnNotFound sets a handler for undefined routes.
	OnNotFound(h http.Handler)
}

type servermux struct {
	ctx        context.Context
	onnotfound http.Handler
	root       *Route
	prefix     string
}

func NewServerMux(ctx context.Context) ServerMux {
	return &servermux{
		ctx:        ctx,
		root:       &Route{},
		onnotfound: defaultOnNotFound,
	}
}

var _ ServerMux = (*servermux)(nil)

func (sm *servermux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cap := make(url.Values)
	r = r.WithContext(context.WithValue(sm.ctx, &routePathParam{}, cap))
	n := sm.root.Search(r.URL.Path, cap)
	if n != nil {
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
	}

	sm.onnotfound.ServeHTTP(w, r)
}

func (sm *servermux) Get(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodGet, h)
}

func (sm *servermux) Post(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodPost, h)
}

func (sm *servermux) Put(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodPut, h)
}

func (sm *servermux) Delete(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodDelete, h)
}

func (sm *servermux) Patch(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodPatch, h)
}

func (sm *servermux) Options(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodOptions, h)
}

func (sm *servermux) Head(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, http.MethodHead, h)
}

func (sm *servermux) Any(path string, h http.Handler) {
	sm.root.Insert(sm.prefix+path, customANY, h)
}

func (sm *servermux) Group(prefix string, fn func(sm ServerMux)) {
	prev := sm.prefix
	sm.prefix = sm.prefix + prefix
	fn(sm)
	sm.prefix = prev
}

func (sm *servermux) OnNotFound(h http.Handler) {
	sm.onnotfound = h
}
