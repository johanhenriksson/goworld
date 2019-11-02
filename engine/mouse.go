package engine

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type MouseButton glfw.MouseButton
type ButtonMap map[MouseButton]bool

/* Mouse Global */
var Mouse = MouseState{
	buttons:     make([]bool, 4),
	lastbuttons: make([]bool, 4),
}

type MouseState struct {
	X, Y        float32
	DX, DY      float32 // frame delta x, y
	lX, lY      float32 // last x, y
	init        bool
	buttons     []bool
	lastbuttons []bool
}

/* Returns true if the given mouse button is held down */
func MouseDownPress(button MouseButton) bool {
	return Mouse.buttons[button] && !Mouse.lastbuttons[button]
}

func MouseDown(button MouseButton) bool {
	return Mouse.buttons[button]
}

/* GLFW Callback - Update mouse state map */
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	btn := MouseButton(button)
	Mouse.buttons[btn] = action != glfw.Release
}

/* GLFW Callback - Update mouse coords */
func MouseMoveCallback(w *glfw.Window, x, y float64, scale float32) {
	Mouse.X = float32(x) * scale
	Mouse.Y = float32(y) * scale
}

/* Updates mouse delta x/y every frame */
func UpdateMouse(dt float32) {
	if Mouse.init {
		Mouse.DX = Mouse.lX - Mouse.X
		Mouse.DY = Mouse.lY - Mouse.Y
	} else {
		Mouse.init = true
	}
	Mouse.lX = Mouse.X
	Mouse.lY = Mouse.Y

	copy(Mouse.lastbuttons[:], Mouse.buttons[:])
}

const (
	/* Mouse Button Button Mouse Button */
	MouseButton1 MouseButton = MouseButton(glfw.MouseButton1)
	MouseButton2 MouseButton = MouseButton(glfw.MouseButton2)
	MouseButton3 MouseButton = MouseButton(glfw.MouseButton3)
	MouseButton4 MouseButton = MouseButton(glfw.MouseButton4)
)
