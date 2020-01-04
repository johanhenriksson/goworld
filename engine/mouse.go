package engine

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
)

// MouseButton refers to a mouse button.
type MouseButton glfw.MouseButton

// ButtonMap holds information about the state of the mouse buttons.
type ButtonMap map[MouseButton]bool

// Mouse contains the current state of the mouse.
var Mouse = MouseState{
	buttons:     make([]bool, 32),
	nextbuttons: make([]bool, 32),
	lastbuttons: make([]bool, 32),
}

// MouseState holds the current mouse state
type MouseState struct {
	X, Y        float32
	DX, DY      float32 // frame delta x, y
	lX, lY      float32 // last x, y
	init        bool
	buttons     []bool
	nextbuttons []bool
	lastbuttons []bool
}

// MouseDownPress returns true if the given button was just pressed down.
func MouseDownPress(button MouseButton) bool {
	return Mouse.buttons[button] && !Mouse.lastbuttons[button]
}

// MouseDown returns true if the given button is being held down.
func MouseDown(button MouseButton) bool {
	return Mouse.buttons[button]
}

// MouseButtonCallback updates the mouse button state map
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	btn := MouseButton(button)
	fmt.Println("press mouse", btn)
	Mouse.nextbuttons[btn] = action != glfw.Release
}

// MouseMoveCallback updates the cursor position
func MouseMoveCallback(w *glfw.Window, x, y float64, scale float32) {
	Mouse.X = float32(x) * scale
	Mouse.Y = float32(y) * scale
}

// UpdateMouse updates mouse delta x/y on every frame
func updateMouse(dt float32) {
	if Mouse.init {
		Mouse.DX = Mouse.lX - Mouse.X
		Mouse.DY = Mouse.lY - Mouse.Y
	} else {
		Mouse.init = true
	}
	Mouse.lX = Mouse.X
	Mouse.lY = Mouse.Y

	copy(Mouse.lastbuttons[:], Mouse.buttons[:])
	copy(Mouse.buttons[:], Mouse.nextbuttons[:])
}

const (
	// MouseButton1 refers to Mouse Button 1
	MouseButton1 MouseButton = MouseButton(glfw.MouseButton1)

	// MouseButton2 refers to Mouse Button 2
	MouseButton2 MouseButton = MouseButton(glfw.MouseButton2)

	// MouseButton3 refers to Mouse Button 3
	MouseButton3 MouseButton = MouseButton(glfw.MouseButton3)

	// MouseButton4 refers to Mouse Button 4
	MouseButton4 MouseButton = MouseButton(glfw.MouseButton4)
)
