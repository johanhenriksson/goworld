package mouse

import "github.com/go-gl/glfw/v3.1/glfw"

var locked bool = false
var lockingEnabled bool = false

func Lock() {
	// actual cursor locking can be awkward, so leave an option to enable it
	// otherwise the cursor will be locked virtually - i.e. only in the sense that
	// mouse events have the Locked flag set to true
	if lockingEnabled {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}
	locked = true
}

func Hide() {
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func Show() {
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	locked = false
}
