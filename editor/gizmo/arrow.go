package gizmo

import (
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/cylinder"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
)

type Arrow struct {
	Object
	Rigidbody *physics.RigidBody
	Head      *cone.Cone
	Body      *cylinder.Cylinder
	Collider  *physics.Compound

	Hover Property[bool]
}

func NewArrow(pool Pool, clr color.T) *Arrow {
	height := float32(1.5)

	coneRadius := height * 0.06
	bodyRadius := coneRadius * 0.1
	coneHeight := 0.2 * height
	bodyHeight := 0.8 * height
	segments := 32

	arrow := NewObject(pool, "Arrow", &Arrow{
		Hover:     NewProperty(false),
		Rigidbody: physics.NewRigidBody(pool, 0),
		Collider:  physics.NewCompound(pool),

		Head: Builder(cone.New(pool, cone.Args{
			Mat:      material.ColoredForward(),
			Radius:   coneRadius,
			Height:   coneHeight,
			Segments: segments,
			Color:    clr,
		})).
			Position(vec3.UnitY).
			Create(),

		Body: Builder(cylinder.New(pool, cylinder.Args{
			Mat:      material.ColoredForward(),
			Radius:   bodyRadius,
			Height:   bodyHeight,
			Segments: segments,
			Color:    clr,
		})).
			Position(vec3.New(0, bodyHeight*0.5, 0)).
			Create(),
	})

	arrow.Rigidbody.Layer.Set(2)

	arrow.Hover.OnChange.Subscribe(func(b bool) {
		scale := vec3.One
		if b {
			scale = vec3.New(1.2, 1.2, 1.2)
		}
		arrow.Head.Transform().SetScale(scale)
	})

	return arrow
}
