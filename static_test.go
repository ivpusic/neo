package neo

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
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

func TestServeIndexPages(t *testing.T) {
	static := Static{}
	static.init()

	static.Serve("/", "./testassets")
	abs, err := filepath.Abs("./testassets")

	match, err := static.match("/")
	assert.Nil(t, err)
	assert.Exactly(t, abs+"/index.html", match)

}
