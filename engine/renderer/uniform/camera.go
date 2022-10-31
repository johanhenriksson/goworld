package uniform

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
	Eye         vec3.T
}
