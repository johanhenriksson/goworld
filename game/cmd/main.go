package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/math/vec3"
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
			generator := chunk.ExampleWorldgen(3141389, 32)
			object.Attach(scene, chunk.NewWorld(32, generator, 500))
			// chonk := chunk.Generate(generator, 32, 0, 0)
			// object.Attach(scene, chunk.NewMesh(chonk))
			// object.Attach(scene, object.Builder(chunk.NewMesh(chonk)).Position(vec3.New(32, 0, 0)).Create())

			// directional light
			object.Attach(
				scene,
				object.Builder(light.NewDirectional(light.DirectionalArgs{
					Intensity: 1.5,
					Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
					Shadows:   true,
					Cascades:  4,
				})).
					Rotation(vec3.New(-30, 0, 0)).
					Position(vec3.New(1, 2, 3)).
					Create())
		}),
	)
}
