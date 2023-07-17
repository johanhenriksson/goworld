package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/gizmo"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render"
)

type Selectable interface {
	object.Object

	Actions() []Action
	Select(mouse.Event)
	Deselect(mouse.Event) bool
}

type Tool interface {
	object.Component
	mouse.Handler
	CanDeselect() bool
}

type Action struct {
	Name     string
	Key      keys.Code
	Modifier keys.Modifier
	Callback func(ToolManager)
}

type ToolManager interface {
	object.Component

	Select(Selectable)
	SelectTool(Tool)
	MoveTool(object.Component)
	Tool() Tool
}

type toolmgr struct {
	object.Object
	scene    object.Object
	selected Selectable
	tool     Tool
	camera   mat4.T
	viewport render.Screen

	// built-in tools
	Mover *gizmo.Mover
}

func NewToolManager() ToolManager {
	return object.New("Tool Manager", &toolmgr{
		Mover: object.Builder(gizmo.NewMover()).
			Active(false).
			Create(),
	})
}

func (m *toolmgr) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := m.viewport.NormalizeCursor(e.Position())

	if m.tool != nil {
		// we have a tool selected.
		// pass on the mouse event
		m.tool.MouseEvent(e)
		if e.Handled() {
			return
		}
	}

	canReselect := m.selected == nil || m.tool == nil || m.tool.CanDeselect()
	if !canReselect {
		return
	}

	// if nothing is selected, or CanDeselect() is true,
	// look for something else to select.
	if e.Button() == mouse.Button1 && e.Action() == mouse.Release {
		// calculate a ray going into the screen
		near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
		far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))

		world := object.GetInParents[*physics.World](m)
		if world == nil {
			return
		}

		hit, _ := world.Raycast(near, far)
		if hit.Shape == nil {
			return
		}

		if object := object.GetInParents[Selectable](hit.Shape); object != nil {
			// todo: pass physics hit info instead
			m.setSelect(e, object)
		} else {
			// deselect
			m.setSelect(e, nil)
		}
	}
}

func (m *toolmgr) KeyEvent(e keys.Event) {
	if e.Action() != keys.Release {
		return
	}
	if m.selected != nil {
		if e.Code() == keys.Escape {
			m.setSelect(mouse.NopEvent(), nil)
			e.Consume()
			return
		}

		// todo: consider all editors
		for _, action := range m.selected.Actions() {
			if action.Key == e.Code() && e.Modifier(action.Modifier) {
				action.Callback(m)
				e.Consume()
			}
		}
	}
}

func (m *toolmgr) Tool() Tool {
	return m.tool
}

func (m *toolmgr) SelectTool(tool Tool) {
	// if we select the same tool twice, deselect it instead
	sameTool := m.tool == tool

	// deselect tool
	if m.tool != nil {
		object.Disable(m.tool)
		m.tool = nil
	}

	// activate the new tool if its different
	if !sameTool && tool != nil {
		m.tool = tool
		object.Enable(m.tool)
	}
}

func (m *toolmgr) Select(obj Selectable) {
	m.setSelect(mouse.NopEvent(), obj)
}

func (m *toolmgr) setSelect(e mouse.Event, object Selectable) bool {
	// todo: detect if the object has been deleted
	// otherwise CanDeselect() will make it impossible to select another object

	// todo: refactor to enable ALL component editors on the object group
	// group := collider.Parent()

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
		object.Select(e)
	}
	return true
}

func (m *toolmgr) PreDraw(args render.Args, scene object.Object) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = args.Viewport
	return nil
}

func (m *toolmgr) MoveTool(obj object.Component) {
	m.Mover.SetTarget(obj.Transform())
	m.SelectTool(m.Mover)
}
