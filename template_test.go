package neo

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type testWriter struct {
	data []byte
}

func (t *testWriter) Write(data []byte) (int, error) {
	t.data = append(t.data, data...)
	return 0, nil
}

func TestCompileTpl(t *testing.T) {
	writer := &testWriter{}
	regex := regexp.MustCompile("(?s)<html>.*\n</html>")

	// before compiling template, tpl is nil
	tpl = nil
	err := renderTpl(writer, "testindex", nil)
	assert.NotNil(t, err)

	writer = &testWriter{}
	compileTpl("./testassets/templates/*")
	assert.Nil(t, writer.data)

	// known template
	err = renderTpl(writer, "testindex", nil)
	assert.Nil(t, err)
	assert.NotNil(t, writer.data)

	match := regex.Match(writer.data[:])
	assert.True(t, match)

	// unknown template
	writer.data = nil
	err = renderTpl(writer, "unknown", nil)
	assert.NotNil(t, err)

	match = regex.Match(writer.data[:])
	assert.False(t, match)

	// wrong path
	assert.Panics(t, func() {
		compileTpl("./testassets/unknown/**")
	})
}

func TestTplExists(t *testing.T) {
	assert.False(t, tplExists("testindex"))
	compileTpl("./testassets/templates/*")
	assert.True(t, tplExists("testindex"))
}
