package builtin

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.World{}, NewWorldEditor)
}

type WorldEditor struct {
	object.T
	target *physics.World
}

func NewWorldEditor(ctx *editor.Context, world *physics.World) *WorldEditor {
	log.Println("create physics world editor")
	return object.New(&WorldEditor{
		target: world,
	})
}

func (e *WorldEditor) Actions() []editor.Action {
	return nil
}

func (e *WorldEditor) Update(scene object.T, dt float32) {
	e.target.Update(scene, dt)
	e.target.DebugDraw()
}
