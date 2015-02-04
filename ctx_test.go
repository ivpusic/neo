package neo

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestMakeCtx(t *testing.T) {
	req := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	ctx := makeCtx(req, w)
	assert.NotNil(t, ctx)
}

func TestCtxReqRes(t *testing.T) {
	req := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	ctx := makeCtx(req, w)

	assert.NotNil(t, ctx.Req)
	assert.NotNil(t, ctx.Req.URL)

	assert.NotNil(t, ctx.Res)
	assert.Equal(t, 404, ctx.Res.Status)
}

func TestCtxData(t *testing.T) {
	req := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	ctx := makeCtx(req, w)

	assert.Nil(t, ctx.Data.Get("data"))
	ctx.Data.Set("data", "value")

	assert.Equal(t, "value", ctx.Data.Get("data"))

	ctx.Data.Del("data")
	assert.Nil(t, ctx.Data.Get("data"))
}

func TestCtxOnEmit(t *testing.T) {
	req := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	ctx := makeCtx(req, w)

	channel := make(chan int)
	ok := false
	ctx.On("event", func(data interface{}) {
		counter := data.(int)
		if counter == 10 {
			ok = true
		}

		channel <- 1
	})

	ctx.Emit("event", 5)
	// wait
	<-channel

	assert.False(t, ok)
	ok = false

	ctx.Emit("event", 10)
	// wait
	<-channel

	assert.True(t, ok)
}
