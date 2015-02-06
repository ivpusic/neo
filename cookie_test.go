package neo

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func makeTestRawCookie(key, value string) *http.Cookie {
	return &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}

func makeTestResCookies() *Cookie {
	return &Cookie{
		"key1": makeTestRawCookie("key1", "value1"),
		"key2": makeTestRawCookie("key2", "value2"),
	}
}

func TestReqCookieGet(t *testing.T) {
	cookies := Cookie{
		"key": makeTestRawCookie("key", "value"),
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
	cookies.SetCustom(makeTestRawCookie("key3", "value3"))
	assert.NotNil(t, cookies.Get("key3"))
}
