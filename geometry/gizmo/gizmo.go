package gizmo

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
)

type Gizmo interface {
	object.T

	Target() transform.T
	SetTarget(transform.T)

	DragStart(mouse.Event, collider.T)
	DragMove(mouse.Event)
	DragEnd(mouse.Event)
}
