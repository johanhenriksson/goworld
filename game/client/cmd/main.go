package main

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	_ "github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/client"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	engine.Run(engine.Args{
		Width:  800,
		Height: 600,
		Title:  "goworld",
	},
		editor.Scene(func(scene object.Object) {
			// todo: should on the scene object by default
			object.Attach(scene, physics.NewWorld())

			object.Attach(scene, gui.New(func() node.T {
				return rect.New("plates", rect.Props{})
			}))

			// add water
			object.Attach(scene, object.Builder(terrain.NewWater(512, 2000)).
				Position(vec3.New(0, -1, 0)).
				Create())

			// create game client
			object.Attach(scene, client.NewManager("127.0.0.1"))

			// directional light
			rot := float32(45)
			object.Attach(
				scene,
				object.Builder(object.Empty("Sun")).
					Attach(light.NewDirectional(light.DirectionalArgs{
						Intensity: 1.3,
						Color:     color.RGB(1, 1, 1),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(1, 2, 3)).
					Rotation(quat.Euler(rot, 0, 0)).
					Attach(script.New(func(scene, self object.Component, dt float32) {
						rot -= dt * 360.0 / 86400.0
						self.Parent().Transform().SetRotation(quat.Euler(rot, 0, 0))
					})).
					Create())
		}),
	)
}
