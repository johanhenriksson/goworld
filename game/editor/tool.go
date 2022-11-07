package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Tool interface {
	object.T

	Use(T, vec3.T, vec3.T)
	Hover(T, vec3.T, vec3.T)
}
