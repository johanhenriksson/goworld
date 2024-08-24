package editor

import (
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type App struct {
	Object
	GUI    gui.Manager
	Tools  *ToolManager
	World  *physics.World
	Player *Player

	objects   Pool
	editors   Component
	workspace Object
}

func NewApp(pool Pool, workspace Object) *App {
	editor := NewObject(pool, "Application", &App{
		World: physics.NewWorld(pool),

		Player:    NewPlayer(pool, vec3.New(-8, 24, -8), quat.Euler(30, 45, 0)),
		objects:   pool,
		editors:   nil,
		workspace: workspace,
	})

	editor.GUI = MakeGUI(pool, editor)
	Attach(editor, editor.GUI)

	// must be attached AFTER gui so that input events are handled in the correct order
	editor.Tools = NewToolManager(pool)
	Attach(editor, editor.Tools)

	// editor.World.Debug(true)
	return editor
}

func (e *App) Update(scene Component, dt float32) {
	e.Object.Update(scene, dt)
	e.Refresh()
}

func (e *App) Refresh() {
	context := &Context{
		Objects: e.objects,
		Camera:  e.Player.Camera.Camera,
		Scene:   e.workspace,
	}
	e.editors = ConstructEditors(context, e.editors, e.workspace)
	if e.editors.Parent() == nil {
		Attach(e, e.editors)
	}
}

func (e *App) Lookup(obj Object) T {
	editor, _ := NewQuery[T]().Where(func(e T) bool {
		return e.Target() == obj
	}).First(e.editors)
	return editor
}
