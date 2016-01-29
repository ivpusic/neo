package neo

import (
	"container/list"
	"errors"

	"github.com/ivpusic/urlregex"
)

///////////////////////////////////////////////////////////////////
// Route
///////////////////////////////////////////////////////////////////

type handler func(*Ctx) (int, error)

// composing route with middlewares. Will be called from ``compose`` fn.
func (h handler) apply(ctx *Ctx, fns []appliable, current int) {
	status, err := h(ctx)

	if err != nil {
		log.Errorf("Error returned from route handler. %s", err.Error())
		ctx.Res.Status = 500
	} else {
		ctx.Res.Status = status
	}

	current++
	if len(fns) > current {
		fns[current].apply(ctx, fns, current)
	}
}

type Route struct {
	*interceptor
	fn      handler
	regex   urlregex.UrlRegex
	fnChain func(*Ctx)
}

///////////////////////////////////////////////////////////////////
// Router
///////////////////////////////////////////////////////////////////

type router struct {
	*interceptor
	*methods
	regions *list.List
}

type Region struct {
	*interceptor
	*methods
}

func (r *router) initRouter() {
	m := &methods{}
	i := &interceptor{[]appliable{}}

	r.methods = m.init()
	r.interceptor = i
	r.regions = list.New()
}

// making new region instance
func (r *router) makeRegion() *Region {
	m := &methods{}
	i := &interceptor{[]appliable{}}

	region := &Region{i, m.init()}
	r.regions.PushBack(region)

	return region
}

// forcing all middlewares to copy to appropriate routes
// also cleans up intrnal structures
func (router *router) flush() {
	// first copy all middlewares to top-level routes
	for method, routesSlice := range router.routes {
		for i, route := range routesSlice {

			route.fnChain = compose(merge(
				router.middlewares,
				route.middlewares,
				[]appliable{route.fn},
			))

			router.routes[method][i] = route
		}
	}

	// then copy all middlewares for region level routes
	// this will copy all b/a from router and region
	// so at the end all b/a for route will be placed directly in route
	//
	// mr = maybeRegion
	for mr := router.regions.Front(); mr != nil; mr = mr.Next() {
		region, ok := mr.Value.(*Region)
		if !ok {
			log.Error("cannot convert element to region")
			continue
		}

		for method, routesSlice := range region.routes {
			for _, route := range routesSlice {

				route.fnChain = compose(merge(
					router.middlewares,
					region.middlewares,
					route.middlewares,
					[]appliable{route.fn},
				))

				router.routes[method] = append(router.routes[method], route)
			}

			// remove method key from region (GET, POST,...)
			delete(region.routes, method)
		}
	}

	// remove regions, we don't need them anymore
	router.regions.Init()
}

// main API function of router
// if will try to invoke fns which are realted to route (if any)
func (router *router) match(req *Request) (*Route, error) {
	// todo: consider using goroutines for searching for route

	method := req.Method
	path := req.URL.Path

	for _, route := range router.routes[method] {
		params, err := route.regex.Match(path)

		if err == nil {
			req.Params = params
			return &route, nil
		}
	}

	return nil, errors.New("[" + req.Method + "] handler for route " + path + " not found")
}
