package engine

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type MouseButton glfw.MouseButton
type ButtonMap map[MouseButton]bool

/* Mouse Global */
var Mouse = MouseState {
    buttons: ButtonMap { },
}

type MouseState struct {
    X, Y    float32
    DX, DY  float32 // frame delta x, y
    lX, lY  float32 // last x, y
    init    bool
    buttons ButtonMap
}

/* Returns true if the given mouse button is held down */
func MouseDown(button MouseButton) bool {
    return Mouse.buttons[button]
}

/* GLFW Callback - Update mouse state map */
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
    btn := MouseButton(button)
    Mouse.buttons[btn] = action != glfw.Release
}

/* GLFW Callback - Update mouse coords */
func MouseMoveCallback(w *glfw.Window, x float64, y float64) {
    Mouse.X = float32(x)
    Mouse.Y = float32(y)
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
}

const (
    /* Mouse Button Button Mouse Button */
    MouseButton1 MouseButton = MouseButton(glfw.MouseButton1)
    MouseButton2 MouseButton = MouseButton(glfw.MouseButton2)
    MouseButton3 MouseButton = MouseButton(glfw.MouseButton3)
    MouseButton4 MouseButton = MouseButton(glfw.MouseButton4)
)
