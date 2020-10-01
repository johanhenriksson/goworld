package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// Component is the general interface for scene object components.
type Component interface {
	Update(float32)
	Draw(render.DrawArgs)
	Base() *Object
}
