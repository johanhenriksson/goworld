package camera

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

// Camera Group
type Object struct {
	object.Object
	*Camera
}

// Camera Component
type Camera struct {
	object.Component
	Args

	Viewport    render.Screen
	Aspect      float32
	Proj        mat4.T
	View        mat4.T
	ViewInv     mat4.T
	ViewProj    mat4.T
	ViewProjInv mat4.T
	Eye         vec3.T
	Forward     vec3.T
}

type Args struct {
	Fov   float32
	Near  float32
	Far   float32
	Clear color.T
}

// New creates a new camera component.
func New(args Args) *Camera {
	return object.NewComponent(&Camera{
		Args:   args,
		Aspect: 1,
	})
}

func NewObject(args Args) *Object {
	return object.New("Camera", &Object{
		Camera: New(args),
	})
}

func (cam *Object) Name() string { return "Camera" }

// Unproject screen space coordinates into world space
func (cam *Camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	pos.Y = 1 - pos.Y
	pos = pos.Scaled(2).Sub(vec3.One)

	// unproject to world space by multiplying inverse view-projection
	return cam.ViewProjInv.TransformPoint(pos)
}

func (cam *Camera) PreDraw(args render.Args, scene object.Object) error {
	// update view & view-projection matrices
	cam.Viewport = args.Viewport
	cam.Aspect = float32(args.Viewport.Width) / float32(args.Viewport.Height)
	cam.Proj = mat4.Perspective(cam.Fov, cam.Aspect, cam.Near, cam.Far)

	// calculate the view matrix.
	// should be the inverse of the cameras transform matrix
	tf := cam.Transform()

	cam.Eye = tf.WorldPosition()
	cam.Forward = tf.Forward()

	cam.ViewInv = tf.Matrix()
	cam.View = cam.ViewInv.Invert()

	cam.ViewProj = cam.Proj.Mul(&cam.View)
	cam.ViewProjInv = cam.ViewProj.Invert()

	return nil
}
