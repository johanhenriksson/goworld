package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type SelectManager interface {
	object.T
	mouse.Handler
}

type selectmgr struct {
	object.T
	scene    object.T
	selected Selectable

	camera   mat4.T
	viewport render.Screen
}

func NewSelectManager() SelectManager {
	return object.New(&selectmgr{})
}

func (m *selectmgr) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := m.viewport.NormalizeCursor(e.Position())

	if e.Button() == mouse.Button1 && e.Action() == mouse.Release {
		// calculate a ray going into the screen
		near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
		far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
		dir := far.Sub(near).Normalized()

		// find Collider children of Selectable objects
		selectables := object.Query[Selectable]().CollectObjects(m.scene)
		colliders := object.Query[collider.T]().Collect(selectables...)

		// return closest hit
		closest, hit := collider.ClosestIntersection(colliders, &physics.Ray{
			Origin: near,
			Dir:    dir,
		})

		if hit {
			if selectable, ok := object.FindInParents[Selectable](closest); ok {
				m.Select(e, selectable, closest)
			}
		} else if m.selected != nil {
			// we hit nothing, deselect
			m.Deselect(e)
		}
	}
}

func (m *selectmgr) Select(e mouse.Event, object Selectable, collider collider.T) {
	m.Deselect(e)
	if m.selected == nil {
		m.selected = object
		object.Select(e, collider)
	}
}

func (m *selectmgr) Deselect(e mouse.Event) {
	if m.selected != nil {
		ok := m.selected.Deselect(e)
		if ok {
			m.selected = nil
		}
	}
}

func (m *selectmgr) PreDraw(args render.Args, scene object.T) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = args.Viewport
	return nil
}
