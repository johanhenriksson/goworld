package engine

import (
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/engine/transform"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Camera represents a 3D camera and its transform.
type Camera struct {
	*transform.T
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
		T:          transform.New(position, vec3.Zero, vec3.One),
		Buffer:     buffer,
		Ratio:      ratio,
		Fov:        fov,
		Near:       near,
		Far:        far,
		Projection: mat4.Perspective(math.DegToRad(fov), ratio, near, far),
	}

	// do an initial update at t=0 to initialize vectors
	cam.Update(0.0)

	return cam
}

// Unproject screen space coordinates into world space
func (cam *Camera) Unproject(pos vec3.T) vec3.T {
	// screen space -> clip space
	pos.Y = 1 - pos.Y
	point := pos.Scaled(2).Sub(vec3.One)

	// unproject to world space by multiplying inverse view-projection
	pvi := cam.Projection.Mul(&cam.View)
	pvi = pvi.Invert()
	return pvi.TransformPoint(point)
}

// Update the camera
func (cam *Camera) Update(dt float32) {
	/* Mouse look */
	if mouse.Down(mouse.Button1) {
		sensitivity := vec2.New(0.08, 0.09)
		rot := cam.Rotation().XY().Sub(mouse.Delta.Swap().Mul(sensitivity))

		// camera angle limits
		rot.X = math.Clamp(rot.X, -89.9, 89.9)
		rot.Y = math.Mod(rot.Y, 360)

		cam.SetRotation(vec3.Extend(rot, 0))
	}

	// Update transform with new position & rotation
	cam.T.Update(nil)

	// update projection matrix in case aspect ratio changed
	ratio := float32(cam.Buffer.Width) / float32(cam.Buffer.Height)
	cam.Projection = mat4.Perspective(math.DegToRad(cam.Fov), ratio, cam.Near, cam.Far)

	// Calculate new view matrix based on forward vector
	lookAt := cam.Position().Add(cam.Forward())
	cam.View = mat4.LookAt(cam.Position(), lookAt)
}

// Use this camera for output.
func (cam *Camera) Use() {
	cam.Buffer.Bind()
}

func (cam *Camera) DrawArgs() DrawArgs {
	vp := cam.Projection.Mul(&cam.View)
	return DrawArgs{
		Projection: cam.Projection,
		View:       cam.View,
		VP:         vp,
		MVP:        vp,
		Transform:  mat4.Ident(),
		Position:   cam.Position(),
	}
}
