package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	object.T

	Type() Type
	LightDescriptor(args render.Args) Descriptor
}

// Descriptor holds rendering information for lights
type Descriptor struct {
	Projection  mat4.T // Light projection matrix
	View        mat4.T // Light view matrix
	ViewProj    mat4.T
	Color       color.T
	Position    vec4.T
	Type        Type
	Range       float32
	Intensity   float32
	Shadows     int32
	Attenuation Attenuation
}

// Attenuation properties for point lights
type Attenuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

var DefaultAttenuation = Attenuation{
	Constant:  1.0,
	Linear:    0.35,
	Quadratic: 0.44,
}
