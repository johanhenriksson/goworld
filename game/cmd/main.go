package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/player"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	defer func() {
		log.Println("exited main defer")
	}()

	engine.Run(engine.Args{
		Width:  1600,
		Height: 1200,
		Title:  "goworld",
	},
		editor.Scene(func(r renderer.T, scene object.Object) {
			world := physics.NewWorld()
			object.Attach(scene, world)

			// generator := chunk.ExampleWorldgen(4, 123123)
			// chonk := chunk.NewWorld(4, generator, 40)
			// object.Attach(scene, chonk)

			// physics boxes
			boxgen := chunk.NewRandomGen()
			for x := 0; x < 3; x++ {
				for z := 0; z < 3; z++ {
					chonk := chunk.Generate(boxgen, 1, 100*x, 100*z)
					object.Builder(physics.NewRigidBody("Box", 5)).
						Position(vec3.New(20+3*float32(x), 30, 15+3*float32(z))).
						Attach(object.Builder(object.Empty("ChunkMesh")).
							Attach(chunk.NewMesh(chonk)).
							Attach(physics.NewMesh()).
							Create()).
						Parent(scene).
						Create()
				}
			}

			// character
			char := player.New()
			char.Transform().SetPosition(vec3.New(5, 16, 5))
			object.Attach(scene, char)

			m := terrain.NewMap(64, 3)
			tile := m.GetTile(0, 0, true)
			tileMesh := terrain.NewMesh(tile)

			meshShape := physics.NewMesh()
			object.Builder(physics.NewRigidBody("Terrain", 0)).
				Position(vec3.New(0, 10, 0)).
				Attach(tileMesh).
				Attach(meshShape).
				Parent(scene).
				Create()

			// directional light
			rot := float32(-40)
			object.Attach(
				scene,
				object.Builder(object.Empty("Sun")).
					Attach(light.NewDirectional(light.DirectionalArgs{
						Intensity: 1.5,
						Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(1, 2, 3)).
					Attach(script.New(func(scene, self object.Component, dt float32) {
						rot -= dt
						self.Parent().Transform().SetRotation(quat.Euler(rot, 0, 0))
					})).
					Create())
		}),
	)
}
