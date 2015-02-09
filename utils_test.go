package neo

import (
	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAppliable struct {
	id int
}

func (t *testAppliable) apply(ctx *Ctx, fns []appliable, current int) {
}

func TestIsDirectory(t *testing.T) {
	is, err := isDirectory("./testassets/templates")
	assert.Nil(t, err)
	assert.True(t, is)

	is, err = isDirectory("./testassets/unknown")
	assert.NotNil(t, err)
	assert.False(t, is)
}

func TestMerge(t *testing.T) {
	ta1 := testAppliable{1}
	ta2 := testAppliable{2}
	ta3 := testAppliable{3}
	ta4 := testAppliable{4}
	ta5 := testAppliable{5}
	ta6 := testAppliable{6}
	ta7 := testAppliable{7}
	ta8 := testAppliable{8}

	firstSlice := []appliable{&ta1, &ta2}
	secondSlice := []appliable{&ta3, &ta4}
	thirdSlice := []appliable{&ta5, &ta6, &ta7, &ta8}

	result := merge(firstSlice, secondSlice, thirdSlice)
	assert.Exactly(t, 8, len(result))

	for i, value := range result {
		val := value.(*testAppliable)
		assert.Equal(t, i+1, val.id)
	}
}

func TestParseLogLevel(t *testing.T) {
	gologLevel, err := parseLogLevel("debug")
	assert.Nil(t, err)
	assert.NotNil(t, gologLevel)
	assert.Equal(t, golog.DEBUG, gologLevel)

	gologLevel, err = parseLogLevel("info")
	assert.Nil(t, err)
	assert.NotNil(t, gologLevel)
	assert.Equal(t, golog.INFO, gologLevel)

	gologLevel, err = parseLogLevel("warn")
	assert.Nil(t, err)
	assert.NotNil(t, gologLevel)
	assert.Equal(t, golog.WARN, gologLevel)

	gologLevel, err = parseLogLevel("error")
	assert.Nil(t, err)
	assert.NotNil(t, gologLevel)
	assert.Equal(t, golog.ERROR, gologLevel)

	gologLevel, err = parseLogLevel("panic")
	assert.Nil(t, err)
	assert.NotNil(t, gologLevel)
	assert.Equal(t, golog.PANIC, gologLevel)

	gologLevel, err = parseLogLevel("unknown")
	assert.NotNil(t, err)
}
