package ezy

import (
	"net/http"
	"path"
	"strings"
)

const (
	dynamic  = ':'
	wildcard = '*'
)

type Route struct {
	parent *Route

	pattern     string
	placeholder string // pattern without param or wildcard identifier

	hasDynamicChild  bool
	hasWildcardChild bool
	sub              map[string]*Route

	hmap map[string]http.Handler
}

func (p *Route) Search(route string) *Route {
	route = path.Clean(route)
	chunks := strings.Split(route, "/")

	current := p
	for i := 0; i < len(chunks); i++ {
		seg := chunks[i]

		next, ok := current.sub[seg]
		if ok {
			current = next
			continue
		}

		if current.hasDynamicChild {
			next, ok = current.sub[string(dynamic)]
			if ok {
				current = next // continues on dynamic route.
				continue
			}
		}

		if current.hasWildcardChild {
			next, ok = current.sub[string(wildcard)]
			if ok {
				// wildcard route should returns immediately.
				return next
			}
		}

		return nil
	}

	return current
}

func (p *Route) Insert(route string, method string, h http.Handler) *Route {
	route = path.Clean(route)
	current := p
	for _, seg := range strings.Split(route, "/") {
		if current.sub == nil {
			current.sub = make(map[string]*Route)
		}

		next, ok := current.sub[seg]
		if ok {
			current = next
			continue
		}

		next = &Route{
			parent:      current,
			pattern:     seg,
			placeholder: seg,
		}

		if len(seg) > 0 {
			switch next.pattern[0] {
			case dynamic:
				current.hasDynamicChild = true
				next.placeholder = next.placeholder[1:]
				current.sub[string(dynamic)] = next
			case wildcard:
				current.hasWildcardChild = true
				next.placeholder = next.placeholder[1:]
				current.sub[string(wildcard)] = next
			}
		}

		current.sub[seg] = next
		current = next
	}

	if current.hmap == nil {
		current.hmap = make(map[string]http.Handler)
	}

	current.hmap[strings.ToUpper(method)] = h

	return current
}
