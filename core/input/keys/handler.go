package keys

import "github.com/go-gl/glfw/v3.1/glfw"

type Handler interface {
	KeyEvent(Event)
}

func KeyCallbackWrapper(handler Handler) glfw.KeyCallback {
	return func(
		w *glfw.Window,
		key glfw.Key,
		scancode int,
		action glfw.Action,
		mods glfw.ModifierKey) {

		handler.KeyEvent(&event{
			code:   Code(key),
			action: Action(action),
			mods:   Modifier(mods),
		})
	}
}

func CharCallbackWrapper(handler Handler) glfw.CharCallback {
	return func(
		w *glfw.Window,
		char rune) {

		handler.KeyEvent(&event{
			char:   char,
			action: Char,
		})
	}
}
