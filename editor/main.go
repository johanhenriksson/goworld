package editor

import (
	"log"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/keys"
	. "github.com/johanhenriksson/goworld/core/object"
)

func WrapScene(f SceneFunc) SceneFunc {
	return func(ctx Pool, scene Object) {
		// create subscene in a child object
		workspace := Empty(ctx, "Workspace")
		f(ctx, workspace)

		editorScene := NewEditorScene(ctx, workspace, true)
		Attach(scene, editorScene)
	}
}

type EditorScene struct {
	Object
	App       *App
	Workspace Object
	Objects   Pool

	playing bool
}

func NewEditorScene(pool Pool, workspace Object, playing bool) *EditorScene {
	app := NewApp(pool, workspace)
	Toggle(app, !playing)

	return NewObject(pool, "Editor", &EditorScene{
		Object:    Scene(pool),
		Objects:   pool,
		App:       app,
		Workspace: workspace,
		playing:   playing,
	})
}

func (s *EditorScene) Replace(workspace Object) {
	parent := s.Parent()
	Detach(s)
	*s = *NewEditorScene(s.Objects, workspace, s.playing)
	Attach(parent, s)
}

func (s *EditorScene) KeyEvent(e keys.Event) {
	if e.Action() == keys.Release && e.Code() == keys.O && e.Modifier(keys.Ctrl) {
		c, err := Load[Object](s.Objects, assets.FS, "scene.scn")
		if err != nil {
			panic(err)
		}
		s.Replace(c)
		log.Println("scene loaded")
	}
	if e.Action() == keys.Release && e.Code() == keys.S && e.Modifier(keys.Ctrl) {
		if err := Save(assets.FS, "scene.scn", s.Workspace); err != nil {
			panic(err)
		}
		log.Println("scene saved")
	}
	if e.Action() == keys.Release && e.Code() == keys.H {
		Toggle(s.App, s.playing)
		s.playing = !s.playing
	} else {
		s.Object.KeyEvent(e)
	}
}

func (s *EditorScene) Update(scene Component, dt float32) {
	if s.playing {
		s.Workspace.Update(scene, dt)
	} else {
		s.App.Update(scene, dt)
		s.App.World.DebugDraw()
	}
}
