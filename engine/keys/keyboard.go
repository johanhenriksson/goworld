package keys

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

// Code represents a keyboard key
type Code glfw.Key

// KeyMap holds information about which keys are currently being held.
type KeyMap map[Code]bool

// Pressed returns true if the given key was just pressed.
func Pressed(key Code) bool {
	var current, last, ok bool
	if current, ok = keyState[key]; !ok {
		current = false
	}
	if last, ok = lastKeyState[key]; !ok {
		last = false
	}
	return current && !last
}

// Released returns true if the given key was just released.
func Released(key Code) bool {
	var current, last, ok bool
	if current, ok = keyState[key]; !ok {
		current = false
	}
	if last, ok = lastKeyState[key]; !ok {
		last = false
	}
	return !current && last
}

// Down returns true if the given key is currently being held down.
func Down(key Code) bool {
	if current, ok := keyState[key]; ok {
		return current
	}
	return false
}

// KeyCallback handles GLFW key callbacks to update the key state map.
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	code := Code(key)
	nextKeyState[code] = action != glfw.Release
}

// global key state
var keyState KeyMap = KeyMap{}
var nextKeyState KeyMap = KeyMap{}
var lastKeyState KeyMap = KeyMap{}

// Update updates key state maps at the end of each frame.
// Should only be called by the engine itself.
func Update(dt float32) {
	for k, v := range keyState {
		lastKeyState[k] = v
	}
	for k, v := range nextKeyState {
		keyState[k] = v
	}
}
