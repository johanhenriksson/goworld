package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Tool interface {
	object.Component
	Use(*Editor, vec3.T, vec3.T)
	Hover(*Editor, vec3.T, vec3.T)
}
