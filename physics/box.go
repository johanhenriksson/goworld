package physics

import (
    "fmt"
    "github.com/ianremmler/ode"
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Box struct {
    X   float32
    Y   float32
    Z   float32
    box ode.Box
}

func (w *World) NewBox(width, height, depth float32) *Box {
    col := w.space.NewBox(odeV3(width, height, depth))

    bx := &Box {
        X: width,
        Y: height,
        Z: depth,
        box: col,
    }

    /* store a handle to the box object so we can
     * send back collision events */
    col.SetData(bx)

    return bx
}

func (b *Box) String() string {
    return fmt.Sprintf("Box [w=%f, h=%f, d=%f]", b.X, b.Y, b.Z)
}

func (b *Box) SetPosition(position mgl.Vec3) {
    b.box.SetPosition(ToOdeVec3(position))
}

func (b *Box) AttachToBody(body *RigidBody) {
    b.box.SetBody(body.body)
}

func (b *Box) OnCollision(other Collider, contact Contact) {
    fmt.Println("Box Collision!", other)
}