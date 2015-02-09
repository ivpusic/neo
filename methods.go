package neo

import (
	"github.com/ivpusic/urlregex"
)

// supported HTTP methods
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
)

// slice of all supported methods. Used for testing.
var methodsSlice = []string{GET, POST, PUT, DELETE, OPTIONS, HEAD}

// map of all known routes.
// Map has two levels. First level is HTTP method, and second level is regex and it's actual route instance.
// So once row in this map could be:
// route := routemap["GET"]["^\/$"], where "^\/$" is actual UrlRegex instance.
type routemap map[string]map[urlregex.UrlRegex]Route

type methods struct {
	routes routemap
}

func (m *methods) init() *methods {
	m.routes = routemap{
		GET:     make(map[urlregex.UrlRegex]Route),
		POST:    make(map[urlregex.UrlRegex]Route),
		PUT:     make(map[urlregex.UrlRegex]Route),
		DELETE:  make(map[urlregex.UrlRegex]Route),
		OPTIONS: make(map[urlregex.UrlRegex]Route),
		HEAD:    make(map[urlregex.UrlRegex]Route),
	}

	return m
}

func (m *methods) add(path string, fn handler, method string) *Route {
	route := Route{&interceptor{[]appliable{}}, fn, nil}
	m.routes[method][urlregex.Pattern(path)] = route

	return &route
}

// Registering route handler for ``GET`` request on provided path
func (m *methods) Get(path string, fn handler) *Route {
	return m.add(path, fn, GET)
}

// Registering route handler for ``POST`` request on provided path
func (m *methods) Post(path string, fn handler) *Route {
	return m.add(path, fn, POST)
}

// Registering route handler for ``PUT`` request on provided path
func (m *methods) Put(path string, fn handler) *Route {
	return m.add(path, fn, PUT)
}

// Registering route handler for ``DELETE`` request on provided path
func (m *methods) Delete(path string, fn handler) *Route {
	return m.add(path, fn, DELETE)
}

// Registering route handler for ``OPTIONS`` request on provided path
func (m *methods) Options(path string, fn handler) *Route {
	return m.add(path, fn, OPTIONS)
}

// Registering route handler for ``HEAD`` request on provided path
func (m *methods) Head(path string, fn handler) *Route {
	return m.add(path, fn, HEAD)
}
