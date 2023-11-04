package client

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"

	// "github.com/johanhenriksson/goworld/game/server"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Entity interface {
	object.Object
	EntityID() server.Identity
}

type entity struct {
	object.Object
	Model object.Object

	id server.Identity
}

func NewEntity(id server.Identity, pos vec3.T) *entity {
	// entity model
	model := cube.New(cube.Args{
		Size: 1,
		Mat:  material.StandardDeferred(),
	})
	colorIdx := uint64(id) % uint64(len(color.DefaultPalette))
	model.SetTexture(texture.Diffuse, color.DefaultPalette[colorIdx])

	return object.New(fmt.Sprintf("Entity %x", id), &entity{
		id: id,
		Model: object.Builder(object.Empty("Model")).
			Attach(model).
			Scale(vec3.New(1, 2, 1)).
			Create(),
	})
}

func (e *entity) EntityID() server.Identity {
	return e.id
}
