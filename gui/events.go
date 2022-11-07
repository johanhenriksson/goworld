package gui

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
)

// Dummy mouse handler that consumes the event and does nothing.
func ConsumeMouse(ev mouse.Event) {
	ev.Consume()
}

// Dummy key handler that consumes the event and does nothing.
func ConsumeKey(ev keys.Event) {}
