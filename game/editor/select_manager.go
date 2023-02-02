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

	Select(Selectable)
	Validate(Selectable) bool
}

type SelectCallback func(Selectable)
type SelectFilter func(Selectable) bool

type selectmgr struct {
	object.T
	scene    object.T
	selected Selectable

	camera   mat4.T
	viewport render.Screen

	filter   SelectFilter
	onSelect SelectCallback
}

func NewSelectManager(onSelect SelectCallback, filter SelectFilter) SelectManager {
	f := func(Selectable) bool { return true }
	if filter != nil {
		f = filter
	}

	return object.New(&selectmgr{
		onSelect: onSelect,
		filter:   f,
	})
}

func (m *selectmgr) Validate(obj Selectable) bool {
	return m.filter(obj)
}

func (m *selectmgr) Select(obj Selectable) {
	changed := m.setSelect(mouse.NopEvent(), obj, nil)
	if changed && m.onSelect != nil {
		m.onSelect(obj)
	}
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
		selectables := object.Query[Selectable]().Where(m.filter).CollectObjects(m.scene)
		colliders := object.Query[collider.T]().Collect(selectables...)

		// return closest hit
		closest, hit := collider.ClosestIntersection(colliders, &physics.Ray{
			Origin: near,
			Dir:    dir,
		})

		if hit {
			if selectable, ok := object.FindInParents[Selectable](closest); ok {
				m.setSelect(e, selectable, closest)
			}
		} else if m.selected != nil {
			// we hit nothing, deselect
			m.setSelect(e, nil, nil)
		}
	}
}

func (m *selectmgr) setSelect(e mouse.Event, object Selectable, collider collider.T) bool {
	// deselect
	if m.selected != nil {
		if !m.selected.Deselect(e) {
			return false
		}
		m.selected = nil
	}
	// select
	if object != nil {
		m.selected = object
		object.Select(e, collider)
	}
	return true
}

func (m *selectmgr) PreDraw(args render.Args, scene object.T) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = args.Viewport
	return nil
}
