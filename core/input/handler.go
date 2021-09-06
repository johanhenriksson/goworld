package input

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
)

type Handler interface {
	KeyHandler
	MouseHandler
}

type KeyHandler interface {
	KeyEvent(keys.Event)
}

type MouseHandler interface {
	MouseEvent(mouse.Event)
}
