package editor

import (
	"log"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/gizmo"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type Tool interface {
	Component
	CanDeselect() bool
	ToolMouseEvent(e mouse.Event, hover physics.RaycastHit)
}

const ToolLayer = physics.Mask(2)

type Action struct {
	Name     string
	Icon     icon.Icon
	Key      keys.Code
	Modifier keys.Modifier
	Callback func(*ToolManager)
}

type ToolManager struct {
	Object
	scene    Object
	selected []T
	tool     Tool
	camera   mat4.T
	viewport draw.Viewport

	// built-in tools
	Mover   *gizmo.Mover
	Rotater *gizmo.Rotater
}

func NewToolManager(pool Pool) *ToolManager {
	return NewObject(pool, "Tool Manager", &ToolManager{
		Mover: Builder(gizmo.NewMover(pool)).
			Active(false).
			Create(),

		Rotater: Builder(gizmo.NewRotater(pool)).
			Active(false).
			Create(),

		selected: make([]T, 0, 16),
	})
}

func (m *ToolManager) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := m.viewport.NormalizeCursor(e.Position())

	// calculate a ray going into the screen
	near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
	far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))

	world := GetInParents[*physics.World](m)
	if world == nil {
		return
	}
	hit, _ := world.Raycast(near, far, 1)

	if m.tool != nil {
		// pass on the mouse event
		m.tool.ToolMouseEvent(e, hit)
		if e.Handled() {
			return
		}
	}

	if hit.Shape == nil {
		return
	}

	canReselect := m.selected == nil || m.tool == nil || m.tool.CanDeselect()
	if !canReselect {
		return
	}

	editor := GetInParents[T](hit.Shape)

	// if nothing is selected, or CanDeselect() is true,
	// look for something else to select.
	if e.Button() == mouse.Button1 && e.Action() == mouse.Release {
		if editor != nil {
			// todo: pass physics hit info
			m.setSelect(e, editor)
		} else {
			// deselect
			m.setSelect(e, nil)
		}
	}
}

func (m *ToolManager) Actions() []Action {
	actions := make([]Action, 0, 16)
	actions = append(actions, Action{
		Name:     "Save",
		Icon:     icon.IconSave,
		Key:      keys.S,
		Modifier: keys.Ctrl,
		Callback: func(m *ToolManager) {
			log.Println("save")
		},
	})
	for _, editor := range m.selected {
		for _, action := range editor.Actions() {
			actions = append(actions, action)
		}
	}
	return actions
}

func (m *ToolManager) KeyEvent(e keys.Event) {
	if e.Action() != keys.Release {
		return
	}
	if m.selected != nil {
		if e.Code() == keys.Escape {
			m.setSelect(mouse.NopEvent(), nil)
			e.Consume()
			return
		}

		for _, editor := range m.selected {
			for _, action := range editor.Actions() {
				if action.Key == e.Code() && e.Modifier(action.Modifier) {
					action.Callback(m)
					e.Consume()
				}
			}
		}
	}
}

func (m *ToolManager) Tool() Tool {
	return m.tool
}

func (m *ToolManager) UseTool(tool Tool) {
	// if we select the same tool twice, deselect it instead
	sameTool := m.tool == tool

	// deselect tool
	if m.tool != nil {
		Disable(m.tool)
		m.tool = nil
	}

	// activate the new tool if its different
	if !sameTool && tool != nil {
		m.tool = tool
		Enable(m.tool)
	}
}

func (m *ToolManager) Select(obj T) {
	m.setSelect(mouse.NopEvent(), obj)
}

func (m *ToolManager) Selected() []T {
	return m.selected
}

func (m *ToolManager) setSelect(e mouse.Event, component T) bool {
	// todo: detect if the object has been deleted
	// otherwise CanDeselect() will make it impossible to select another object

	// todo: refactor to enable ALL component editors on the object group
	// group := collider.Parent()

	// editors := Children(group)

	// deselect
	if m.selected != nil {
		// deselect tool
		m.UseTool(nil)

		// deselect object
		for _, editor := range m.selected {
			if !editor.Deselect(e) {
				log.Println("editor", editor.Name(), "attempted to abort deselection")
			}
		}
		m.selected = m.selected[:0]
	}

	// select
	if component != nil {
		group := component
		_, ok := component.Target().(Object)
		if !ok {
			group, ok = component.Parent().(T)
			if !ok {
				return true
			}
		}

		// select game object editor
		group.Select(e)
		m.selected = m.selected[:0]
		m.selected = append(m.selected, group)

		// select child component editors
		for child := range group.Children() {
			if childEdit, ok := child.(T); ok {
				if _, isObject := childEdit.Target().(Object); isObject {
					continue
				}
				childEdit.Select(e)
				m.selected = append(m.selected, childEdit)
			}
		}
	}
	return true
}

func (m *ToolManager) PreDraw(args draw.Args, scene Object) error {
	m.scene = scene
	m.camera = args.Camera.ViewProj
	m.viewport = args.Camera.Viewport
	return nil
}

func (m *ToolManager) MoveTool(obj Component) {
	m.UseTool(m.Mover)
	m.Mover.SetTarget(obj.Transform())
}

func (m *ToolManager) RotateTool(obj Component) {
	m.UseTool(m.Rotater)
	m.Rotater.SetTarget(obj.Transform())
}
