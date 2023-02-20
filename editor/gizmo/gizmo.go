package gizmo

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Gizmos are basically tools that may have mouse-interactive
// 3D components. Perhaps this should merge with the general Tool
// interface?
type Gizmo interface {
	object.T

	DragStart(e mouse.Event, collider collider.T)
	DragEnd(e mouse.Event)
	DragMove(e mouse.Event)

	Camera() mat4.T
	Viewport() render.Screen
	Dragging() bool
}

// Implements mouse dragging behaviour for Gizmos. Move to editor.ToolManager?
func HandleMouse(m Gizmo, e mouse.Event) {
	vpi := m.Camera()
	vpi = vpi.Invert()
	cursor := m.Viewport().NormalizeCursor(e.Position())

	if m.Dragging() {
		if e.Action() == mouse.Release {
			m.DragEnd(e)
			e.Consume()
		} else {
			m.DragMove(e)
			e.Consume()
		}
	} else if e.Action() == mouse.Press {
		near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
		far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
		dir := far.Sub(near).Normalized()

		// return closest hit
		// find Collider children of Selectable objects
		colliders := object.Query[collider.T]().Collect(m)

		closest, hit := collider.ClosestIntersection(colliders, &physics.Ray{
			Origin: near,
			Dir:    dir,
		})

		if hit {
			m.DragStart(e, closest)
			e.Consume()
		} else if m.Dragging() {
			// we hit nothing, deselect
			m.DragEnd(e)
			e.Consume()
		}
	}
}
