package neo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func makeTestCookie(key, value string) *http.Cookie {
	return &http.Cookie{
		Name:  key,
		Value: value,
	}
}

func makeTestHttpRequest(body io.Reader) *http.Request {
	req, err := http.NewRequest(GET, "/some/path", body)
	req.AddCookie(makeTestCookie("key1", "value1"))

	if err != nil {
		return nil
	}

	return req
}

func TestMakeRequest(t *testing.T) {
	req := makeRequest(makeTestHttpRequest(nil))
	assert.NotNil(t, req)

	cookie, err := req.Cookie("key1")
	assert.NotNil(t, cookie)
	assert.Nil(t, err)

	assert.Equal(t, GET, req.Method)
	assert.Equal(t, "/some/path", req.URL.Path)

	cookie, err = req.Cookie("keyunknown")
	assert.Nil(t, cookie)
	assert.NotNil(t, err)
}

type testStruct struct {
	Name string
	Age  int
}

func TestJsonBody(t *testing.T) {
	// invalid
	httpRequest := makeTestHttpRequest(bytes.NewBuffer([]byte("{'Name': \"SomeName\", \"Age\": 20}")))
	request := makeRequest(httpRequest)

	person := &testStruct{}
	err := request.JsonBody(person)

	assert.NotNil(t, err)

	// valid
	httpRequest = makeTestHttpRequest(bytes.NewBuffer([]byte("{\"Name\": \"SomeName\", \"Age\": 20}")))
	request = makeRequest(httpRequest)

	person = &testStruct{}
	err = request.JsonBody(person)

	assert.Nil(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, "SomeName", person.Name)
	assert.Equal(t, 20, person.Age)
}
