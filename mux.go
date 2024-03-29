package pi

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
)

var defaultNotFoundHandler HandlerFunc = func(ctx Context) error {
	return ctx.Code(404)
}

var defaultErrorFormatter = func(ctx Context, err error) {
	ctx.Error(http.StatusInternalServerError, &ErrorResult{
		Error:        "unknown",
		ErrorMessage: err.Error(),
	})
}

type ServerMux interface {
	http.Handler

	Route(path string) Route
	Group(prefix string, fn func(sm ServerMux))
	SetNotFoundHandler(h HandlerFunc)
	SetErrorFormatter(fn func(ctx Context, err error))
	Use(c func(next HandlerFunc) HandlerFunc)
}

var _ ServerMux = (*servermux)(nil)

type servermux struct {
	notFoundHandler HandlerFunc
	root            *_route
	capcap          *sync.Pool
	errorFormater   func(ctx Context, err error)
	prefix          string
	cc              []func(next HandlerFunc) HandlerFunc
}

func NewServerMux() ServerMux {
	return &servermux{
		root:            createRootRoute(),
		notFoundHandler: defaultNotFoundHandler,
		errorFormater:   defaultErrorFormatter,
		capcap: &sync.Pool{
			New: func() any {
				return make(url.Values)
			},
		},
	}
}

func (sm *servermux) SetErrorFormatter(fn func(ctx Context, err error)) {
	sm.errorFormater = fn
}

func (sm *servermux) SetNotFoundHandler(h HandlerFunc) {
	sm.notFoundHandler = h
}

func (sm *servermux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cap := sm.capcap.Get().(url.Values)
	defer func() {
		for k := range cap {
			cap[k] = cap[k][:0] // reset slice to empty, but the keys in map will keep
		}
		sm.capcap.Put(cap)
	}()

	ctx := createContext(w, r, cap)

	var err error
	n := sm.root.Search(r.URL.Path, cap) // 1 allocs/op
	if n == nil {
		err = ErrHandlerNotFound
	} else {
		// 2 allocs/op
		err = n.Invoke(ctx)
	}

	if errors.Is(err, ErrHandlerNotFound) {
		err = sm.notFoundHandler(ctx)
	}

	if err != nil {
		sm.errorFormater(ctx, err)
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

func (sm *servermux) Use(c func(next HandlerFunc) HandlerFunc) {
	sm.cc = append(sm.cc, c)
}
