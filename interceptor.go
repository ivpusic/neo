package neo

type Next func()

type appliable interface {
	apply(ctx *Ctx, fns []appliable, current int)
}

// Composing multiple middlewares into one chained function.
func compose(fns []appliable) func(*Ctx) {
	return func(ctx *Ctx) {
		fns[0].apply(ctx, fns, 0)
	}
}

// Representing middleware handler.
// It accepts Neo Context and Next function for calling next middleware in chain.
type Middleware func(*Ctx, Next)

func (m Middleware) apply(ctx *Ctx, fns []appliable, current int) {
	m(ctx, func() {
		current++
		if len(fns) > current {
			fns[current].apply(ctx, fns, current)
		}
	})
}

// contains list of functions for intercepting requests before and after route handle call
type interceptor struct {
	middlewares []appliable
}

// Adding new middleware into chain of middlewares.
func (m *interceptor) Use(fn Middleware) {
	m.middlewares = append(m.middlewares, appliable(&fn))
}
