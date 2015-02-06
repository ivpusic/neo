package neo

import (
	"github.com/ivpusic/neo/ebus"
	"net/http"
)

///////////////////////////////////////////////////////////////////
// Context data representation and operations
///////////////////////////////////////////////////////////////////

// Map which will hold your Request contextual data.
type CtxData map[string]interface{}

// Get context data on key
func (r CtxData) Get(key string) interface{} {
	return r[key]
}

// Add new contextual data on key with provided value
func (r CtxData) Set(key string, value interface{}) {
	r[key] = value
}

// Delete contextual data on key
func (r CtxData) Del(key string) {
	delete(r, key)
}

///////////////////////////////////////////////////////////////////
// Context
///////////////////////////////////////////////////////////////////

// Representing context for this request.
type Ctx struct {
	ebus.EBus

	// Wrapped http.Request object.
	Req *Request

	// Object with utility methods for writing responses to http.ResponseWriter
	Res *Response

	// Context data map
	Data CtxData
}

// Will make default contextual data based on provided request and ResponseWriter
func makeCtx(req *http.Request, w http.ResponseWriter) *Ctx {
	request := makeRequest(req)
	response := makeResponse(req, w)
	// 404 - OK by default
	// this can be overided by middlewares and handler fns
	response.Status = 404

	ctx := &Ctx{
		Req:  request,
		Res:  response,
		Data: CtxData{},
	}
	ctx.InitEBus()

	return ctx
}
