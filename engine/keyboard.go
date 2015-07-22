package engine

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type KeyCode glfw.Key
type KeyMap map[KeyCode]bool

/* Global key state */
var keyState KeyMap = KeyMap { }

/* Returns true if the given key is being held */
func KeyDown(key KeyCode) bool {
    return keyState[key]
}

/* GLFW Callback - Updates key state map */
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    code := KeyCode(key)
    keyState[code] = action != glfw.Release
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
