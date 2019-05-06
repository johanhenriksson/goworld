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
	"time"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render"

	opensimplex "github.com/ojrac/opensimplex-go"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	WIDTH  = 1600
	HEIGHT = 1000
)

func main() {
	app := engine.NewApplication("voxels", WIDTH, HEIGHT)

	/* grab a reference to the geometry render pass */
	geom_pass := app.Render.Get("geometry").(*engine.GeometryPass)
	light_pass := app.Render.Get("light").(*engine.LightPass)

	/* create a camera */

	width, height := app.Window.GetBufferSize()
	app.Scene.Camera = engine.CreateCamera(10, 20, 10, float32(width), float32(height), 65.0, 0.1, 500.0)
	app.Scene.Camera.Transform.Rotation[1] = 130.0
	app.Scene.Lights = []engine.Light{
		{ // directional light
			Attenuation: engine.Attenuation{
				Constant:  0.01,
				Linear:    0,
				Quadratic: 1.0,
			},
			Color: mgl.Vec3{0.35, 0.35, 0.35},
			Range: 4,
			Type:  engine.DirectionalLight,
			Projection: mgl.Ortho(-32, 32, 0, 64, -32, 64),
			Position: mgl.Vec3{-11, 16, -11},
		},
		{ // centered point light
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color: mgl.Vec3{1, 1, 1},
			Range: 20,
			Type:  engine.PointLight,
			Position: mgl.Vec3{ 32, 36, 32},
		},
	}

	w := app.Scene.World
	w.NewPlane(0, 1, 0, 0)

	obj2 := app.Scene.NewObject(-2, 0, -2)
	chk2 := game.NewColorChunk(obj2, 64)
	generateChunk(chk2) // populate with random data
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

	versiontext := fmt.Sprintf("goworld | %s", time.Now())
	watermark := app.UI.NewText(versiontext, render.Color{1,1,1,1}, WIDTH - 200, 0, 0)
	app.UI.Append(watermark)

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
		R: 72,
		G: 140,
		B: 54,
	}

	/* Fill chunk with voxels */
	size := chk.Size

	rockNoise := NewNoise(3001, 1.0 / 19.0)
	grassNoise := NewNoise(314159, 1.0 / 28.0)

	grassHeight := size / 4

	for z := 0; z < size; z++ {
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				vr := rockNoise.Sample(x, y, z)
				vg := grassNoise.Sample(x, 0, z)
				gh := int(vg * 9)

				var vtype *game.ColorVoxel = nil
				if y < grassHeight {
					vtype = rock2
				}

				if y == grassHeight {
					vtype = grass
				}
				if y < grassHeight + gh && y > grassHeight {
					vtype = grass
				}
				if vr < -0.3 {
					vtype = rock
				}
				chk.Set(x, y, z, vtype)
			}
		}
	}
}

type Noise struct {
	opensimplex.Noise
	Seed int
	Freq float64
}

func NewNoise(seed int, freq float64) *Noise {
	return &Noise{
		Noise: opensimplex.New(int64(seed)),
		Seed: seed,
		Freq: freq,
	}
}

func (n *Noise) Sample(x, y, z int) float64 {
	fx, fy, fz := float64(x) * n.Freq, float64(y) * n.Freq, float64(z) * n.Freq
	return n.Eval3(fx, fy, fz)
}