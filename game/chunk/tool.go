package chunk

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Tool interface {
	object.T

	Use(Editor, vec3.T, vec3.T)
	Hover(Editor, vec3.T, vec3.T)
}
