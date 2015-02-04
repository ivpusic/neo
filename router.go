package neo

import (
	"container/list"
	"errors"
)

///////////////////////////////////////////////////////////////////
// Route
///////////////////////////////////////////////////////////////////

type handler func(*Ctx)

// composing route with middlewares. Will be called from ``compose`` fn.
func (h handler) apply(ctx *Ctx, fns []appliable) {
	h(ctx)

	if len(fns) > 0 {
		fns[0].apply(ctx, fns[1:])
	}
}

type Route struct {
	*interceptor
	fn      handler
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
	for method, routesmap := range router.routes {
		for path, route := range routesmap {

			route.fnChain = compose(merge(
				router.middlewares,
				route.middlewares,
				[]appliable{route.fn},
			))

			router.routes[method][path] = route
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

		for method, routesmap := range region.routes {
			for path, route := range routesmap {

				route.fnChain = compose(merge(
					router.middlewares,
					region.middlewares,
					route.middlewares,
					[]appliable{route.fn},
				))

				router.routes[method][path] = route

				// remove route from region
				// we are doing this because route reference is now in top level
				delete(region.routes[method], path)
			}

			// remove method key from region (GET, POST,...)
			delete(region.routes, method)
		}

		// remove region from router region list
		// in case of failed cast to region, this won't be invoked
		router.regions.Remove(mr)
	}
}

// main API function of router
// if will try to invoke fns which are realted to route (if any)
func (router *router) match(req *Request) (*Route, error) {
	// todo: consider using goroutines for searching for route

	method := req.Method
	path := req.URL.Path

	for k, _ := range router.routes[method] {
		params, err := k.Match(path)

		if err == nil {
			route := router.routes[method][k]
			req.Params = params
			return &route, nil
		}
	}

	return nil, errors.New("[" + req.Method + "] handler for route " + path + " not found")
}
