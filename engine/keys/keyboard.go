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

// GLFW Keycodes
const (
	A Code = 65
	B Code = 66
	C Code = 67
	D Code = 68
	E Code = 69
	F Code = 70
	G Code = 71
	H Code = 72
	I Code = 73
	J Code = 74
	K Code = 75
	L Code = 76
	M Code = 77
	N Code = 78
	O Code = 79
	P Code = 80
	Q Code = 81
	R Code = 82
	S Code = 83
	T Code = 84
	U Code = 85
	V Code = 86
	W Code = 87
	X Code = 88
	Y Code = 89
	Z Code = 90

	Enter       = Code(glfw.KeyEnter)
	Escape      = Code(glfw.KeyEscape)
	Backspace   = Code(glfw.KeyBackspace)
	Space       = Code(glfw.KeySpace)
	LeftShift   = Code(glfw.KeyLeftShift)
	LeftControl = Code(glfw.KeyLeftControl)
)
