package ebus

import (
	"container/list"
	"github.com/ivpusic/golog"
)

var log = golog.GetLogger("github.com/ivpusic/neo")

// Handler for event
type EvHandler func(interface{})

// Map with name of event and list of listeners.
type events map[string]list.List

type EBus struct {
	eventList events
}

func (b *EBus) InitEBus() {
	b.eventList = events{}
}

// Registering listener for provided event.
func (e *EBus) On(event string, fn EvHandler) {
	log.Debugf("registering event `%s`", event)
	topic, ok := e.eventList[event]

	if !ok {
		topic = list.List{}
	}

	topic.PushBack(fn)
	e.eventList[event] = topic
}

// Emitting event with data. One goroutine per emit will be created.
func (e *EBus) Emit(event string, data interface{}) {
	go func() {
		log.Debugf("emmiting event `%s`", event)

		if topic, ok := e.eventList[event]; ok {
			log.Debug("found listeners for event " + event)

			for el := topic.Front(); el != nil; el = el.Next() {
				fn := el.Value.(EvHandler)
				fn(data)
			}
		} else {
			log.Debugf("listeners for event %s not found", event)
		}
	}()
}
