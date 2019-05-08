package engine

import (
	"github.com/johanhenriksson/goworld/render"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Attenuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

type Light struct {
	Attenuation Attenuation
	Position    mgl.Vec3
	Color       mgl.Vec3
	Range       float32
	Type        LightType

	Projection mgl.Mat4 // Shadow projection matrix
	ProjWidth  float32
	ProjHeight float32
	ShadowMap  *render.Texture
}

type LightType int32

const PointLight LightType = 1
const DirectionalLight LightType = 2
