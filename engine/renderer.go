package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Renderer interface {
	Draw(args render.Args, scene object.T)
	Buffers() BufferOutput
	Destroy()
}
type BufferOutput interface {
	SamplePosition(cursor vec2.T) (vec3.T, bool)
	SampleNormal(cursor vec2.T) (vec3.T, bool)
}
