package main

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
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
		return rect.New("gui", rect.Props{
			Children: []node.T{
				makeMenu(),
				rect.New("gui-main", rect.Props{
					Style: rect.Style{
						Grow: style.Grow(1),
					},
					Children: []node.T{
						makeSidebar(scene, r),
					},
				}),
			},
		})
	}))
}

func makeMenu() node.T {
	return menu.Menu("gui-menu", menu.Props{
		Style: menu.Style{
			Color:      color.RGB(0.76, 0.76, 0.76),
			HoverColor: color.RGB(0.85, 0.85, 0.85),
			TextColor:  color.Black,
		},
	})
}

func makeSidebar(scene object.T, r renderer.T) node.T {
	return rect.New("sidebar", rect.Props{
		OnMouseDown: gui.ConsumeMouse,
		Style: rect.Style{
			Layout: style.Column{},
			Width:  style.Pct(15),
			Height: style.Pct(100),
			Color:  color.RGBA(0.1, 0.1, 0.11, 0.85),
		},
		Children: []node.T{
			image.New("logo", image.Props{
				Image: texture.PathRef("textures/shit_logo.png"),
				Style: image.Style{
					Width:  style.Pct(100),
					Height: style.Auto{},
				},
			}),

			// content placeholder
			rect.New("sidebar:content", rect.Props{}),

			ObjectListEntry("scene-graph", ObjectListEntryProps{
				Object: scene,
				OnSelect: func(obj object.T) {
					fmt.Println("selected", obj.Name())
				},
			}),
		},
	})
}

type SelectObjectHandler func(object.T)

type ObjectListEntryProps struct {
	Object   object.T
	OnSelect SelectObjectHandler
}

func ObjectListEntry(key string, props ObjectListEntryProps) node.T {
	return node.Component(key, props, nil, func(props ObjectListEntryProps) node.T {
		obj := props.Object
		clr := color.White
		if !obj.Active() {
			clr = color.RGB(0.7, 0.7, 0.7)
		}

		open, setOpen := hooks.UseState(false)
		icon := "+"
		if open {
			icon = "-"
		}

		title := rect.New("title-row", rect.Props{
			Style: rect.Style{
				Layout: style.Row{},
			},
			Children: []node.T{
				label.New("toggle", label.Props{
					Text: icon,
					OnClick: func(e mouse.Event) {
						setOpen(!open)
					},
					Style: label.Style{
						Color: clr,
					},
				}),
				label.New("title", label.Props{
					Text: obj.Name(),
					OnClick: func(e mouse.Event) {
						if props.OnSelect != nil {
							props.OnSelect(obj)
						}
					},
					Style: label.Style{
						Color: clr,
					},
				}),
			},
		})

		nodes := make([]node.T, 0, len(obj.Children())+1)
		nodes = append(nodes, title)

		if open {
			for idx, obj := range obj.Children() {
				key := fmt.Sprintf("object%d:%s", idx, obj.Name())
				nodes = append(nodes, ObjectListEntry(key, ObjectListEntryProps{
					Object:   obj,
					OnSelect: props.OnSelect,
				}))
			}
		}

		return rect.New(key, rect.Props{
			Style: rect.Style{
				Padding: style.Rect{
					Left: 5,
				},
			},
			Children: nodes,
		})
	})
}
