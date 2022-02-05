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
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/ui"

	"github.com/lsfn/ode"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func ObjectHierarchy(idx int, obj object.T) node.T {
	children := make([]node.T, len(obj.Children())+1)
	clr := color.White
	if !obj.Active() {
		clr = color.RGB(0.7, 0.7, 0.7)
	}
	children[0] = label.New("title", &label.Props{
		Text:  obj.Name(),
		Color: clr,
	})
	for i, child := range obj.Children() {
		children[i+1] = ObjectHierarchy(i, child)
	}
	return rect.New(fmt.Sprintf("object%d", idx), &rect.Props{
		Layout: layout.Column{
			Padding: 4,
		},
		Children: children,
	})
}

func main() {
	fmt.Println("goworld")

	// cpu profiling
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		fmt.Println("writing cpu profiling output to", *cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	w := ode.NewWorld()
	w.SetGravity(ode.Vector3{0, -9.82, 0})

	scene := scene.New()

	wnd, err := window.New(window.Args{
		Title:        "goworld2",
		Width:        1600,
		Height:       1200,
		InputHandler: scene,
	})
	if err != nil {
		panic(err)
	}

	renderer := engine.NewRenderer(wnd)

	// attach GUI manager first.
	// this will give it input priority
	guim := gui.New(func() node.T {
		return rect.New("sidebar", &rect.Props{
			Layout: layout.Column{},
			Width:  dimension.Percent(15),
			Height: dimension.Percent(100),
			Children: []node.T{
				palette.New("palette", &palette.Props{
					Palette: color.DefaultPalette,
					OnPick: func(clr color.T) {
						editor := query.New[editor.T]().First(scene)
						if editor == nil {
							panic("could not find editor")
						}

						editor.SelectColor(clr)
					},
				}),
				image.New("diffuse", &image.Props{
					Image:  renderer.Geometry.Buffer.Diffuse(),
					Invert: true,
				}),
				image.New("normals", &image.Props{
					Image:  renderer.Geometry.Buffer.Normal(),
					Invert: true,
				}),
				image.New("position", &image.Props{
					Image:  renderer.Geometry.Buffer.Position(),
					Invert: true,
				}),
				rect.New("objects", &rect.Props{
					Color:    color.RGBA(0, 0, 0, 0.5),
					Children: []node.T{ObjectHierarchy(0, scene)},
				}),
			},
		})
	})

	scene.Adopt(guim)

	uim := ui.NewManager(1600, 1200)
	renderer.Append("ui", uim)
	renderer.Append("gui", guim)

	scene.Attach(light.NewDirectional(light.DirectionalArgs{
		Intensity: 1.2,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(0.95, -1.6, 1.05),
		Shadows:   true,
	}))
	scene.Attach(light.NewDirectional(light.DirectionalArgs{
		Intensity: 0.4,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(-1.2, -1.05, 1.12),
		Shadows:   true,
	}))

	gizmo := mover.New(mover.Args{})
	gizmo.Transform().SetPosition(vec3.New(-1, 0, -1))
	scene.Adopt(gizmo)

	// create chunk
	world := game.NewWorld(31481234, 16)
	chunk := world.AddChunk(0, 0)

	// first person controls
	player := game.NewPlayer(vec3.New(1, 22, 1), func(player *game.Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Flying = true
	player.Eye.Transform().SetRotation(vec3.New(22, 135, 0))
	scene.SetCamera(player.Camera)
	scene.Adopt(player)

	// create editor
	edit := editor.NewEditor(chunk, player.Camera, renderer.Geometry.Buffer)
	scene.Adopt(edit.Object())
	// uim.Attach(edit.Palette)

	// cube := geometry.NewColorCube(render.Blue, 3)
	// cube.Position = vec3.New(-2, -2, -2)

	// THIS DOES NOT ATTACH THE OBJECT
	// IT ATTACHES THE COMPONENT DIRECTLY!
	// scene.Attach(cube)

	// particles := engine.NewParticleSystem(vec3.New(3, 9, 3))
	// scene.Add(particles)

	fmt.Println("Ok")

	for !wnd.ShouldClose() {
		scene.Update(0.030)
		renderer.Draw(scene)
		wnd.SwapBuffers()
	}
}
