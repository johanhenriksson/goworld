package gizmo

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/cylinder"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Arrow struct {
	object.Object
	Rigidbody *physics.RigidBody
	Head      *cone.Cone
	Body      *cylinder.Cylinder
	Collider  *physics.Compound

	Hover object.Property[bool]
}

func NewArrow(clr color.T) *Arrow {
	height := float32(1.5)

	coneRadius := height * 0.1
	bodyRadius := coneRadius * 0.25
	coneHeight := 0.33 * height
	bodyHeight := 0.67 * height
	segments := 32

	mat := &material.Def{
		Shader:       "forward/color",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
	}

	arrow := object.New("Arrow", &Arrow{
		Hover:     object.NewProperty(false),
		Rigidbody: physics.NewRigidBody(0),
		Collider:  physics.NewCompound(),

		Head: object.Builder(cone.NewObject(cone.Args{
			Mat:      mat,
			Radius:   coneRadius,
			Height:   coneHeight,
			Segments: segments,
			Color:    clr,
		})).
			Position(vec3.UnitY).
			Create(),

		Body: object.Builder(cylinder.NewObject(cylinder.Args{
			Mat:      mat,
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
