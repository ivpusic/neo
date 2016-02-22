package neo

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	static := Static{}
	static.init()

	assert.NotNil(t, static.serving)

	static.Serve("/some/url", "./testassets")

	abs, err := filepath.Abs("./testassets")
	assert.Nil(t, err)

	match, err := static.match("/some/url/test.txt")
	assert.Nil(t, err)
	assert.Exactly(t, abs+"/test.txt", match)

	match, err = static.match("/some/unknown/url")
	assert.NotNil(t, err)
	assert.Equal(t, "", match)
}

func TestCanBeServed(t *testing.T) {
	static := Static{}
	static.init()

	assert.NotNil(t, static.serving)
	static.Serve("/some/url", "./testassets")

	abs, err := filepath.Abs("./testassets")
	assert.Nil(t, err)

	assert.False(t, static.canBeServed(abs))
	assert.True(t, static.canBeServed(abs+"/index.html"))
	assert.False(t, static.canBeServed(abs+"/unknown.html"))
}

func TestServeIndexPages(t *testing.T) {
	static := Static{}
	static.init()

	static.Serve("/", "./testassets")
	abs, err := filepath.Abs("./testassets")

	// it should server index.html inside main directory
	match, err := static.match("/")
	assert.Nil(t, err)
	assert.Exactly(t, abs+"/index.html", match)

	// it should serve index.html inside child directories
	match, err = static.match("/dir")
	assert.Nil(t, err)
	assert.Exactly(t, abs+"/dir/index.html", match)

	// it should not server index.tpl
	match, err = static.match("/templates")
	assert.NotNil(t, err)
}
