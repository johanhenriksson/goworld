package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
)

type MouseHandler func(MouseEvent)

// MouseEvent represents a mouse event as it propagates through the component hierarchy
type MouseEvent struct {
	UI     *Manager
	Point  mgl.Vec2
	Button engine.MouseButton
}

// KeyEvent represents raw key event
type KeyEvent struct {
	UI    *Manager
	Key   engine.KeyCode
	Press bool
}

type Component interface {
	Draw(render.DrawArgs)
	HandleMouse(MouseEvent) bool
	HandleKey(KeyEvent)
	HandleInput(rune)
	Focus()
	Blur()
}
