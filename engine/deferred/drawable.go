package deferred

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/material"
)

type Drawable interface {
	object.Component

	DrawDeferred(render.Args) error
	Material() material.T
}

type ShadowDrawable interface {
	object.Component

	DrawShadow(render.Args) error
}
