package ui

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
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

	Draw(render.DrawArgs)

	HandleMouse(MouseEvent) bool
	HandleKey(KeyEvent)
	HandleInput(rune)
	Focus()
	Blur()
}
