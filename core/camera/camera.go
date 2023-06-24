package camera

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	object.T

	Unproject(vec3.T) vec3.T
	View() mat4.T
	ViewInv() mat4.T
	Projection() mat4.T
	ViewProj() mat4.T
	ViewProjInv() mat4.T
	ClearColor() color.T
	Near() float32
	Far() float32
	Fov() float32
	Viewport() render.Screen
}

// camera represents a 3D camera and its transform.
type camera struct {
	object.T

	viewport render.Screen
	aspect   float32
	fov      float32
	near     float32
	far      float32
	clear    color.T

	proj  mat4.T
	view  mat4.T
	viewi mat4.T
	vp    mat4.T
	vpi   mat4.T
	eye   vec3.T
	fwd   vec3.T
}

// New creates a new camera component.
func New(fov, near, far float32, clear color.T) T {
	return object.New(&camera{
		aspect: 1,
		fov:    fov,
		near:   near,
		far:    far,
		clear:  clear,
	})
}

func (cam *camera) Name() string { return "Camera" }

// Unproject screen space coordinates into world space
func (cam *camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	pos.Y = 1 - pos.Y
	pos = pos.Scaled(2).Sub(vec3.One)

	// unproject to world space by multiplying inverse view-projection
	return cam.vpi.TransformPoint(pos)
}

// Update the camera
func (cam *camera) Update(scene object.T, dt float32) {
}

func (cam *camera) PreDraw(args render.Args, scene object.T) error {
	// update view & view-projection matrices
	cam.viewport = args.Viewport
	cam.aspect = float32(args.Viewport.Width) / float32(args.Viewport.Height)
	cam.proj = mat4.Perspective(cam.fov, cam.aspect, cam.near, cam.far)

	// calculate the view matrix.
	// should be the inverse of the cameras transform matrix
	tf := cam.Transform()

	cam.eye = tf.WorldPosition()
	cam.fwd = tf.Forward()

	cam.viewi = tf.Matrix()
	cam.view = cam.viewi.Invert()

	cam.vp = cam.proj.Mul(&cam.view)
	cam.vpi = cam.vp.Invert()

	return nil
}

func (cam *camera) View() mat4.T            { return cam.view }
func (cam *camera) ViewInv() mat4.T         { return cam.viewi }
func (cam *camera) Projection() mat4.T      { return cam.proj }
func (cam *camera) ViewProj() mat4.T        { return cam.vp }
func (cam *camera) ViewProjInv() mat4.T     { return cam.vpi }
func (cam *camera) Near() float32           { return cam.near }
func (cam *camera) Far() float32            { return cam.far }
func (cam *camera) Fov() float32            { return cam.fov }
func (cam *camera) Viewport() render.Screen { return cam.viewport }

func (cam *camera) ClearColor() color.T { return cam.clear }
