package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/editor"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

func NewVoxelRenderer(target vulkan.Target) renderer.T {
	return renderer.NewGraph(
		target,
		[]pass.DeferredSubpass{
			game.NewVoxelSubpass(target),
		},
		[]pass.DeferredSubpass{
			game.NewVoxelShadowpass(target),
		},
	)
}

func main() {
	defer func() {
		log.Println("Clean exit")
	}()

	engine.Run(engine.Args{
		Width:    1600,
		Height:   1200,
		Title:    "goworld: vulkan",
		Renderer: NewVoxelRenderer,
	},
		func(r renderer.T, scene object.T) {
			// make some chonks
			world := chunk.NewWorld(31481284, 8)
			chonk := world.AddChunk(0, 0)
			chonk2 := world.AddChunk(1, 0)

			object.Builder(object.Empty("Chunk")).
				Attach(chunk.NewMesh(chonk, nil)).
				Parent(scene).
				Create()
			object.Builder(object.Empty("Chunk2")).
				Attach(chunk.NewMesh(chonk2, nil)).
				Position(vec3.New(8, 0, 0)).
				Parent(scene).
				Create()

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
