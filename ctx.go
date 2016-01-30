package neo

import (
	"net/http"

	"github.com/ivpusic/neo/ebus"
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

// Purpose of this struct is to be used by authentication and authorization middlewares.
// Saving session is common in many web applications, and this struct is trying to provide
// type-safe version of it.
type Session struct {
	// is user authenticated?
	Authenticated bool
	// data about logged user
	User interface{}
	// other data
	Data interface{}
}

// Representing context for this request.
type Ctx struct {
	ebus.EBus

	// Wrapped http.Request object.
	Req *Request

	// Object with utility methods for writing responses to http.ResponseWriter
	Res *Response

	// Context data map
	Data CtxData

	// general purpose session instance
	Session Session
}

// Will make default contextual data based on provided request and ResponseWriter
func makeCtx(req *http.Request, w http.ResponseWriter) *Ctx {
	request := makeRequest(req)
	response := makeResponse(req, w)
	// StatusNotFound - OK by default
	// this can be overided by middlewares and handler fns
	response.Status = http.StatusNotFound

	ctx := &Ctx{
		Req:  request,
		Res:  response,
		Data: CtxData{},
	}
	ctx.InitEBus()

	return ctx
}
