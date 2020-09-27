package engine

import (
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
)

// not sure if this camera should also support orthographic projection
// or if it's better to create a separate one

// Camera represents a 3D camera and its transform.
type Camera struct {
	*Transform
	Fov        float32
	Ratio      float32
	Near       float32
	Far        float32
	Buffer     *render.FrameBuffer
	Projection mat4.T
	View       mat4.T
	Clear      render.Color
}

// CreateCamera creates a new camera object.
func CreateCamera(buffer *render.FrameBuffer, position vec3.T, fov, near, far float32) *Camera {
	ratio := float32(buffer.Width) / float32(buffer.Height)
	cam := &Camera{
		Transform:  CreateTransform(position),
		Buffer:     buffer,
		Ratio:      ratio,
		Fov:        fov,
		Near:       near,
		Far:        far,
		Projection: mat4.Perspective(math.DegToRad(fov), ratio, near, far),
	}

	/* do an initial update at t=0 to initialize vectors */
	cam.Update(0.0)

	return cam
}

// todo
/* Project world space coordinates to screen space */
// func (cam *Camera) Project(mgl.Vec3) mgl.Vec2 { }

// Unproject screen space coordinates into world space
func (cam *Camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	point := vec4.T{
		pos.X*2 - 1,
		(1-pos.Y)*2 - 1,
		pos.Z*2 - 1,
		1.0,
	}

	/* Multiply by inverse view-projection matrix */
	pvi := cam.Projection.Mul(&cam.View)
	pvi = pvi.Inv()
	world := pvi.MulVec4(point)

	/* World space coord */
	return vec3.T{
		world.X / world.W,
		world.Y / world.W,
		world.Z / world.W,
	}
}

// Update the camera
func (cam *Camera) Update(dt float32) {
	/* Mouse look */
	if mouse.Down(mouse.Button1) {
		rx := cam.Transform.Rotation.X - mouse.DY*0.08
		ry := cam.Transform.Rotation.Y - mouse.DX*0.09

		/* Camera angle limits */
		/* -90 < rx < 90 */
		rx = math.Max(-90.0, math.Min(90.0, rx))

		/* -180 < ry < 180 */
		if ry > 180.0 {
			ry -= 360.0
		}
		if ry < -180.0 {
			ry += 360.0
		}
		cam.Transform.Rotation.X = rx
		cam.Transform.Rotation.Y = ry
	}

	/* Update transform with new position & rotation */
	cam.Transform.Update(dt)

	/* Calculate new view matrix based on forward vector */
	lookAt := cam.Transform.Position.Add(cam.Transform.Forward)
	cam.LookAt(lookAt)
}

// LookAt orients the camera towards a point in 3D space.
func (cam *Camera) LookAt(target vec3.T) {
	cam.View = mat4.LookAt(cam.Transform.Position, target)
}

// Use this camera for output.
func (cam *Camera) Use() {
	cam.Buffer.Bind()
}
