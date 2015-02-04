package aux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func makeTestCookie(key, value string) *http.Cookie {
	return &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}

func makeTestResCookies() *ResCookie {
	return &ResCookie{
		"key1": makeTestCookie("key1", "value1"),
		"key2": makeTestCookie("key2", "value2"),
	}
}

func TestReqCookieGet(t *testing.T) {
	cookies := ReqCookie{
		"key": makeTestCookie("key", "value"),
	}

	assert.NotNil(t, cookies.Get("key"))
	assert.Nil(t, cookies.Get("unknown"))
}

func TestResCookieGet(t *testing.T) {
	cookies := makeTestResCookies()
	assert.NotNil(t, cookies.Get("key1"))
	assert.NotNil(t, cookies.Get("key2"))
	assert.Nil(t, cookies.Get("key3"))
}

func TestResCookieSet(t *testing.T) {
	cookies := makeTestResCookies()
	cookies.Set("key3", "value3")
	assert.NotNil(t, cookies.Get("key3"))
	assert.Nil(t, cookies.Get("key4"))
}

func TestResCookieDel(t *testing.T) {
	cookies := makeTestResCookies()

	assert.NotNil(t, cookies.Get("key1"))
	cookies.Del("key1")
	assert.Nil(t, cookies.Get("key1"))
	assert.NotNil(t, cookies.Get("key2"))
}

func TestResCookieSetCustom(t *testing.T) {
	cookies := makeTestResCookies()

	assert.Nil(t, cookies.Get("key3"))
	cookies.SetCustom(makeTestCookie("key3", "value3"))
	assert.NotNil(t, cookies.Get("key3"))
}
