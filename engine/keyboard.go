package engine

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

// KeyCode represents a keyboard key
type KeyCode glfw.Key

// KeyMap holds information about which keys are currently being held.
type KeyMap map[KeyCode]bool

// KeyPressed returns true if the given key was just pressed.
func KeyPressed(key KeyCode) bool {
	var current, last, ok bool
	if current, ok = keyState[key]; !ok {
		current = false
	}
	if last, ok = lastKeyState[key]; !ok {
		last = false
	}
	return current && !last
}

// KeyReleased returns true if the given key was just released.
func KeyReleased(key KeyCode) bool {
	var current, last, ok bool
	if current, ok = keyState[key]; !ok {
		current = false
	}
	if last, ok = lastKeyState[key]; !ok {
		last = false
	}
	return !current && last
}

// KeyDown returns true if the given key is currently being held down.
func KeyDown(key KeyCode) bool {
	if current, ok := keyState[key]; ok {
		return current
	}
	return false
}

// KeyCallback handles GLFW key callbacks to update the key state map.
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	code := KeyCode(key)
	nextKeyState[code] = action != glfw.Release
}

// global key state
var keyState KeyMap = KeyMap{}
var nextKeyState KeyMap = KeyMap{}
var lastKeyState KeyMap = KeyMap{}

// updateKeyboard updates key state maps at the end of each frame.
// Should only be called by the engine itself.
func updateKeyboard(dt float32) {
	for k, v := range keyState {
		lastKeyState[k] = v
	}
	for k, v := range nextKeyState {
		keyState[k] = v
	}
}

// GLFW Keycodes
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

	KeyEnter     = KeyCode(glfw.KeyEnter)
	KeyEscape    = KeyCode(glfw.KeyEscape)
	KeyBackspace = KeyCode(glfw.KeyBackspace)
)
