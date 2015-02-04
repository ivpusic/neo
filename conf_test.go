package neo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	assert.True(t, true)
	conf := new(Conf)
	conf.Parse("./testassets/testconf.toml")

	assert.Equal(t, "go run main.go", conf.Hotreload.Command)
	assert.Equal(t, []string{".", ".."}, conf.Hotreload.Watch)
	assert.Equal(t, []string{"/some/path"}, conf.Hotreload.Ignore)

	assert.Equal(t, ":3000", conf.App.Addr)
	assert.Equal(t, "DEBUG", conf.App.Logger.Level)
	assert.Equal(t, "application-test", conf.App.Logger.Name)

	assert.Equal(t, "INFO", conf.Neo.Logger.Level)

	assert.Panics(t, func() {
		conf.Parse("./testassets/testconf_inv.toml")
	})
}
