package main

import (
	"log"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/game/editor"
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
		editor.Scene,
	)
}
