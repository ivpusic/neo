package neo

type Next func()

type appliable interface {
	apply(ctx *Ctx, fns []appliable)
}

// Composing multiple middlewares into one chained function.
func compose(fns []appliable) func(*Ctx) {
	return func(ctx *Ctx) {
		fns[0].apply(ctx, fns[1:])
	}
}

// Representing middleware handler.
// It accepts Neo Context and Next function for calling next middleware in chain.
type middleware func(*Ctx, Next)

func (m middleware) apply(ctx *Ctx, fns []appliable) {
	m(ctx, func() {
		if len(fns) > 0 {
			fns[0].apply(ctx, fns[1:])
		}
	})
}

// contains list of functions for intercepting requests before and after route handle call
type interceptor struct {
	middlewares []appliable
}

// Adding new middleware into chain of middlewares.
func (m *interceptor) Use(fn middleware) {
	m.middlewares = append(m.middlewares, appliable(&fn))
}
