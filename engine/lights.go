package engine

import (
	"github.com/johanhenriksson/goworld/render"

	mgl "github.com/go-gl/mathgl/mgl32"
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
	Position    mgl.Vec3
	Color       mgl.Vec3
	Range       float32
	Intensity   float32
	Type        LightType
	Shadows     bool

	Projection mgl.Mat4 // Shadow projection matrix
	ProjWidth  float32
	ProjHeight float32
	ShadowMap  *render.Texture
}

// LightType indicates which kind of light. Point, Directional etc
type LightType int32

const AmbientLight LightType = 0

// PointLight is a normal light casting rays in all directions.
const PointLight LightType = 1

// DirectionalLight is a directional light source, casting parallell rays.
const DirectionalLight LightType = 2
