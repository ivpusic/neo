package neo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegion(t *testing.T) {
	app := App()

	region1 := app.Region()
	assert.NotNil(t, region1)

	region2 := app.Region()
	assert.NotNil(t, region2)
}
