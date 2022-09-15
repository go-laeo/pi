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

var _routePathParam *routePathParam

var defaultOnNotFound http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(404)
})

type ServerMux interface {
	http.Handler

	Get(path string, h HandlerFunc)
	Post(path string, h HandlerFunc)
	Put(path string, h HandlerFunc)
	Delete(path string, h HandlerFunc)
	Patch(path string, h HandlerFunc)
	Options(path string, h HandlerFunc)
	Head(path string, h HandlerFunc)

	// Any insert a handler for path without check http method.
	Any(path string, h HandlerFunc)

	// Group insert routes with same prefix.
	Group(prefix string, fn func(sm ServerMux))

	// OnNotFound sets a handler for undefined routes.
	OnNotFound(h HandlerFunc)

	Use(c func(next HandlerFunc) HandlerFunc)
}

type servermux struct {
	ctx        context.Context
	onnotfound http.Handler
	root       *Route
	prefix     string
	cc         []func(next HandlerFunc) HandlerFunc
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
	r = r.WithContext(context.WithValue(sm.ctx, _routePathParam, cap))
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

func (sm *servermux) Get(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodGet, h.Connect(sm.cc...))
}

func (sm *servermux) Post(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodPost, h.Connect(sm.cc...))
}

func (sm *servermux) Put(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodPut, h.Connect(sm.cc...))
}

func (sm *servermux) Delete(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodDelete, h.Connect(sm.cc...))
}

func (sm *servermux) Patch(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodPatch, h.Connect(sm.cc...))
}

func (sm *servermux) Options(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodOptions, h.Connect(sm.cc...))
}

func (sm *servermux) Head(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, http.MethodHead, h.Connect(sm.cc...))
}

func (sm *servermux) Any(path string, h HandlerFunc) {
	sm.root.Insert(sm.prefix+path, customANY, h.Connect(sm.cc...))
}

func (sm *servermux) Group(prefix string, fn func(sm ServerMux)) {
	prev := sm.prefix
	prevCC := sm.cc
	sm.prefix = sm.prefix + prefix
	fn(sm)
	sm.prefix = prev
	sm.cc = prevCC
}

func (sm *servermux) OnNotFound(h HandlerFunc) {
	sm.onnotfound = h
}

func (sm *servermux) Use(c func(next HandlerFunc) HandlerFunc) {
	sm.cc = append(sm.cc, c)
}
