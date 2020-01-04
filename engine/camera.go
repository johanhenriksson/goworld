package engine

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

// not sure if this camera should also support orthographic projection
// or if it's better to create a separate one

/* Perspective Camera. */
type Camera struct {
	*Transform
	Fov        float32
	Ratio      float32
	Near       float32
	Far        float32
	Buffer     *render.FrameBuffer
	Projection mgl.Mat4
	View       mgl.Mat4
	Clear      render.Color
}

func CreateCamera(buffer *render.FrameBuffer, x, y, z, fov, near, far float32) *Camera {
	ratio := float32(buffer.Width) / float32(buffer.Height)
	cam := &Camera{
		Transform:  CreateTransform(x, y, z),
		Buffer:     buffer,
		Ratio:      ratio,
		Fov:        fov,
		Near:       near,
		Far:        far,
		Projection: mgl.Perspective(mgl.DegToRad(fov), ratio, near, far),
		//Projection: mgl.Ortho(-width/2,width/2,-height/2,height/2,-100,100),
	}

	/* do an initial update at t=0 to initialize vectors */
	cam.Update(0.0)

	return cam
}

// todo
/* Project world space coordinates to screen space */
// func (cam *Camera) Project(mgl.Vec3) mgl.Vec2 { }

/* Unproject screen space coordinates into world space */
func (cam *Camera) Unproject(pos mgl.Vec3) mgl.Vec3 {
	// screen space -> clip space
	point := mgl.Vec4{
		pos.X()*2 - 1,
		(1-pos.Y())*2 - 1,
		pos.Z()*2 - 1,
		1.0,
	}

	/* Multiply by inverse view-projection matrix */
	pvi := cam.Projection.Mul4(cam.View)
	pvi = pvi.Inv()
	world := pvi.Mul4x1(point)

	/* World space coord */
	return mgl.Vec3{
		world.X() / world.W(),
		world.Y() / world.W(),
		world.Z() / world.W(),
	}
}

func (c *Camera) Update(dt float32) {
	/* Handle keyboard input */
	move := false
	dir := mgl.Vec3{}
	if KeyDown(KeyW) && !KeyDown(KeyS) {
		dir[2] += 1
		move = true
	}
	if KeyDown(KeyS) && !KeyDown(KeyW) {
		dir[2] -= 1
		move = true
	}
	if KeyDown(KeyA) && !KeyDown(KeyD) {
		dir[0] -= 1
		move = true
	}
	if KeyDown(KeyD) && !KeyDown(KeyA) {
		dir[0] += 1
		move = true
	}
	if KeyDown(KeyE) && !KeyDown(KeyQ) {
		dir[1] += 1
		move = true
	}
	if KeyDown(KeyQ) && !KeyDown(KeyE) {
		dir[1] -= 1
		move = true
	}

	if move {
		/* Calculate normalized movement vector */
		dv := 12.0 * dt /* magic number: movement speed */
		dir = dir.Normalize().Mul(dv)

		right := c.Transform.Right.Mul(dir[0])
		up := mgl.Vec3{0, dir[1], 0}
		forward := c.Transform.Forward.Mul(dir[2])

		/* Translate camera */
		c.Transform.Translate(right.Add(up.Add(forward)))
	}

	/* Mouse look */
	if MouseDown(MouseButton1) {
		rx := c.Transform.Rotation[0] - Mouse.DY*0.08
		ry := c.Transform.Rotation[1] - Mouse.DX*0.09

		/* Camera angle limits */
		/* -90 < rx < 90 */
		rx = float32(math.Max(-90.0, math.Min(90.0, float64(rx))))

		/* -180 < ry < 180 */
		if ry > 180.0 {
			ry -= 360.0
		}
		if ry < -180.0 {
			ry += 360.0
		}
		c.Transform.Rotation[0] = rx
		c.Transform.Rotation[1] = ry
	}

	/* Update transform with new position & rotation */
	c.Transform.Update(dt)

	/* Calculate new view matrix based on forward vector */
	lookAt := c.Transform.Position.Add(c.Transform.Forward)
	c.LookAt(lookAt)
}

func (c *Camera) LookAt(target mgl.Vec3) {
	c.View = mgl.LookAtV(c.Transform.Position, target, mgl.Vec3{0, 1, 0})
}

func (c *Camera) Use() {
	c.Buffer.Bind()
}
