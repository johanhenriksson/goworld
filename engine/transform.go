package engine

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
    Matrix      mgl.Mat4
    Position    mgl.Vec3
    Rotation    mgl.Vec3
    Scale       mgl.Vec3
    Forward     mgl.Vec3
    Right       mgl.Vec3
    Up          mgl.Vec3
}

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

func CreateTransform(x, y, z float32) *Transform {
    t := &Transform {
        Matrix:   mgl.Ident4(),
        Position: mgl.Vec3 { x,y,z },
        Rotation: mgl.Vec3 { 0,0,0 },
        Scale:    mgl.Vec3 { 1,1,1 },
    }
    t.Update(0)
    return t
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

func (t *Transform) Update(dt float32) {
    t.Position[1] += 0.5 * dt
    t.Matrix     = mgl.Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z())

    t.Right[0]   = t.Matrix[4*0+0]
    t.Right[1]   = t.Matrix[4*1+0]
    t.Right[2]   = t.Matrix[4*2+0]
    t.Up[0]      = t.Matrix[4*0+1]
    t.Up[1]      = t.Matrix[4*1+1]
    t.Up[2]      = t.Matrix[4*2+1]
    t.Forward[0] = -t.Matrix[4*0+2]
    t.Forward[1] = -t.Matrix[4*1+2]
    t.Forward[2] = -t.Matrix[4*2+2]
}

func (c *Camera) Update(dt float32) {
    c.Transform.Update(dt)

    lookAt := c.Transform.Position.Add(c.Transform.Forward);
	c.View = mgl.LookAtV(c.Transform.Position, lookAt, mgl.Vec3{0, 1, 0})
}
