package srv

import (
	"fmt"
	"reflect"
)

type Event interface {
	Source() any
}

type EventHandler func(Event)

type Observable interface {
	Subscribe(any, EventHandler) func()
	Unsubscribe(any)
}

// Emitter is an observable that can emit events to subscribers
type Emitter struct {
	subscribers map[any]EventHandler
}

var _ Observable = &Emitter{}

func (e *Emitter) Subscribe(target any, sub EventHandler) func() {
	if e.subscribers == nil {
		e.subscribers = make(map[any]EventHandler)
	}
	e.subscribers[target] = sub
	return func() {
		e.Unsubscribe(target)
	}
}

func (e *Emitter) Unsubscribe(target any) {
	delete(e.subscribers, target)
}

func (e *Emitter) Emit(ev Event) {
	for _, callback := range e.subscribers {
		callback(ev)
	}
}

func EventBuffer() chan Event {
	return make(chan Event, 1024)
}

func Dump(object any) string {
	t := reflect.TypeOf(object)
	return fmt.Sprintf("%s %+v", t.Name(), object)
}
