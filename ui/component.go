package ui

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Component interface {
	GetStyle() Style

	Width() float32
	Height() float32
	ZIndex() float32
	Resize(vec2.T) vec2.T
	Flow(vec2.T) vec2.T
	SetPosition(vec2.T)
	Children() []Component

	Draw(engine.DrawArgs)

	HandleMouse(MouseEvent) bool
	HandleKey(KeyEvent)
	HandleInput(rune)
	Focus()
	Blur()
}
