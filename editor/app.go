package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type App struct {
	object.Object
	GUI    gui.Manager
	Tools  *ToolManager
	World  *physics.World
	Player *Player

	objects   object.Pool
	editors   object.Component
	workspace object.Object
}

func NewApp(pool object.Pool, workspace object.Object) *App {
	editor := object.New(pool, "Application", &App{
		World: physics.NewWorld(pool),

		Player:    NewPlayer(pool, vec3.New(-8, 24, -8), quat.Euler(30, 45, 0)),
		objects:   pool,
		editors:   nil,
		workspace: workspace,
	})

	editor.GUI = MakeGUI(pool, editor)
	object.Attach(editor, editor.GUI)

	// must be attached AFTER gui so that input events are handled in the correct order
	editor.Tools = NewToolManager(pool)
	object.Attach(editor, editor.Tools)

	// editor.World.Debug(true)
	return editor
}

func (e *App) Update(scene object.Component, dt float32) {
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
		object.Attach(e, e.editors)
	}
}

func (e *App) Lookup(obj object.Object) T {
	editor, _ := object.NewQuery[T]().Where(func(e T) bool {
		return e.Target() == obj
	}).First(e.editors)
	return editor
}
