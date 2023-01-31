package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func Scene(render renderer.T, scene object.T) {
	// move the existing scene into a child object
	workspace := object.Empty("Workspace")
	for _, existing := range scene.Children() {
		object.Attach(workspace, existing)
	}

	// collision support
	editor := object.Empty("Editor")
	object.Attach(editor, MakeGUI(workspace))
	object.Attach(editor, NewGizmoManager())
	object.Attach(editor, NewSelectManager())

	// first person controls
	player := game.NewPlayer(vec3.New(0, 20, -11), nil)
	player.Eye.Transform().SetRotation(vec3.New(-30, 0, 0))
	object.Attach(editor, player)

	// mover gizmo
	mv := object.Builder(mover.New(mover.Args{})).
		Position(vec3.New(1, 10, 1)).
		Parent(editor).
		Create()
	mv.SetActive(false)

	// construct editors
	context := &Context{
		Camera: player.Camera,
		Render: render,
	}
	object.Attach(editor, ConstructEditors(context, workspace, mv))

	// attach editor & game to scene
	object.Attach(scene, editor)
	object.Attach(scene, workspace)
}
