package neo

import (
	"net/http"
)

// Map of cookies which will be sent to client.
type Cookie map[string]*http.Cookie

func (cookies Cookie) Get(key string) *http.Cookie {
	return cookies[key]
}

func (cookies Cookie) Set(key string, value string) {
	cookies[key] = &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}

func (cookies Cookie) Del(key string) {
	delete(cookies, key)
}

// Set *http.Cookie instance to Response
func (cookies Cookie) SetCustom(cookie *http.Cookie) {
	cookies[cookie.Name] = cookie
}
