package mouse

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/johanhenriksson/goworld/math/vec2"
)

// Button refers to a mouse button.
type Button glfw.MouseButton

// ButtonMap holds information about the state of the mouse buttons.
type ButtonMap map[Button]bool

var Position = vec2.Zero
var Delta = vec2.Zero

var last = vec2.Zero
var initialized = false
var buttons = make([]bool, 32)
var nextbuttons = make([]bool, 32)
var lastbuttons = make([]bool, 32)

// Pressed returns true if the given button was just pressed down.
func Pressed(button Button) bool {
	return buttons[button] && !lastbuttons[button]
}

// Down returns true if the given button is being held down.
func Down(button Button) bool {
	return buttons[button]
}

// Up returns true if the given button is not being held down.
func Up(button Button) bool {
	return !buttons[button]
}

// ButtonCallback updates the mouse button state map
func ButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	btn := Button(button)
	nextbuttons[btn] = action != glfw.Release
}

// MoveCallback updates the cursor position
func MoveCallback(w *glfw.Window, x, y float64, scale float32) {
	Position = vec2.New(float32(x), float32(y)).Scaled(scale)
}

// Update updates mouse delta x/y on every frame
func Update(dt float32) {
	if initialized {
		Delta = last.Sub(Position)
	} else {
		initialized = true
	}
	last = Position

	copy(lastbuttons[:], buttons[:])
	copy(buttons[:], nextbuttons[:])
}

const (
	// Button1 refers to Mouse Button 1
	Button1 Button = Button(glfw.MouseButton1)

	// Button2 refers to Mouse Button 2
	Button2 Button = Button(glfw.MouseButton2)

	// Button3 refers to Mouse Button 3
	Button3 Button = Button(glfw.MouseButton3)

	// Button4 refers to Mouse Button 4
	Button4 Button = Button(glfw.MouseButton4)
)
