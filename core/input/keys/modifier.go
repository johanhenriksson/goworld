package keys

import "github.com/go-gl/glfw/v3.1/glfw"

type Modifier glfw.ModifierKey

const (
	Shift = Modifier(glfw.ModSuper)
	Ctrl  = Modifier(glfw.ModControl)
	Alt   = Modifier(glfw.ModAlt)
	Super = Modifier(glfw.ModSuper)
)
