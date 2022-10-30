package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/game/editor"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

func main() {
	defer func() {
		log.Println("Clean exit")
	}()

	backend := vulkan.New("goworld: vulkan", 0)

	engine.Run(engine.Args{
		Backend: backend,
		Width:   1600,
		Height:  1200,
		Title:   "goworld: vulkan",
		Renderer: func() renderer.T {
			// todo: this causes a resource leak - destroy it somehow
			voxelCache := cache.NewMeshCache(backend)
			return renderer.New(
				backend,
				[]pass.DeferredSubpass{game.NewVoxelSubpass(backend, voxelCache)},
				[]pass.DeferredSubpass{game.NewVoxelShadowpass(backend, voxelCache)},
			)
		},
	},
		makeGui,
		func(r renderer.T, scene object.T) {
			game.CreateScene(scene, r.Buffers())

			// mesh := game.NewChunkMesh(chunk)
			// chunkobj := object.New("chunk", mesh)
			// scene.Adopt(chunkobj)

			// create editor
			// edit := editor.NewEditor(chunk, player.Camera, r.Buffers())
			// scene.Adopt(edit.Object())

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
		},
	)
}

func makeGui(r renderer.T, scene object.T) {
	scene.Attach(gui.New(func() node.T {
		return rect.New("sidebar", rect.Props{
			OnMouseDown: func(e mouse.Event) {},
			Style: rect.Style{
				Layout: style.Column{},
				Width:  style.Pct(15),
				Height: style.Pct(100),
				Color:  color.Black.WithAlpha(0.5),
			},
			Children: []node.T{
				palette.New("palette", palette.Props{
					Palette: color.DefaultPalette,
					OnPick: func(clr color.T) {
						editor := query.New[editor.T]().First(scene)
						if editor == nil {
							panic("could not find editor")
						}

						editor.SelectColor(clr)
					},
				}),
			},
		})
	}))
}