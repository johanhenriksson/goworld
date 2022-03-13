package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/vkrender"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/color"
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
			player.Transform().SetPosition(vec3.New(0, 20, -11))
			player.Eye.Transform().SetRotation(vec3.New(-30, 0, 0))

			mesh := game.NewChunkMesh(chunk)
			chunkobj := object.New("chunk", mesh)
			scene.Adopt(chunkobj)

			object.Build("light1").
				Position(vec3.New(10, 9, 13)).
				Attach(light.NewPoint(light.PointArgs{
					Attenuation: light.DefaultAttenuation,
					Color:       color.Red,
					Range:       15,
					Intensity:   15,
				})).
				Parent(scene).
				Create()

			object.Build("light2").
				Position(vec3.New(10-16, 9, 13)).
				Attach(light.NewPoint(light.PointArgs{
					Attenuation: light.DefaultAttenuation,
					Color:       color.Blue,
					Range:       15,
					Intensity:   15,
				})).
				Parent(scene).
				Create()

			scene.Adopt(player)
		},
	})
}
