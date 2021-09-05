package mouse

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Handler interface {
	MouseEvent(Event)
}

type MouseWrapper interface {
	Button(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey)
	Move(w *glfw.Window, x, y float64)
	Scroll(w *glfw.Window, x, y float64)
}

type wrapper struct {
	Handler
	position vec2.T
}

func NewWrapper(handler Handler) MouseWrapper {
	return &wrapper{
		Handler: handler,
	}
}

func (mw *wrapper) Button(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	mw.MouseEvent(&event{
		action:   Action(action),
		button:   Button(button),
		mods:     keys.Modifier(mod),
		position: mw.position,
	})
}

func (mw *wrapper) Move(w *glfw.Window, x, y float64) {
	// calculate framebuffer scale relative to window
	width, _ := w.GetSize()
	fwidth, fheight := w.GetFramebufferSize()
	scale := float32(fwidth) / float32(width)

	// calculate framebuffer position & mouse delta
	pos := vec2.New(float32(x), float32(y)).Scaled(scale)
	dt := pos.Sub(mw.position)
	mw.position = pos

	// discard events that occur outside of the window bounds
	if pos.X < 0 || pos.Y < 0 || int(pos.X) > fwidth || int(pos.Y) > fheight {
		return
	}

	// submit event to handler
	mw.MouseEvent(&event{
		action:   Move,
		position: pos,
		delta:    dt,
	})
}

func (mw *wrapper) Scroll(w *glfw.Window, x, y float64) {
	mw.MouseEvent(&event{
		action: Scroll,
		scroll: vec2.New(float32(x), float32(y)),
	})
}
