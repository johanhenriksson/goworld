package camera

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	object.Component

	Unproject(vec3.T) vec3.T
	View() mat4.T
	ViewInv() mat4.T
	Projection() mat4.T
	ViewProj() mat4.T
	ViewProjInv() mat4.T
	ClearColor() render.Color
}

// camera represents a 3D camera and its transform.
type camera struct {
	object.Component

	fov    float32
	aspect float32
	near   float32
	far    float32
	clear  render.Color

	proj  mat4.T
	view  mat4.T
	viewi mat4.T
	vp    mat4.T
	vpi   mat4.T
}

// New creates a new camera component.
func New(aspect, fov, near, far float32, clear render.Color) T {
	return &camera{
		Component: object.NewComponent(),

		aspect: aspect,
		fov:    fov,
		near:   near,
		far:    far,
		clear:  clear,

		proj: mat4.Perspective(math.DegToRad(fov), aspect, near, far),
	}
}

// Unproject screen space coordinates into world space
func (cam *camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	pos.Y = 1 - pos.Y
	point := pos.Scaled(2).Sub(vec3.One)

	// unproject to world space by multiplying inverse view-projection
	vpi := cam.vp.Invert()
	return vpi.TransformPoint(point)
}

// Update the camera
func (cam *camera) Update(dt float32) {
	// update view & view-projection matrices

	// Calculate new view matrix based on position & forward vector
	// why is this different from the parent objects world matrix?
	position := cam.Transform().WorldPosition()
	lookAt := position.Add(cam.Transform().Forward())
	cam.view = mat4.LookAt(position, lookAt)
	cam.viewi = cam.view.Invert()

	cam.vp = cam.proj.Mul(&cam.view)
	cam.vpi = cam.vp.Invert()
}

func (cam *camera) View() mat4.T        { return cam.view }
func (cam *camera) ViewInv() mat4.T     { return cam.viewi }
func (cam *camera) Projection() mat4.T  { return cam.proj }
func (cam *camera) ViewProj() mat4.T    { return cam.vp }
func (cam *camera) ViewProjInv() mat4.T { return cam.vpi }

func (cam *camera) ClearColor() render.Color { return cam.clear }
