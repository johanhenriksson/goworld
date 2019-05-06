package main

/*
 * Copyright (C) 2016 Johan Henriksson
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

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render"

	opensimplex "github.com/ojrac/opensimplex-go"
)

const (
	WIDTH  = 800
	HEIGHT = 600
)

func main() {
	app := engine.NewApplication("voxels", WIDTH, HEIGHT)

	/* grab a reference to the geometry render pass */
	geom_pass := app.Render.Get("geometry").(*engine.GeometryPass)
	light_pass := app.Render.Get("light").(*engine.LightPass)

	/* create a camera */

	width, height := app.Window.GetBufferSize()
	app.Scene.Camera = engine.CreateCamera(-3, 2, -3, float32(width), float32(height), 65.0, 0.1, 500.0)
	app.Scene.Camera.Transform.Rotation[1] = 130.0

	w := app.Scene.World
	w.NewPlane(0, 1, 0, 0)

	obj2 := app.Scene.NewObject(5, 0, 5)
	chk2 := game.NewColorChunk(obj2, 32)
	generateChunk(chk2) // populate with random data
	chk2.Set(0, 0, 0, &game.ColorVoxel{R: 255, G: 0, B: 0})
	chk2.Set(1, 0, 0, &game.ColorVoxel{R: 0, G: 255, B: 0})
	chk2.Set(2, 0, 0, &game.ColorVoxel{R: 0, G: 0, B: 255})
	chk2.Compute()
	geom_pass.Material.SetupVertexPointers()
	app.Scene.Add(obj2)

	game.NewPlacementGrid(obj2)

	fmt.Println("goworld")

	// buffer display window
	bufferWindow := func(title string, texture *render.Texture, x, y float32, depth bool) {
		win_color := render.Color{0.15, 0.15, 0.15, 0.8}
		text_color := render.Color{1, 1, 1, 1}

		win := app.UI.NewRect(win_color, x, y, 250, 280, -10)
		label := app.UI.NewText(title, text_color, 0, 0, -21)
		win.Append(label)

		if depth {
			img := app.UI.NewDepthImage(texture, 0, 30, 250, 250, -20)
			img.Quad.FlipY()
			win.Append(img)
		} else {
			img := app.UI.NewImage(texture, 0, 30, 250, 250, -20)
			img.Quad.FlipY()
			win.Append(img)
		}

		/* attach UI element */
		app.UI.Append(win)
	}

	bufferWindow("Diffuse", geom_pass.Buffer.Diffuse, 30, 30, false)
	bufferWindow("Normal", geom_pass.Buffer.Normal, 30, 340, false)
	bufferWindow("Shadowmap", light_pass.Shadows.Output, 30, 650, true)

	/* Render loop */
	app.UpdateFunc = func(dt float32) {
		if engine.KeyReleased(engine.KeyF) {
			fmt.Println("raycast")
			w.Raycast(10, app.Scene.Camera.Position, app.Scene.Camera.Forward)
		}
	}

	app.Run()
}

func generateChunk(chk *game.ColorChunk) {
	/* Define voxels */
	rock2 := &game.ColorVoxel{
		R: 137,
		G: 131,
		B: 119,
	}
	rock := &game.ColorVoxel{
		R: 173,
		G: 169,
		B: 158,
	}
	grass := &game.ColorVoxel{
		R: 122,
		G: 189,
		B: 64,
	}

	/* Fill chunk with voxels */
	f := 1.0 / 30.0
	size := chk.Size
	simplex := opensimplex.New(3001)
	for z := 0; z < size; z++ {
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				fx, fy, fz := float64(x)*f, float64(y)*f, float64(z)*f
				v := simplex.Eval3(fx, fy, fz)
				var vtype *game.ColorVoxel = nil
				if y < size/4 {
					vtype = rock2
				}
				if y == size/4 {
					vtype = grass
				}
				if v < -0.3 {
					vtype = rock
				}
				chk.Set(x, y, z, vtype)
			}
		}
	}
}
