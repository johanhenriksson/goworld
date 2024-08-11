package draw

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Camera struct {
	Proj        mat4.T
	View        mat4.T
	ViewProj    mat4.T
	ProjInv     mat4.T
	ViewInv     mat4.T
	ViewProjInv mat4.T
	Position    vec3.T
	Forward     vec3.T
	Viewport    Viewport
	Near        float32
	Far         float32
	Aspect      float32
	Fov         float32
}
