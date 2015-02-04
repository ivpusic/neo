package aux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitEbus(t *testing.T) {
	ebus := EBus{}
	ebus.InitEBus()
	assert.NotNil(t, ebus.eventList)
}

func TestOn(t *testing.T) {
	ebus := EBus{}
	ebus.InitEBus()

	ebus.On("test", func(data interface{}) {})

	topic, ok := ebus.eventList["test"]
	assert.True(t, ok)
	assert.Exactly(t, 1, topic.Len())

	_, ok = ebus.eventList["unknown"]
	assert.False(t, ok)
}

func TestEmit(t *testing.T) {
	ebus := EBus{}
	ebus.InitEBus()
	counter := 0

	ebus.On("test", func(data interface{}) {
		counter += data.(int)
	})

	ebus.Emit("test", 10)
	assert.Exactly(t, 10, counter)

	ebus.Emit("unknown", 10)
	assert.Exactly(t, 10, counter)

	ebus.Emit("test", 5)
	assert.Exactly(t, 15, counter)
}
