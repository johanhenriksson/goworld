package engine

import (
    "fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type MouseButton glfw.MouseButton
type ButtonMap map[MouseButton]bool

var lastMouseX float32 = 0
var lastMouseY float32 = 0

var Mouse = MouseState {
    X: 0, Y: 0,
    DX: 0, DY: 0,
    lX: 0, lY: 0,
    init: false,
    buttons: ButtonMap { },
}

type MouseState struct {
    X, Y    float32
    DX, DY  float32
    lX, lY  float32
    init    bool
    buttons ButtonMap
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
    btn := MouseButton(button)
    Mouse.buttons[btn] = action != glfw.Release
    fmt.Println("Mouse", btn, Mouse.buttons[btn])
}

func MouseMoveCallback(w *glfw.Window, x float64, y float64) {
    Mouse.X = float32(x)
    Mouse.Y = float32(y)
}

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

func MouseDown(button MouseButton) bool {
    if s, ok := Mouse.buttons[button]; ok {
        return s
    }
    return false
}

const (
    /* Mouse Button Button Mouse Button */
    MouseButton1 MouseButton = MouseButton(glfw.MouseButton1)
    MouseButton2 MouseButton = MouseButton(glfw.MouseButton2)
    MouseButton3 MouseButton = MouseButton(glfw.MouseButton3)
    MouseButton4 MouseButton = MouseButton(glfw.MouseButton4)
)
