package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Editor struct {
	object.G
	GUI    gui.Manager
	Tools  ToolManager
	Player *Player

	editors   object.T
	workspace object.G
	render    renderer.T
}

func NewEditor(render renderer.T, workspace object.G) *Editor {
	editor := object.Group("Editor", &Editor{
		GUI:   MakeGUI(render),
		Tools: NewToolManager(),

		Player:    NewPlayer(vec3.New(0, 25, -11), quat.Euler(-10, 30, 0)),
		editors:   nil,
		workspace: workspace,
		render:    render,
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

	return editor
}

func (e *Editor) Update(scene object.T, dt float32) {
	e.G.Update(scene, dt)

	context := &Context{
		Camera: e.Player.Camera.T,
		Render: e.render,
		Root:   scene,
		Scene:  e.workspace,
	}
	e.editors = ConstructEditors(context, e.editors, e.workspace)
	if e.editors.Parent() == nil {
		object.Attach(e, e.editors)
	}
}

func Scene(f engine.SceneFunc) engine.SceneFunc {
	return func(render renderer.T, scene object.G) {
		// create subscene in a child object
		workspace := object.Empty("Workspace")
		f(render, workspace)

		editorScene := NewEditorScene(render, workspace)
		object.Attach(scene, editorScene)
	}
}

type EditorScene struct {
	object.G
	Editor    object.G
	Workspace object.G

	playing bool
}

func NewEditorScene(render renderer.T, workspace object.G) *EditorScene {
	return object.Group("EditorScene", &EditorScene{
		Editor:    NewEditor(render, workspace),
		Workspace: workspace,
	})
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.H {
		s.playing = !s.playing
		s.Editor.SetActive(!s.playing)
	} else {
		s.G.KeyEvent(e)
	}
}

func (s *EditorScene) Update(scene object.T, dt float32) {
	s.Editor.Update(scene, dt)
	if s.playing {
		s.Workspace.Update(scene, dt)
	}
}
