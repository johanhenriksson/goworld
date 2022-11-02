package keys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Handler interface {
	KeyEvent(Event)
}

var focused Handler

func KeyCallbackWrapper(handler Handler) glfw.KeyCallback {
	return func(
		w *glfw.Window,
		key glfw.Key,
		scancode int,
		action glfw.Action,
		mods glfw.ModifierKey,
	) {
		ev := &event{
			code:   Code(key),
			action: Action(action),
			mods:   Modifier(mods),
		}
		if focused != nil {
			focused.KeyEvent(ev)
		} else {
			handler.KeyEvent(ev)
		}
	}
}

func CharCallbackWrapper(handler Handler) glfw.CharCallback {
	return func(
		w *glfw.Window,
		char rune,
	) {
		ev := &event{
			char:   char,
			action: Char,
		}
		if focused != nil {
			focused.KeyEvent(ev)
		} else {
			handler.KeyEvent(ev)
		}
	}
}

func Focus(handler Handler) {
	focused = handler
}
