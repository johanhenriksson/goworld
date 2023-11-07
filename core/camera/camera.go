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
	pos = pos.Scaled(2).Sub(vec3.One) // transforms from [0,1] to [-1,1]

	// unproject to world space by multiplying inverse view-projection
	return cam.ViewProjInv.TransformPoint(pos)
}

// Project world coordinates into screen space
func (cam *Camera) Project(pos vec3.T) vec3.T {
	p := cam.ViewProj.TransformPoint(pos)
	p = p.Add(vec3.One).Scaled(0.5) // transforms from [-1,1] to [0,1]
	return p.Mul(vec3.Extend(cam.Viewport.Size(), 1)).Scaled(1 / cam.Viewport.Scale)
}

func (cam *Camera) RenderArgs(screen render.Screen) render.Args {
	// todo: passing the global viewport allows the camera to modify the actual render viewport

	// update view & view-projection matrices
	cam.Viewport = screen
	cam.Aspect = float32(cam.Viewport.Width) / float32(cam.Viewport.Height)
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

	return render.Args{
		Viewport:   cam.Viewport,
		Near:       cam.Near,
		Far:        cam.Far,
		Fov:        cam.Fov,
		Projection: cam.Proj,
		View:       cam.View,
		ViewInv:    cam.ViewInv,
		VP:         cam.ViewProj,
		VPInv:      cam.ViewProjInv,
		MVP:        cam.ViewProj,
		Clear:      cam.Clear,
		Position:   cam.Transform().WorldPosition(),
		Forward:    cam.Transform().Forward(),
	}
}
