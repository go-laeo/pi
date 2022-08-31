package ezy

import (
	"net/http"
	"path"
	"strings"
)

type Route struct {
	lit  string
	name string
	sub  map[string]*Route
	hmap map[string]http.Handler
}

func (p *Route) Search(route string) *Route {
	route = path.Clean(route)
	root := p
	for _, seg := range strings.Split(route, "/") {
		n, ok := root.sub[seg]
		if ok {
			root = n
			continue
		}

		return nil
	}

	return root
}

func (p *Route) Insert(route string, method string, h http.Handler) *Route {
	route = path.Clean(route)
	root := p
	for _, seg := range strings.Split(route, "/") {
		if root.sub == nil {
			root.sub = make(map[string]*Route)
		}

		n, ok := root.sub[seg]
		if ok {
			root = n
			continue
		}

		n = &Route{lit: seg}
		root.sub[seg] = n
		root = n
	}

	if root.hmap == nil {
		root.hmap = make(map[string]http.Handler)
	}

	root.hmap[strings.ToUpper(method)] = h

	return root
}

func (p *Route) Alias(name string) *Route {
	p.name = name
	return p
}
