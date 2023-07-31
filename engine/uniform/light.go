package uniform

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

type Light struct {
	ViewProj    [4]mat4.T
	Shadowmap   [4]uint32
	Distance    [4]float32
	Color       color.T
	Position    vec4.T
	Type        light.Type
	Intensity   float32
	Range       float32
	Attenuation light.Attenuation
}
