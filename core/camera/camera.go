package camera

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

// Camera Group
type Object struct {
	object.Object
	Camera *Camera
}

// Camera Component
type Camera struct {
	object.Component

	Fov  object.Property[float32]
	Near object.Property[float32]
	Far  object.Property[float32]

	state draw.Camera
}

type Args struct {
	Fov   float32
	Near  float32
	Far   float32
	Clear color.T
}

// New creates a new camera component.
func New(pool object.Pool, args Args) *Camera {
	return object.NewComponent(pool, &Camera{
		Fov:  object.NewProperty(args.Fov),
		Near: object.NewProperty(args.Near),
		Far:  object.NewProperty(args.Far),
	})
}

func NewObject(pool object.Pool, args Args) *Object {
	return object.NewObject(pool, "Camera", &Object{
		Camera: New(pool, args),
	})
}

func (cam *Object) Name() string { return "Camera" }

// Unproject screen space coordinates into world space
func (cam *Camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	pos.Y = 1 - pos.Y
	pos = pos.Scaled(2).Sub(vec3.One) // transforms from [0,1] to [-1,1]

	// unproject to world space by multiplying inverse view-projection
	return cam.state.ViewProjInv.TransformPoint(pos)
}

// Project world coordinates into screen space
func (cam *Camera) Project(pos vec3.T) vec3.T {
	p := cam.state.ViewProj.TransformPoint(pos)
	p = p.Add(vec3.One).Scaled(0.5) // transforms from [-1,1] to [0,1]
	scale := float32(1)
	if cam.state.Viewport.Scale > 1 {
		scale /= cam.state.Viewport.Scale
	}
	return p.Mul(vec3.Extend(cam.state.Viewport.Size(), 1)).Scaled(scale)
}

// Recalculate camera matrices based on the current transform and viewport
func (cam *Camera) Refresh(viewport draw.Viewport) draw.Camera {
	// todo: passing the global viewport allows the camera to modify the actual render viewport

	cam.state.Viewport = viewport
	cam.state.Aspect = viewport.Aspect()
	cam.state.Near = cam.Near.Get()
	cam.state.Far = cam.Far.Get()
	cam.state.Fov = cam.Fov.Get()

	// update view & view-projection matrices
	cam.state.Proj = mat4.Perspective(cam.state.Fov, cam.state.Aspect, cam.state.Near, cam.state.Far)

	// calculate the view matrix.
	// should be the inverse of the cameras transform matrix
	tf := cam.Transform()

	cam.state.Position = tf.WorldPosition()
	cam.state.Forward = tf.Forward()

	cam.state.ViewInv = tf.Matrix()
	cam.state.View = cam.state.ViewInv.Invert()

	cam.state.ViewProj = cam.state.Proj.Mul(&cam.state.View)
	cam.state.ViewProjInv = cam.state.ViewProj.Invert()

	return cam.state
}
