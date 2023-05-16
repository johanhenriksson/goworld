package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	defer func() {
		log.Println("Clean exit")
	}()

	engine.Run(engine.Args{
		Width:  1600,
		Height: 1200,
		Title:  "goworld",
	},
		editor.Scene(func(r renderer.T, scene object.T) {
			world := physics.NewWorld()
			world.Debug(false)
			object.Attach(scene, world)

			generator := chunk.ExampleWorldgen(4, 123123)
			chonk := chunk.NewWorld(4, generator, 40)
			object.Attach(scene, chonk)

			// physics boxes
			boxgen := chunk.NewRandomGen()
			for x := 0; x < 5; x++ {
				for z := 0; z < 5; z++ {
					chonk := chunk.Generate(boxgen, 1, 100*x, 100*z)
					object.Builder(physics.NewRigidBody(5, physics.NewBox(vec3.New(0.5, 0.5, 0.5)))).
						Position(vec3.New(20+3*float32(x), 60, 15+3*float32(z))).
						Attach(object.Builder(chunk.NewMesh(chonk)).
							Position(vec3.New(-0.5, -0.5, -0.5)).
							Create()).
						Parent(world).
						Create()
				}
			}

			m := terrain.NewMap(64, 3)
			tile := m.GetTile(0, 0, true)
			tileMesh := terrain.NewMesh(tile)
			tileMesh.RefreshSync()

			meshShape := physics.NewMesh(tileMesh.Mesh())
			object.Builder(physics.NewRigidBody(0, meshShape)).
				Position(vec3.New(0, 10, 0)).
				Attach(tileMesh).
				Parent(world).
				Create()

			// directional light
			rot := float32(-40)
			object.Attach(
				scene,
				object.Builder(light.NewDirectional(light.DirectionalArgs{
					Intensity: 1.5,
					Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
					Shadows:   true,
					Cascades:  4,
				})).
					Position(vec3.New(1, 2, 3)).
					Attach(script.New(func(scene, self object.T, dt float32) {
						rot -= dt
						self.Parent().Transform().SetRotation(vec3.New(rot, 0, 0))
					})).
					Create())
		}),
	)
}
