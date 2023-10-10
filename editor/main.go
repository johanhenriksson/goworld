package editor

import (
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
	App *App

	playing bool
}

func NewEditorScene(workspace object.Object) *EditorScene {
	editor := object.New("Editor", &EditorScene{
		Object: object.Scene(),
		App:    NewApp(workspace),
	})
	object.Attach(editor, workspace)
	return editor
}

func (s *EditorScene) Play() {
	s.playing = true
	object.Disable(s.App)
}

func (s *EditorScene) Pause() {
	s.playing = false
	object.Enable(s.App)
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.H {
		if s.playing {
			s.Pause()
		} else {
			s.Play()
		}
	} else {
		s.Object.KeyEvent(e)
	}
}

func (s *EditorScene) Update(scene object.Component, dt float32) {
	if s.playing {
		s.Object.Update(scene, dt)
	} else {
		s.App.Update(scene, dt)
	}
}
