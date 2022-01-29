package render

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/texture"
)

// Attenuation properties for point lights
type Attenuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

// Light holds information about a generic light
type Light struct {
	Attenuation Attenuation
	Position    vec3.T
	Color       vec3.T
	Range       float32
	Intensity   float32
	Type        LightType
	Shadows     bool

	Projection mat4.T // Shadow projection matrix
	ShadowMap  texture.T
}

// LightType indicates which kind of light. Point, Directional etc
type LightType int32

// AmbientLight is the background light applied to everything.
const AmbientLight LightType = 0

// PointLight is a normal light casting rays in all directions.
const PointLight LightType = 1

// DirectionalLight is a directional light source, casting parallell rays.
const DirectionalLight LightType = 2
