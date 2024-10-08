package gizmo

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

// Gizmos are basically tools that may have mouse-interactive
// 3D components. Perhaps this should merge with the general Tool
// interface?
type Gizmo interface {
	object.Component

	DragStart(e mouse.Event, hit physics.RaycastHit)
	DragEnd(e mouse.Event, hit physics.RaycastHit)
	DragMove(e mouse.Event, hit physics.RaycastHit)
	Hover(bool, physics.Shape)

	Camera() mat4.T
	Viewport() draw.Viewport
	Dragging() bool
}

// Implements mouse dragging behaviour for Gizmos. Move to editor.ToolManager?
func HandleMouse(m Gizmo, e mouse.Event, hit physics.RaycastHit) {
	vpi := m.Camera()
	vpi = vpi.Invert()
	cursor := m.Viewport().NormalizeCursor(e.Position())

	near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
	far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))

	world := object.GetInParents[*physics.World](m)
	hit, ok := world.Raycast(near, far, 2)

	if m.Dragging() {
		if e.Action() == mouse.Release && e.Button() == mouse.Button1 {
			m.DragEnd(e, hit)
			e.Consume()
		} else {
			m.DragMove(e, hit)
			e.Consume()
		}
	} else if e.Action() == mouse.Press && e.Button() == mouse.Button1 {
		if ok {
			m.DragStart(e, hit)
			e.Consume()
		} else if m.Dragging() {
			// we hit nothing, deselect
			m.DragEnd(e, hit)
			e.Consume()
		}
	} else if e.Action() == mouse.Move {
		m.Hover(ok, hit.Shape)
	}
}
