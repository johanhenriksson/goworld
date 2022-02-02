package deferred

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
)

type DeferredDrawable interface {
	object.Component
	DrawDeferred(render.Args)
}
