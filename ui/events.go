package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
)

type MouseHandler func(MouseEvent)

// MouseEvent represents a mouse event as it propagates through the component hierarchy
type MouseEvent struct {
	UI     *Manager
	Point  mgl.Vec2
	Button mouse.Button
}

// KeyEvent represents raw key event
type KeyEvent struct {
	UI    *Manager
	Key   keys.Code
	Press bool
}
