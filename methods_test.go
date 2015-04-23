package neo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	assert.True(t, true)
}

func TestMethods(t *testing.T) {
	m := new(methods)
	m.init()
	fn := func(this *Ctx) (int, error) { return 200, nil }

	m.Get("/some", fn)
	assert.NotNil(t, m.routes["GET"][0])

	m.Post("/some", fn)
	assert.NotNil(t, m.routes["POST"][0])

	m.Put("/some", fn)
	assert.NotNil(t, m.routes["PUT"][0])

	m.Delete("/some", fn)
	assert.NotNil(t, m.routes["DELETE"][0])

	m.Options("/some", fn)
	assert.NotNil(t, m.routes["OPTIONS"][0])

	m.Head("/some", fn)
	assert.NotNil(t, m.routes["HEAD"][0])
}
