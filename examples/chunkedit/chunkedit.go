package main

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	engine.Run(engine.Args{
		Width:  1600,
		Height: 1200,
		Title:  "goworld: chunk edit",
	},
		editor.Scene(func(scene object.Object) {
			chk := chunk.New("editable", 16)
			chunk := chunk.NewMesh(chk)
			object.Attach(scene, chunk)

			object.Builder(object.Empty("Chunk")).
				Attach(physics.NewRigidBody(0)).
				Attach(physics.NewMesh()).
				Attach(chunk).
				Parent(scene).
				Create()

			// directional light
			rot := float32(40)
			object.Attach(
				scene,
				object.Builder(object.Empty("Sun")).
					Attach(light.NewDirectional(light.DirectionalArgs{
						Intensity: 1.3,
						Color:     color.RGB(1, 1, 1),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(16, 16, -4)).
					Rotation(quat.Euler(rot, 0, 0)).
					Attach(script.New(func(scene, self object.Component, dt float32) {
						rot -= 0.5 * dt
						self.Parent().Transform().SetRotation(quat.Euler(rot, 0, 0))
					})).
					Create())
		}),
	)
}
