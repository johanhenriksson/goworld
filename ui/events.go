package ui

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type MouseHandler func(MouseEvent)

// MouseEvent represents a mouse event as it propagates through the component hierarchy
type MouseEvent struct {
	UI     *Manager
	Point  vec2.T
	Button mouse.Button
}

// KeyEvent represents raw key event
type KeyEvent struct {
	UI    *Manager
	Key   keys.Code
	Press bool
}
