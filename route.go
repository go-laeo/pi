package ezy

import (
	"path"
	"strings"
)

type Route struct {
	lit string
	h   any
	sub map[string]*Route
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

func (p *Route) Insert(route string, method string, h HandlerFunc[any]) *Route {
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

	root.h = h

	return root
}
