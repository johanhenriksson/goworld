package builtin

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.World{}, NewWorldEditor)
}

type WorldEditor struct {
	*editor.ComponentEditor
	target *physics.World
}

func NewWorldEditor(ctx *editor.Context, world *physics.World) *WorldEditor {
	world.Debug(false)
	return object.NewComponent(ctx.Objects, &WorldEditor{
		ComponentEditor: editor.NewComponentEditor(ctx.Objects, world),
		target:          world,
	})
}

func (e *WorldEditor) Update(scene object.Component, dt float32) {
	e.target.DebugDraw()
}
