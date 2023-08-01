package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type Editor struct {
	object.Object
	GUI    gui.Manager
	Tools  ToolManager
	World  *physics.World
	Player *Player

	editors   object.Component
	workspace object.Object
}

func NewEditor(workspace object.Object) *Editor {
	editor := object.New("Editor", &Editor{
		GUI:   MakeGUI(),
		Tools: NewToolManager(),
		World: physics.NewWorld(),

		Player:    NewPlayer(vec3.New(0, 25, -11), quat.Euler(-10, 30, 0)),
		editors:   nil,
		workspace: workspace,
	})

	object.Attach(editor, SidebarFragment(
		gui.FragmentFirst,
		func() node.T {
			return ObjectList("scene-graph", ObjectListProps{
				Scene:       workspace,
				EditorRoot:  editor,
				ToolManager: editor.Tools,
			})
		},
	))

	// editor.World.Debug(true)
	return editor
}

func (e *Editor) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)

	context := &Context{
		Camera: e.Player.Camera.Camera,
		Root:   scene,
		Scene:  e.workspace,
	}
	e.editors = ConstructEditors(context, e.editors, e.workspace)
	if e.editors.Parent() == nil {
		object.Attach(e, e.editors)
	}
}

func Scene(f engine.SceneFunc) engine.SceneFunc {
	return func(scene object.Object) {
		// create subscene in a child object
		workspace := object.Empty("Workspace")
		f(workspace)

		editorScene := NewEditorScene(workspace)
		object.Attach(scene, editorScene)
	}
}

type EditorScene struct {
	object.Object
	Editor    *Editor
	Workspace object.Object

	playing bool
}

func NewEditorScene(workspace object.Object) *EditorScene {
	return object.New("EditorScene", &EditorScene{
		Object:    object.Scene(),
		Editor:    NewEditor(workspace),
		Workspace: workspace,
	})
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.H {
		object.Toggle(s.Editor, s.playing)
		s.playing = !s.playing
	} else {
		s.Object.KeyEvent(e)
	}
}

func (s *EditorScene) Update(scene object.Component, dt float32) {
	if s.playing {
		s.Workspace.Update(scene, dt)
	} else {
		s.Editor.Update(scene, dt)
		s.Editor.World.DebugDraw()
	}
}
