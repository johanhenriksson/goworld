package keys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Code represents a keyboard key
type Code glfw.Key

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

	Key0 = Code(glfw.Key0)
	Key1 = Code(glfw.Key1)
	Key2 = Code(glfw.Key2)
	Key3 = Code(glfw.Key3)
	Key4 = Code(glfw.Key4)
	Key5 = Code(glfw.Key5)
	Key6 = Code(glfw.Key6)
	Key7 = Code(glfw.Key7)
	Key8 = Code(glfw.Key8)
	Key9 = Code(glfw.Key9)

	Enter        = Code(glfw.KeyEnter)
	Escape       = Code(glfw.KeyEscape)
	Backspace    = Code(glfw.KeyBackspace)
	Delete       = Code(glfw.KeyDelete)
	Space        = Code(glfw.KeySpace)
	LeftShift    = Code(glfw.KeyLeftShift)
	RightShift   = Code(glfw.KeyRightShift)
	LeftControl  = Code(glfw.KeyLeftControl)
	RightControl = Code(glfw.KeyRightControl)
	LeftAlt      = Code(glfw.KeyLeftAlt)
	RightAlt     = Code(glfw.KeyRightAlt)
	LeftSuper    = Code(glfw.KeyLeftSuper)
	RightSuper   = Code(glfw.KeyRightSuper)
	LeftArrow    = Code(glfw.KeyLeft)
	RightArrow   = Code(glfw.KeyRight)
	UpArrow      = Code(glfw.KeyUp)
	DownArrow    = Code(glfw.KeyDown)
	NumpadEnter  = Code(glfw.KeyKPEnter)
)
