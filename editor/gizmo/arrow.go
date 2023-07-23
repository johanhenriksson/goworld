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
}

func NewArrow(clr color.T) *Arrow {
	height := float32(1.5)

	coneRadius := height * 0.1
	bodyRadius := coneRadius * 0.25
	coneHeight := 0.33 * height
	bodyHeight := 0.67 * height
	segments := 32

	mat := &material.Def{
		Shader:       "color_f",
		Subpass:      "forward",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
	}

	return object.New("Arrow", &Arrow{
		Rigidbody: physics.NewRigidBody(0),

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
			Attach(physics.NewBox(vec3.New(bodyRadius, bodyHeight*0.5, bodyRadius))).
			Create(),

		Collider: physics.NewCompound(),
	})
}
