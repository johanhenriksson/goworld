package light

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

// Attenuation properties for point lights
type Attenuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

// T holds information about a generic light
type T struct {
	Attenuation Attenuation
	Position    vec3.T
	Color       color.T
	Range       float32
	Intensity   float32
	Type        Type
	Shadows     bool

	Projection mat4.T // Shadow projection matrix
	ShadowMap  texture.T
}

// Type indicates which kind of light. Point, Directional etc
type Type int32

// AmbientLight is the background light applied to everything.
const Ambient Type = 0

// PointLight is a normal light casting rays in all directions.
const Point Type = 1

// DirectionalLight is a directional light source, casting parallell rays.
const Directional Type = 2

var DefaultAttenuation = Attenuation{
	Constant:  1.0,
	Linear:    0.35,
	Quadratic: 0.44,
}
