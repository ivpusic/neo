package neo

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func makeReq(path, method string) *Request {
	httpRequest := &http.Request{
		Method: method,
		URL: &url.URL{
			Path: path,
		},
	}

	return &Request{
		Request: httpRequest,
	}
}

func testHandlerFn(this *Ctx) {
	log.Debug("route handler called")
}

func testInterceptorFn(this *Ctx, next Next) {
	log.Debug("middleware interceptor handler called")
	next()
}

func TestRouterUse(t *testing.T) {
	router := &router{}
	router.initRouter()

	router.Use(testInterceptorFn)
	router.Use(testInterceptorFn)
	router.Use(testInterceptorFn)
	router.Use(testInterceptorFn)
	router.Use(testInterceptorFn)

	assert.Equal(t, 5, len(router.middlewares), "All middlewares should be added")
	router.flush()
	assert.Equal(t, 5, len(router.middlewares), "Middlewares should be cleaned up after flush")
}

func TestRouterMatch(t *testing.T) {
	testPath := "/some/path"

	for _, method := range methodsSlice {
		router := &router{}
		counter := 0
		router.initRouter()

		fn := func(this *Ctx, next Next) {
			log.Debug("middleware interceptor handler called")
			counter++
			next()
			counter++
		}

		handler := func(this *Ctx) (int, error) {
			counter++
			return 200, nil
		}

		router.Use(fn)
		router.Use(fn)

		// make route
		router.add(testPath, handler, method).Use(fn)
		router.flush()

		route, err := router.match(makeReq(testPath, method))
		assert.Nil(t, err)
		assert.NotNil(t, route)

		req := makeTestHttpRequest(nil)
		w := httptest.NewRecorder()
		ctx := makeCtx(req, w)

		route.fnChain(ctx)
		assert.Equal(t, 7, counter)

		router.initRouter()
		for _, tmpMethod := range methodsSlice {
			if tmpMethod != method {
				route, err := router.match(makeReq(testPath, tmpMethod))
				assert.NotNil(t, err)
				assert.Nil(t, route)
			}
		}
	}
}

func TestMakeRegion(t *testing.T) {
	router := &router{}
	router.initRouter()

	region := router.makeRegion()
	assert.NotNil(t, region)
	assert.NotNil(t, region.interceptor)
	assert.NotNil(t, region.methods)
	assert.Exactly(t, 1, router.regions.Len())

	region = router.makeRegion()
	assert.NotNil(t, region)
	assert.NotNil(t, region.interceptor)
	assert.NotNil(t, region.methods)
	assert.Exactly(t, 2, router.regions.Len())
}

func TestRegionMatch(t *testing.T) {
	router := &router{}
	router.initRouter()
	region := router.makeRegion()
	counter := 0
	testPath := "/some"

	fn := func(this *Ctx, next Next) {
		log.Debug("middleware interceptor handler called")
		counter++
		next()
		counter++
	}

	handler := func(this *Ctx) (int, error) {
		counter++
		return 200, nil
	}

	region.Use(fn)
	region.Use(fn)
	_route := region.Get("/some", handler)
	_route.Use(fn)
	_route.Use(fn)

	router.flush()

	route, err := router.match(makeReq(testPath, GET))
	assert.Nil(t, err)

	req := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	ctx := makeCtx(req, w)

	route.fnChain(ctx)
	assert.Equal(t, 9, counter)

	route, err = router.match(makeReq("/some/unknown/path", GET))
	assert.NotNil(t, err)
	assert.Nil(t, route)
}
