package keys

import "github.com/go-gl/glfw/v3.3/glfw"

type Modifier glfw.ModifierKey

const (
	NoMod = Modifier(0)
	Shift = Modifier(glfw.ModShift)
	Ctrl  = Modifier(glfw.ModControl)
	Alt   = Modifier(glfw.ModAlt)
	Super = Modifier(glfw.ModSuper)
)
