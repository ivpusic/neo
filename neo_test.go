package neo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	app := App()
	assert.NotNil(t, app)
}
