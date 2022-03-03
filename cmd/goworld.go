package main

/*
 * Copyright (C) 2016-2022 Johan Henriksson
 *
 * goworld is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * goworld is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with goworld. If not, see <http://www.gnu.org/licenses/>.
 */

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/geometry/gltf"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	engine.Run(engine.Args{
		Title:     "goworld",
		Width:     1600,
		Height:    1200,
		SceneFunc: makeScene,
	})
}

// the renderer object should probably not be exposed to the scene at all
// currently, access to the geometry buffer is needed for 2 things:
//
// - object picking (editor)
// - framebuffer debug windows
//
// in both cases, what we need to access is framebuffer textures

func makeScene(renderer engine.Renderer, scene object.T) {
	glrender := renderer.(*engine.GLRenderer)
	makeGui(glrender, scene)

	// create voxel chunk scene
	player, chunk := game.CreateScene(renderer, scene)

	// create editor
	edit := editor.NewEditor(chunk, player.Camera, glrender.Geometry.Buffer)
	scene.Adopt(edit.Object())

	// little gizmo thingy
	gizmo := mover.New(mover.Args{})
	gizmo.Transform().SetPosition(vec3.New(-1, 0, -1))
	scene.Adopt(gizmo)

	asd := gltf.Load(assets.GetMaterial("color.d"), "models/stack.glb")
	scene.Adopt(asd)
}

func makeGui(renderer *engine.GLRenderer, scene object.T) {
	scene.Attach(gui.New(func() node.T {
		return rect.New("sidebar", rect.Props{
			OnMouseDown: func(e mouse.Event) {},
			Style: rect.Style{
				Layout: style.Column{},
				Width:  style.Pct(15),
				Height: style.Pct(100),
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
				image.New("diffuse", image.Props{
					Image:  renderer.Geometry.Buffer.Diffuse(),
					Invert: true,
				}),
				image.New("normals", image.Props{
					Image:  renderer.Geometry.Buffer.Normal(),
					Invert: true,
				}),
				image.New("position", image.Props{
					Image:  renderer.Geometry.Buffer.Position(),
					Invert: true,
				}),
				rect.New("objects", rect.Props{
					Style: rect.Style{
						Color: color.Black.WithAlpha(0.9),
					},
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
		Text: obj.Name(),
		Style: label.Style{
			Color: clr,
		},
	})
	i := 1
	for j, cmp := range obj.Components() {
		children[i] = label.New(fmt.Sprintf("component%d:%s", j, cmp.Name()), label.Props{
			Text: fmt.Sprintf("+ %s", cmp.Name()),
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
