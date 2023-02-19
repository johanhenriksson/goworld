package editor

import (
	"log"

	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Selectable interface {
	object.T
	Select(mouse.Event, collider.T)
	Deselect(mouse.Event) bool
	Actions() []Action
}

type Tool interface {
	mouse.Handler
	CanDeselect() bool
	SetActive(bool)
}

type Action struct {
	Key      keys.Code
	Modifier keys.Modifier
	Name     string
	Tool     Tool
}

type SelectManager interface {
	object.T
	mouse.Handler

	Select(Selectable)
	SelectTool(Tool)
}

type SelectCallback func(Selectable)
type SelectFilter func(Selectable) bool

type selectmgr struct {
	object.T
	scene    object.T
	selected Selectable
	tool     Tool

	camera   mat4.T
	viewport render.Screen
}

func NewSelectManager() SelectManager {
	return object.New(&selectmgr{})
}

func (m *selectmgr) Select(obj Selectable) {
	m.setSelect(mouse.NopEvent(), obj, nil)
}

func (m *selectmgr) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := m.viewport.NormalizeCursor(e.Position())

	if m.tool != nil {
		// we have a tool selected.
		// pass on the mouse event
		m.tool.MouseEvent(e)
	}

	canReselect := m.selected == nil || m.tool == nil || m.tool.CanDeselect()
	if !canReselect {
		return
	}

	// if nothing is selected, or CanDeselect() is true,
	// look for something else to select.
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
				m.setSelect(e, selectable, closest)
			}
		} else if m.selected != nil {
			// we hit nothing, deselect
			m.setSelect(e, nil, nil)
		}
	}
}

func (m *selectmgr) KeyEvent(e keys.Event) {
	if e.Action() != keys.Release {
		return
	}
	if m.selected != nil {
		if e.Code() == keys.Escape {
			m.setSelect(mouse.NopEvent(), nil, nil)
			e.Consume()
			return
		}
		for _, action := range m.selected.Actions() {
			if action.Key == e.Code() && e.Modifier(action.Modifier) {
				// woot!
				if m.tool != action.Tool {
					log.Println("select tool", action.Name)
					m.SelectTool(action.Tool)
				} else {
					// deselect
					log.Println("deselect tool", action.Name)
					m.SelectTool(nil)
				}
				e.Consume()
			}
		}
	}
}

func (m *selectmgr) SelectTool(tool Tool) {
	// deselect tool
	if m.tool != nil {
		m.tool.SetActive(false)
		m.tool = nil
	}
	if tool != nil {
		m.tool = tool
		m.tool.SetActive(true)
	}
}

func (m *selectmgr) setSelect(e mouse.Event, object Selectable, collider collider.T) bool {
	// todo: detect if the object has been deleted
	// otherwise CanDeselect() will make it impossible to select another object

	// deselect
	if m.selected != nil {
		// deselect tool
		m.SelectTool(nil)

		// deselect object
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
