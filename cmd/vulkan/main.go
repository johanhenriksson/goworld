package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/vkrender"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
)

func main() {
	defer func() {
		log.Println("Clean exit")
	}()

	backend := vulkan.New("goworld: vulkan", 0)
	defer backend.Destroy()

	engine.Run(engine.Args{
		Backend: backend,
		Width:   1600,
		Height:  1200,
		Title:   "goworld: vulkan",
		Renderer: func() engine.Renderer {
			return vkrender.NewRenderer(backend)
		},
		SceneFunc: func(r engine.Renderer, scene object.T) {
			player, chunk := game.CreateScene(r, scene)

			mesh := game.NewChunkMesh(chunk)
			chunkobj := object.New("chunk", mesh)
			scene.Adopt(chunkobj)

			scene.Adopt(player)
		},
	})
}