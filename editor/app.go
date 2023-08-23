package editor

import (
	"log"
	"os"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/gizmo"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type App struct {
	object.Object
	GUI    gui.Manager
	World  *physics.World
	Player *Player

	// Tools
	Mover   *gizmo.Mover
	Rotater *gizmo.Rotater

	selected  []T
	tool      Tool
	editors   object.Component
	workspace object.Object
}

func NewApp(workspace object.Object) *App {
	editor := object.New("Application", &App{
		World: physics.NewWorld(),

		Player: NewPlayer(vec3.New(-8, 24, -8), quat.Euler(30, 45, 0)),

		Mover: object.Builder(gizmo.NewMover()).
			Active(false).
			Create(),

		Rotater: object.Builder(gizmo.NewRotater()).
			Active(false).
			Create(),

		editors:   nil,
		workspace: workspace,
		selected:  make([]T, 0, 16),
	})

	editor.GUI = MakeGUI(editor)
	object.Attach(editor, editor.GUI)

	// editor.World.Debug(true)
	return editor
}

func (e *App) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.Refresh()
}

func (e *App) Open(scene object.Object) {
	// replace workspace
	object.Detach(e.workspace)
	object.Attach(e.Parent(), scene)
	e.workspace = scene

	// recreate editors
	object.Detach(e.editors)
	e.editors = nil
}

func (e *App) Load(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	scene, err := object.Load(fp)
	if err != nil {
		return err
	}

	e.Open(scene.(object.Object))
	return nil
}

func (e *App) Save(path string) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if err := object.Save(fp, e.workspace); err != nil {
		return err
	}
	return nil
}

func (e *App) Refresh() {
	context := &Context{
		Camera: e.Player.Camera.Camera,
		Scene:  e.workspace,
	}
	e.editors = ConstructEditors(context, e.editors, e.workspace)
	if e.editors.Parent() == nil {
		object.Attach(e, e.editors)
	}
}

func (e *App) Lookup(obj object.Object) T {
	editor, _ := object.NewQuery[T]().Where(func(e T) bool {
		return e.Target() == obj
	}).First(e.editors)
	return editor
}

func (m *App) Actions() []Action {
	actions := make([]Action, 0, 16)
	actions = append(actions, Action{
		Name:     "Open",
		Icon:     icon.IconFileOpen,
		Key:      keys.O,
		Modifier: keys.Ctrl,
		Callback: func(m *App) {
			if err := m.Load("scene.goworld"); err != nil {
				panic("failed to load: " + err.Error())
			}
			log.Println("loaded")
		},
	})
	actions = append(actions, Action{
		Name:     "Save",
		Icon:     icon.IconSave,
		Key:      keys.S,
		Modifier: keys.Ctrl,
		Callback: func(m *App) {
			if err := m.Save("scene.goworld"); err != nil {
				panic("failed to save: " + err.Error())
			}
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

func (m *App) Tool() Tool {
	return m.tool
}

func (m *App) UseTool(tool Tool) {
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

func (m *App) Select(obj T) {
	m.setSelect(mouse.NopEvent(), obj)
}

func (m *App) Selected() []T {
	return m.selected
}

func (m *App) setSelect(e mouse.Event, component T) bool {
	// todo: detect if the object has been deleted
	// otherwise CanDeselect() will make it impossible to select another object

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
		_, ok := component.Target().(object.Object)
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
		for _, child := range group.Children() {
			if childEdit, ok := child.(T); ok {
				if _, isObject := childEdit.Target().(object.Object); isObject {
					continue
				}
				childEdit.Select(e)
				m.selected = append(m.selected, childEdit)
			}
		}
	}
	return true
}

func (m *App) MoveTool(obj object.Component) {
	m.UseTool(m.Mover)
	m.Mover.SetTarget(obj.Transform())
}

func (m *App) RotateTool(obj object.Component) {
	m.UseTool(m.Rotater)
	m.Rotater.SetTarget(obj.Transform())
}

func (m *App) MouseEvent(e mouse.Event) {
	m.toolMouseEvent(e)
	m.Object.MouseEvent(e)
}

func (m *App) toolMouseEvent(e mouse.Event) {
	vpi := m.Player.Camera.ViewProjInv
	cursor := m.Player.Camera.Viewport.NormalizeCursor(e.Position())

	// calculate a ray going into the screen
	near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
	far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))

	hit, _ := m.World.Raycast(near, far, 1)

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

	editor := object.GetInParents[T](hit.Shape)

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

func (m *App) KeyEvent(e keys.Event) {
	if m.selected != nil {
		if e.Code() == keys.Escape {
			m.setSelect(mouse.NopEvent(), nil)
			e.Consume()
			return
		}

		for _, action := range m.Actions() {
			if action.Key == e.Code() && e.Modifier(action.Modifier) {
				if e.Action() == keys.Release {
					action.Callback(m)
				}
				e.Consume()
			}
		}
	}

	if !e.Handled() {
		m.Object.KeyEvent(e)
	}
}
