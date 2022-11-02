package main

import (
	"fmt"
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
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type voxrender struct {
	renderer.T
	voxelCache cache.MeshCache
}

func NewVoxelRenderer(backend vulkan.T) renderer.T {
	voxelCache := cache.NewSharedMeshCache(backend, 16_777_216)
	return &voxrender{
		voxelCache: voxelCache,
		T: renderer.New(
			backend,
			[]pass.DeferredSubpass{
				game.NewVoxelSubpass(backend, voxelCache),
			},
			[]pass.DeferredSubpass{
				game.NewVoxelShadowpass(backend, voxelCache),
			},
		),
	}
}

func (r *voxrender) Draw(args render.Args, scene object.T) {
	r.T.Draw(args, scene)
	r.voxelCache.Tick()
}

func (r *voxrender) Destroy() {
	r.T.Destroy()
	r.voxelCache.Destroy()
}

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
			return NewVoxelRenderer(backend)
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
				Color:  color.RGBA(0.1, 0.1, 0.11, 0.85),
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
				image.New("cat", image.Props{
					Image: texture.PathRef("textures/kitten.png"),
					Style: image.Style{
						Width:  style.Pct(100),
						Height: style.Auto{},
					},
				}),
				textbox.New("testybox", textbox.Props{
					Style: textbox.Style{
						Bg: rect.Style{
							Color:   color.White,
							Padding: style.Px(4),
						},
						Text: label.Style{
							Color: color.Black,
							Width: style.Pct(100),
							Grow:  style.Grow(1),
						},
					},
				}),
				rect.New("objects", rect.Props{
					Children: []node.T{ObjectListEntry(0, scene)},
				}),
			},
		})
	}))
}

func ObjectListEntry(idx int, obj object.T) node.T {
	children := make([]node.T, len(obj.Children())+len(obj.Components())+1)
	clr := color.White
	if !obj.Active() {
		clr = color.RGB(0.7, 0.7, 0.7)
	}
	children[0] = label.New("title", label.Props{
		Text: "+ " + obj.Name(),
		Style: label.Style{
			Color: clr,
		},
	})
	i := 1
	for j, cmp := range obj.Components() {
		children[i] = label.New(fmt.Sprintf("component%d:%s", j, cmp.Name()), label.Props{
			Text: cmp.Name(),
			Style: label.Style{
				Color: clr,
			},
		})
		i++
	}
	for j, child := range obj.Children() {
		children[i] = ObjectListEntry(j, child)
		i++
	}
	return rect.New(fmt.Sprintf("object%d:%s", idx, obj.Name()), rect.Props{
		Style: rect.Style{
			Padding: style.Rect{
				Left: 5,
			},
		},
		Children: children,
	})
}
