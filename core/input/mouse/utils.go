package mouse

import "github.com/go-gl/glfw/v3.1/glfw"

func Lock() {
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func Hide() {
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func Show() {
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}
