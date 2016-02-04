package physics

import (
    "github.com/ianremmler/ode"
)

type RigidBox struct {
    *RigidBody
    collider ode.Box

    /* Size */
    X float32
    Y float32
    Z float32
}

type RigidBody struct {
    Mass    float32

    mass    *ode.Mass
}

func NewRigidBody(mass float32) *RigidBody {


    m := ode.NewMass()
    m.Adjust(float64(mass))

    rb := &RigidBody{
        Mass: mass,
        mass: m,
    }
    return rb
}

func (rb *RigidBody) setBox(density, x, y, z float32) {
    rb.mass.SetBox(float64(density), ode.V3(float64(x), float64(y), float64(z)))
    rb.mass.Adjust(float64(rb.Mass))
}

