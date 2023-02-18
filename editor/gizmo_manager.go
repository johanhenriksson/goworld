package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/gizmo"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type GizmoManager interface {
	object.T
	mouse.Handler
}

type gizmomgr struct {
	object.T
	scene  object.T
	active gizmo.Gizmo

	camera   mat4.T
	viewport render.Screen
}

func NewGizmoManager() GizmoManager {
	return object.New(&gizmomgr{})
}

func (m *gizmomgr) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := m.viewport.NormalizeCursor(e.Position())

	if m.active != nil {
		if e.Action() == mouse.Release {
			m.active.DragEnd(e)
			m.active = nil
			e.Consume()
		} else {
			m.active.DragMove(e)
			e.Consume()
		}
	} else if e.Action() == mouse.Press {
		near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
		far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
		dir := far.Sub(near).Normalized()

		// return closest hit
		// find Collider children of Selectable objects
		gizmos := object.Query[gizmo.Gizmo]().CollectObjects(m.scene)
		colliders := object.Query[collider.T]().Collect(gizmos...)

		closest, hit := collider.ClosestIntersection(colliders, &physics.Ray{
			Origin: near,
			Dir:    dir,
		})

		if hit {
			if giz, ok := object.FindInParents[gizmo.Gizmo](closest); ok {
				m.active = giz
				m.active.DragStart(e, closest)
				e.Consume()
			}
		} else if m.active != nil {
			// we hit nothing, deselect
			m.active.DragEnd(e)
			e.Consume()
		}
	}
}

func (m *gizmomgr) PreDraw(args render.Args, scene object.T) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = args.Viewport
	return nil
}
