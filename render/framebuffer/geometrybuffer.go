package framebuffer

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Geometry interface {
	T
	Diffuse() texture.T
	Normal() texture.T
	Position() texture.T
	Depth() texture.T
	SampleNormal(vec2.T) (vec3.T, bool)
}
