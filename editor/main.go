package editor

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
)

func Scene(f object.SceneFunc) object.SceneFunc {
	return func(scene object.Object) {
		// create subscene in a child object
		workspace := object.Empty("Workspace")
		f(workspace)

		editorScene := NewEditorScene(workspace, true)
		object.Attach(scene, editorScene)
	}
}

type EditorScene struct {
	object.Object
	App       *App
	Workspace object.Object

	playing bool
}

func NewEditorScene(workspace object.Object, playing bool) *EditorScene {
	app := NewApp(workspace)
	object.Toggle(app, !playing)

	return object.New("Editor", &EditorScene{
		Object:    object.Scene(),
		App:       app,
		Workspace: workspace,
		playing:   playing,
	})
}

func (s *EditorScene) Replace(workspace object.Object) {
	parent := s.Parent()
	object.Detach(s)
	*s = *NewEditorScene(workspace, s.playing)
	object.Attach(parent, s)
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.O && e.Modifier(keys.Ctrl) {
		c, err := object.Load("scene.scn")
		if err != nil {
			panic(err)
		}
		s.Replace(c.(object.Object))
		log.Println("scene loaded")
	}
	if e.Action() == keys.Release && e.Code() == keys.S && e.Modifier(keys.Ctrl) {
		if err := object.Save("scene.scn", s.Workspace); err != nil {
			panic(err)
		}
		log.Println("scene saved")
	}
	if e.Action() == keys.Release && e.Code() == keys.H {
		object.Toggle(s.App, s.playing)
		s.playing = !s.playing
	} else {
		s.Object.KeyEvent(e)
	}
}

func (s *EditorScene) Update(scene object.Component, dt float32) {
	if s.playing {
		s.Workspace.Update(scene, dt)
	} else {
		s.App.Update(scene, dt)
		s.App.World.DebugDraw()
	}
}
