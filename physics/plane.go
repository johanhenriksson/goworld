package physics

import (
	"fmt"
	"github.com/johanhenriksson/ode"
)

type Plane struct {
	X     float32
	Y     float32
	Z     float32
	C     float32
	plane ode.Plane
}

func (w *World) NewPlane(x, y, z, c float32) *Plane {
	col := w.space.NewPlane(odeV4(x, y, z, c))

	p := &Plane{
		X:     x,
		Y:     y,
		Z:     z,
		C:     c,
		plane: col,
	}

	/* store a handle to the box object so we can
	 * send back collision events */
	col.SetData(p)

	return p
}

func (p *Plane) String() string {
	return fmt.Sprintf("Plane [x=%.1f, y=%.1f, z=%.1f, c=%.1f]", p.X, p.Y, p.Z, p.C)
}

func (p *Plane) AttachToBody(body *RigidBody) {
	p.plane.SetBody(body.body)
}

func (p *Plane) OnCollision(other Collider, contact Contact) {
	fmt.Println("Plane Collision!", other)
}
