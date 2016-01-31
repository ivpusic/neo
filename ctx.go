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

	// list of errors happened in current http request
	Errors []error
}

// Append error to list of errors for current http request.
// This function can be used from any part of your code which has access to current context
func (c *Ctx) Error(err error) {
	c.Errors = append(c.Errors, err)
}

// HasErrors checks if we have errors in current http request
func (c *Ctx) HasErrors() bool {
	return len(c.Errors) > 0
}

// Will make default contextual data based on provided request and ResponseWriter
func makeCtx(req *http.Request, w http.ResponseWriter) *Ctx {
	request := makeRequest(req)
	response := makeResponse(req, w)
	// 404 - OK by default
	// this can be overided by middlewares and handler fns
	response.Status = 404

	ctx := &Ctx{
		Req:    request,
		Res:    response,
		Data:   CtxData{},
		Errors: []error{},
	}
	ctx.InitEBus()

	return ctx
}
