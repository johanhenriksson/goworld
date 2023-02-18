package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/gizmo/mover"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Editor struct {
	object.T
	GUI       gui.Manager
	GizmoMgr  GizmoManager
	SelectMgr SelectManager
	Player    *game.Player

	editors   object.T
	workspace object.T
	render    renderer.T
}

func NewEditor(render renderer.T, workspace object.T) *Editor {
	editor := object.New(&Editor{
		GUI:       MakeGUI(render),
		GizmoMgr:  NewGizmoManager(),
		SelectMgr: NewSelectManager(nil, nil),

		Player:    game.NewPlayer(vec3.New(0, 20, -11), nil),
		editors:   nil,
		workspace: workspace,
		render:    render,
	})

	object.Attach(editor, SidebarFragment(
		gui.FragmentFirst,
		func() node.T {
			return ObjectList("scene-graph", ObjectListProps{
				Scene:         workspace,
				EditorRoot:    editor,
				SelectManager: editor.SelectMgr,
			})
		},
	))

	editor.Player.Eye.Transform().SetRotation(vec3.New(-30, 0, 0))

	return editor
}

func (e *Editor) Update(scene object.T, dt float32) {
	e.T.Update(scene, dt)

	context := &Context{
		Camera: e.Player.Camera,
		Render: e.render,
		Root:   scene,
	}
	e.editors = ConstructEditors(context, e.editors, e.workspace)
	if e.editors.Parent() == nil {
		object.Attach(e, e.editors)
	}
}

func Scene(f engine.SceneFunc) engine.SceneFunc {
	return func(render renderer.T, scene object.T) {
		// create subscene in a child object
		workspace := object.Empty("Workspace")
		f(render, workspace)

		editor := NewEditor(render, workspace)

		// mover gizmo
		mv := object.Builder(mover.New(mover.Args{})).
			Position(vec3.New(1, 40, 1)).
			Parent(editor).
			Create()
		// mv.SetActive(false)
		object.Attach(editor, mv)

		// attach editor & game to scene
		object.Attach(scene, editor)
		object.Attach(scene, workspace)
	}
}
