package editor

import (
	"log"
	"os"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
)

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
	App       *App
	Workspace object.Object

	playing bool
}

func NewEditorScene(workspace object.Object) *EditorScene {
	return object.New("Editor", &EditorScene{
		Object:    object.Scene(),
		App:       NewApp(workspace),
		Workspace: workspace,
	})
}

func (s *EditorScene) Replace(workspace object.Object) {
	parent := s.Parent()
	object.Detach(s)
	*s = *NewEditorScene(workspace)
	object.Attach(parent, s)
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.O && e.Modifier(keys.Ctrl) {
		fp, err := os.Open("scene.scn")
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		c, err := object.Load(fp)
		if err != nil {
			panic(err)
		}
		s.Replace(c.(object.Object))
		log.Println("scene loaded")
	}
	if e.Action() == keys.Release && e.Code() == keys.S && e.Modifier(keys.Ctrl) {
		fp, err := os.OpenFile("scene.scn", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		if err := object.Save(fp, s.Workspace); err != nil {
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
