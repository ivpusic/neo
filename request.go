package neo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Purpose of this struct is to be used by authentication and authorization middlewares.
// Saving session is common in many web applications, and this struct is trying to provide
// type-safe version of it.
type Session struct {
	Authenticated bool
	Id            int
	Ids           string
	Role          string
	Rights        []string
	Data          map[string]interface{}
}

// Wrapped http.Request. It contains utility methods for dealing with content of incomming http.Request instance.
type Request struct {
	*http.Request
	Params UrlParam
	// general purpose session instance
	Session Session
}

// Make cookie map.
func buildReqCookies(cookies []*http.Cookie) Cookie {
	result := Cookie{}

	for _, cookie := range cookies {
		result[cookie.Name] = cookie
	}

	return result
}

// Make new Request instance
func makeRequest(req *http.Request) *Request {
	request := &Request{
		Request: req,
	}

	return request
}

// Parse incomming Body as JSON
func (r *Request) JsonBody(instance interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	return decoder.Decode(instance)
}

func (r *Request) StringBody() (string, error) {
	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
