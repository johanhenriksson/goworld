package widget

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
)

func Find(widget T, key string) T {
	if widget.Key() == key {
		return widget
	}
	for _, child := range widget.Children() {
		if hit := Find(child, key); hit != nil {
			return hit
		}
	}
	return nil
}

func SimulateClick(widget T, button mouse.Button) {
	handler, ok := widget.(mouse.Handler)
	if !ok {
		panic("widget does not implement mouse.Handler")
	}

	press := mouse.NewButtonEvent(button, mouse.Press, widget.Position(), 0, false)
	handler.MouseEvent(press)
	release := mouse.NewButtonEvent(button, mouse.Release, widget.Position(), 0, false)
	handler.MouseEvent(release)
}

type MouseHandler interface {
	MouseEvent(mouse.Event) (MouseHandler, float32)

	MouseEnter(e mouse.Event)
	MouseExit(e mouse.Event)
	MouseMove(e mouse.Event)
	MouseClick(e mouse.Event)
}
