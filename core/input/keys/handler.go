package keys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Callback func(Event)

type Handler interface {
	KeyEvent(Event)
}

type FocusHandler interface {
	Handler
	FocusEvent()
	BlurEvent()
}

var focused FocusHandler

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

func Focus(handler FocusHandler) {
	if focused == handler {
		return
	}
	if focused != nil {
		focused.BlurEvent()
	}
	focused = handler
	if focused != nil {
		focused.FocusEvent()
	}
}
