package engine

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type KeyCode glfw.Key
type KeyState int
type KeyMap map[KeyCode]*Key

type Key struct {
    Code    KeyCode
    Pressed bool
}

func KeyDown(key KeyCode) bool {
    state, ok := keyState[key]
    return ok && state.Pressed
}

var keyState KeyMap = KeyMap { }

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    code := KeyCode(key)
    state, ok := keyState[code]

    if !ok {
        state = &Key {
            Code: code,
        }
        keyState[code] = state
    }

    switch action {
    case glfw.Press:
        state.Pressed = true
    case glfw.Release:
        state.Pressed = false
    }
}

/* GLFW Keycodes */
const (
    KeyA KeyCode = 65
    KeyB KeyCode = 66
    KeyC KeyCode = 67
    KeyD KeyCode = 68
    KeyE KeyCode = 69
    KeyF KeyCode = 70
    KeyG KeyCode = 71
    KeyH KeyCode = 72
    KeyI KeyCode = 73
    KeyJ KeyCode = 74
    KeyK KeyCode = 75
    KeyL KeyCode = 76
    KeyM KeyCode = 77
    KeyN KeyCode = 78
    KeyO KeyCode = 79
    KeyP KeyCode = 80
    KeyQ KeyCode = 81
    KeyR KeyCode = 82
    KeyS KeyCode = 83
    KeyT KeyCode = 84
    KeyU KeyCode = 85
    KeyV KeyCode = 86
    KeyW KeyCode = 87
    KeyX KeyCode = 88
    KeyY KeyCode = 89
    KeyZ KeyCode = 90
)
