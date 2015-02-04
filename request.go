package neo

import (
	"encoding/json"
	"github.com/ivpusic/neo/aux"
	"net/http"
)

// Wrapped http.Request. It contains utility methods for dealing with content of incomming http.Request instance.
type Request struct {
	*http.Request
	Params aux.UrlParam
}

// Make cookie map.
func buildReqCookies(cookies []*http.Cookie) aux.Cookie {
	result := aux.Cookie{}

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
