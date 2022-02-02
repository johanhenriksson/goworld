package deferred

import (
	"github.com/johanhenriksson/goworld/render"
)

type DeferredDrawable interface {
	DrawDeferred(render.Args)
}
