package neo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	app := App()
	assert.NotNil(t, app)
}

func TestAssert(t *testing.T) {
	assert.Panics(t, func() {
		Assert(false, 200, []byte("some message"))
	})

	assert.NotPanics(t, func() {
		Assert(true, 200, []byte("some message"))
	})
}

func TestAssertNil(t *testing.T) {
	assert.Panics(t, func() {
		AssertNil(&testPerson{"Name", 20}, 200, []byte("some message"))
	})

	assert.NotPanics(t, func() {
		AssertNil(nil, 200, []byte("some message"))
	})
}

func TestAssertNotNil(t *testing.T) {
	assert.NotPanics(t, func() {
		AssertNotNil(&testPerson{"Name", 20}, 200, []byte("some message"))
	})

	assert.Panics(t, func() {
		AssertNotNil(nil, 200, []byte("some message"))
	})
}
