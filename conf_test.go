package neo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.True(t, true)
	conf := new(Conf)
	conf.Parse("./testassets/testconf.toml")

	assert.Equal(t, []string{"test", "Godeps"}, conf.Hotreload.Ignore)
	assert.Equal(t, []string{".go"}, conf.Hotreload.Suffixes)

	assert.Equal(t, ":3000", conf.App.Addr)
	assert.Equal(t, "DEBUG", conf.App.Logger.Level)
	assert.Equal(t, "application-test", conf.App.Logger.Name)

	assert.Equal(t, "INFO", conf.Neo.Logger.Level)

	assert.Panics(t, func() {
		conf.Parse("./testassets/testconf_inv.toml")
	})
}
