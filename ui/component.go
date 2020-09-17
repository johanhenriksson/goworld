package ui

import (
	"github.com/johanhenriksson/goworld/render"
)

type Component interface {
	GetStyle() Style

	Width() float32
	Height() float32
	ZIndex() float32
	Resize(Size) Size
	Flow(Size) Size
	SetPosition(float32, float32)
	Children() []Component

	Draw(render.DrawArgs)

	HandleMouse(MouseEvent) bool
	HandleKey(KeyEvent)
	HandleInput(rune)
	Focus()
	Blur()
}
