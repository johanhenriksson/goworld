package physics

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/ode"
)

type RigidBox struct {
	*RigidBody
	collider Collider
	center   mgl.Vec3

	/* Size */
	X float32
	Y float32
	Z float32
}

type RigidBody struct {
	Mass float32
	body ode.Body
	mass *ode.Mass
}

func (w *World) NewRigidBody(mass float32) *RigidBody {
	m := ode.NewMass()
	m.Adjust(float64(mass))

	body := w.world.NewBody()

	rb := &RigidBody{
		Mass: mass,
		body: body,
		mass: m,
	}
	return rb
}

func (rb *RigidBody) Position() mgl.Vec3 {
	return FromOdeVec3(rb.body.Position())
}

func (rb *RigidBody) SetPosition(position mgl.Vec3) {
	rb.body.SetPosition(ToOdeVec3(position))
}

func (rb *RigidBody) Rotation() mgl.Vec3 {
	return FromOdeRotation(rb.body.Rotation())
}

func (rb *RigidBody) setBox(density, x, y, z float32) {
	rb.mass.SetBox(float64(density), odeV3(x, y, z))
	rb.mass.Adjust(float64(rb.Mass))
}

func (w *World) NewRigidBox(mass, x, y, z float32) *RigidBox {
	rb := w.NewRigidBody(mass)
	rb.setBox(1, x, y, z)

	col := w.NewBox(x, y, z)
	col.AttachToBody(rb)

	box := &RigidBox{
		RigidBody: rb,
		collider:  col,
		X:         x, Y: y, Z: z,
		center: mgl.Vec3{x / 2, y / 2, z / 2},
	}

	return box
}

func (rb *RigidBox) Position() mgl.Vec3 {
	return rb.RigidBody.Position().Sub(rb.center)
}

func (rb *RigidBox) SetPosition(position mgl.Vec3) {
	rb.body.SetPosition(ToOdeVec3(position.Add(rb.center)))
}
