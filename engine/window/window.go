package window

import (
	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/engine"
)

type ResizeHandler func(width, height int)

type Window interface {
	engine.Target

	Title() string
	SetTitle(string)

	Poll()
	ShouldClose() bool
	Destroy()

	SetInputHandler(input.Handler)
}

type WindowArgs struct {
	Title         string
	Width         int
	Height        int
	Frames        int
	Vsync         bool
	Debug         bool
	InputHandler  input.Handler
	ResizeHandler ResizeHandler
}
