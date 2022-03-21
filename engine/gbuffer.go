package engine

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type BufferOutput interface {
	SamplePosition(cursor vec2.T) (vec3.T, bool)
	SampleNormal(cursor vec2.T) (vec3.T, bool)
}
