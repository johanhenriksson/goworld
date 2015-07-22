package engine

import (
    "math"
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
    *Transform
    Width       float32
    Height      float32
    Fov         float32
    Ratio       float32
    Near        float32
    Far         float32
    Projection  mgl.Mat4
    View        mgl.Mat4
}

func CreateCamera(x, y, z, width, height, fov, near, far float32) *Camera {
    cam := &Camera {
        Transform: CreateTransform(x, y, z),
        Width: width,
        Height: height,
        Ratio: float32(width) / float32(height),
        Fov: fov,
        Near: near,
        Far: far,
        Projection: mgl.Perspective(mgl.DegToRad(fov), width/height, near, far),
    }
    cam.Update(0.0)
    return cam
}

func (c *Camera) Update(dt float32) {
    c.Transform.Update(dt)

    dv := 5.0 * dt

    move := false
    dir := mgl.Vec3 { }
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
        dir = dir.Normalize().Mul(dv)
        right := c.Transform.Right.Mul(dir[0])
        up := mgl.Vec3{0, dir[1], 0}
        forward := c.Transform.Forward.Mul(dir[2])
        c.Transform.Translate(right.Add(up.Add(forward)))
    }

    rx := c.Transform.Rotation[0] - Mouse.DY * 0.08
    ry := c.Transform.Rotation[1] - Mouse.DX * 0.09

    /* -90 < rx < 90 */
    /* -180 < ry < 180 */
    rx = float32(math.Max(-90.0, math.Min(90.0, float64(rx))))
    if ry > 180.0 {
        ry -= 360.0
    }
    if ry < -180.0 {
        ry += 360.0
    }

    if MouseDown(MouseButton1) {
        c.Transform.Rotation[0] = rx
        c.Transform.Rotation[1] = ry
    }

    lookAt := c.Transform.Position.Add(c.Transform.Forward);
	c.View = mgl.LookAtV(c.Transform.Position, lookAt, mgl.Vec3{0,1,0})
}
