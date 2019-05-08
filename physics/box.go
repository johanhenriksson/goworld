package physics

import (
	"fmt"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/ode"
)

type Box struct {
	X        float32
	Y        float32
	Z        float32
	Center   mgl.Vec3
	box      ode.Box
	Callback CollisionCallback
}

func (w *World) NewBox(width, height, depth float32) *Box {
	col := w.space.NewBox(odeV3(width, height, depth))

	bx := &Box{
		X: width,
		Y: height,
		Z: depth,
		Center: mgl.Vec3{
			width / 2,
			height / 2,
			depth / 2,
		},
		box: col,
	}

	/* store a handle to the box object so we can
	 * send back collision events */
	col.SetData(bx)

	return bx
}

func (b *Box) String() string {
	return fmt.Sprintf("Box [w=%.1f, h=%.1f, d=%.1f]", b.X, b.Y, b.Z)
}

func (b *Box) SetPosition(position mgl.Vec3) {
	b.box.SetPosition(ToOdeVec3(position))
}

func (b *Box) AttachToBody(body *RigidBody) {
	b.box.SetBody(body.body)
	//b.box.SetOffsetPosition(ToOdeVec3(b.Center))
}

func (b *Box) Destroy() {
	b.box.Destroy()
}

func (b *Box) OnCollision(other Collider, contact Contact) {
	if b.Callback != nil {
		b.Callback(other, contact)
	}
}
