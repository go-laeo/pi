package pi

import (
	"net/http"
	"net/url"
	"sync"
)

var defaultOnNotFound HandlerFunc = func(ctx Context) error {
	return ctx.Code(404)
}

type ServerMux interface {
	http.Handler

	Route(path string) Route

	// Group insert routes with same prefix.
	Group(prefix string, fn func(sm ServerMux))

	// OnNotFound sets a handler for undefined routes.
	OnNotFound(h HandlerFunc)

	Use(c func(next HandlerFunc) HandlerFunc)
}

var _ ServerMux = (*servermux)(nil)

type servermux struct {
	onnotfound HandlerFunc
	root       *_route
	capcap     *sync.Pool
	prefix     string
	cc         []func(next HandlerFunc) HandlerFunc
}

func NewServerMux() ServerMux {
	return &servermux{
		root:       createRootRoute(),
		onnotfound: defaultOnNotFound,
		capcap: &sync.Pool{
			New: func() any {
				return make(url.Values)
			},
		},
	}
}

func (sm *servermux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cap := sm.capcap.Get().(url.Values)
	defer func() {
		for k := range cap {
			cap[k] = cap[k][:0] // reset slice to empty, but the keys in map will keep
		}
		sm.capcap.Put(cap)
	}()

	ctx := createContext(w, r, cap, sm)

	n := sm.root.Search(r.URL.Path, cap) // 1 allocs/op
	if n == nil {
		sm.onnotfound(ctx)
		return
	}

	// 2 allocs/op
	ok := n.Invoke(ctx)
	if !ok {
		sm.onnotfound(ctx)
	}
}

func (sm *servermux) Route(path string) Route {
	return sm.root.Insert(sm.prefix+path, sm.cc...)
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
