package neo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBefore(t *testing.T) {
	assert.True(t, true)
}

func testMiddleware(ctx *Ctx, next Next) {
	counter := ctx.Data.Get("counter").(int)
	counter = counter + 1
	ctx.Data.Set("counter", counter)

	next()
}

func testMiddlewareWithDownstream(ctx *Ctx, next Next) {
	counter := ctx.Data.Get("counter-downstream").(int)
	counter = counter + 1
	ctx.Data.Set("counter-downstream", counter)

	next()

	counter = counter + 1
	ctx.Data.Set("counter-downstream", counter)
}

func testMiddlewareWithoutDownstream(ctx *Ctx, next Next) {
	counter := ctx.Data.Get("counter").(int)
	counter = counter + 1
	ctx.Data.Set("counter", counter)
}

func TestCompose(t *testing.T) {
	var mdw Middleware = testMiddleware

	middlewares := []appliable{
		mdw,
		mdw,
		mdw,
		mdw,
	}

	fn := compose(middlewares)
	assert.NotNil(t, fn)

	ctx := &Ctx{
		Data: CtxData{},
	}
	ctx.Data.Set("counter", 0)

	fn(ctx)

	counter := ctx.Data.Get("counter").(int)
	assert.Exactly(t, 4, counter)
}

func TestComposeWithDownstream(t *testing.T) {
	var downstream Middleware = testMiddlewareWithDownstream
	var normalMdw Middleware = testMiddleware

	middlewares := []appliable{
		normalMdw,
		downstream,
		normalMdw,
		normalMdw,
	}

	fn := compose(middlewares)
	assert.NotNil(t, fn)

	ctx := &Ctx{
		Data: CtxData{},
	}
	ctx.Data.Set("counter", 0)
	ctx.Data.Set("counter-downstream", 0)

	fn(ctx)

	counter := ctx.Data.Get("counter-downstream").(int)
	assert.Exactly(t, 2, counter)
}

func TestComposeWithoutDownstream(t *testing.T) {
	var noDownstream Middleware = testMiddlewareWithoutDownstream

	middlewares := []appliable{
		noDownstream,
		noDownstream,
		noDownstream,
		noDownstream,
	}

	fn := compose(middlewares)
	assert.NotNil(t, fn)

	ctx := &Ctx{
		Data: CtxData{},
	}
	ctx.Data.Set("counter", 0)

	fn(ctx)

	assert.Exactly(t, 1, ctx.Data.Get("counter").(int))
}

func TestInterceptorUse(t *testing.T) {
	i := &interceptor{[]appliable{}}
	i.Use(testMiddleware)
	i.Use(testMiddleware)
	i.Use(testMiddleware)

	assert.Exactly(t, 3, len(i.middlewares))
}
