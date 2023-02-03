package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/editor"
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
		Title:  "goworld: vulkan",
	},
		func(r renderer.T, scene object.T) {
			// make some chonks
			generator := chunk.ExampleWorldgen(3141389, 32)
			object.Attach(scene, chunk.NewWorld(32, generator, 40))

			// directional light
			object.Attach(scene, light.NewDirectional(light.DirectionalArgs{
				Intensity: 1.6,
				Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
				Direction: vec3.New(0.95, -1.9, 1.05),
				Shadows:   true,
			}))
		},
		editor.Scene,
	)
}
