package neo

import (
	"github.com/ivpusic/urlregex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	assert.True(t, true)
}

func TestMethods(t *testing.T) {
	m := new(methods)
	m.init()
	fn := func(this *Ctx) {}

	m.Get("/some", fn)
	assert.NotNil(t, m.routes["GET"][*urlregex.Pattern("/some")])

	m.Post("/some", fn)
	assert.NotNil(t, m.routes["POST"][*urlregex.Pattern("/some")])

	m.Put("/some", fn)
	assert.NotNil(t, m.routes["PUT"][*urlregex.Pattern("/some")])

	m.Delete("/some", fn)
	assert.NotNil(t, m.routes["DELETE"][*urlregex.Pattern("/some")])

	m.Options("/some", fn)
	assert.NotNil(t, m.routes["OPTIONS"][*urlregex.Pattern("/some")])

	m.Head("/some", fn)
	assert.NotNil(t, m.routes["HEAD"][*urlregex.Pattern("/some")])
}
