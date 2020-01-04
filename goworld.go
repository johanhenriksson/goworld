package main

/*
 * Copyright (C) 2016-2019 Johan Henriksson
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
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/render"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	WIDTH  = 1200
	HEIGHT = 800
)

func main() {
	fmt.Println("goworld")

	app := engine.NewApplication("voxels", WIDTH, HEIGHT)

	/* grab a reference to the geometry render pass */
	geoPass := app.Render.Get("geometry").(*engine.GeometryPass)

	// create a camera
	camera := engine.CreateCamera(&render.ScreenBuffer, 10, 10, 20, 65.0, 0.1, 1500.0)
	camera.Rotation[0] = 38
	camera.Rotation[1] = 230
	camera.Clear = render.Color4(0.141, 0.128, 0.118, 1.0)
	camera.Clear = render.Color4(0, 0, 0, 1)
	camera.Clear = render.Color4(0.368, 0.611, 0.800, 1.0)
	//camera.Clear = render.Color{0.973, 0.945, 0.776, 1.0}

	app.Scene.Camera = camera
	app.Scene.Lights = []engine.Light{
		{ // directional light
			Intensity:  0.8,
			Color:      mgl.Vec3{0.9 * 0.973, 0.9 * 0.945, 0.9 * 0.776},
			Type:       engine.DirectionalLight,
			Projection: mgl.Ortho(-200, 300, -30, 200, -200, 760),
			Position:   mgl.Vec3{-2, 1, -1},
		},
		{ // centered point light
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:     mgl.Vec3{0.517, 0.506, 0.447},
			Intensity: 1.0,
			Range:     70,
			Type:      engine.PointLight,
			Position:  mgl.Vec3{65, 27, 65},
		},
	}

	w := app.Scene.World
	w.NewPlane(0, 1, 0, 0)

	csize := 32
	ccount := 10

	fmt.Print("generating chunks... ")
	chunks := make([][]*game.ColorChunk, ccount)
	for cx := 0; cx < ccount; cx++ {
		chunks[cx] = make([]*game.ColorChunk, ccount)
		for cz := 0; cz < ccount; cz++ {
			obj := app.Scene.NewObject(float32(cx*csize), 0, float32(cz*csize))
			chk := game.NewColorChunk(obj, csize)
			chk.Seed = 31481234
			chk.Ox, chk.Oy, chk.Oz = cx*csize, 0, cz*csize
			generateChunk(chk, cx*csize, 0, cz*csize) // populate with random data
			chk.Compute()
			geoPass.Material.SetupVertexPointers() // wtfff
			app.Scene.Add(obj)

			chunks[cx][cz] = chk
			fmt.Printf("(%d,%d) ", cx, cz)
		}
	}
	fmt.Println("done")

	// this composition system sucks
	//game.NewPlacementGrid(chunks[0])

	// buffer display window
	winColor := render.Color4(0.15, 0.15, 0.15, 0.8)
	textColor := render.Color4(1, 1, 1, 1)

	bufferWindow := func(title string, texture *render.Texture, x, y float32, depth bool) {
		win := app.UI.NewRect(winColor, x, y, 240, 190, -10)
		label := app.UI.NewText(title, textColor, 0, 0, -21)
		win.Append(label)

		if depth {
			img := app.UI.NewDepthImage(texture, 0, 30, 240, 160, -20)
			img.Quad.FlipY()
			win.Append(img)
		} else {
			img := app.UI.NewImage(texture, 0, 30, 240, 160, -20)
			img.Quad.FlipY()
			win.Append(img)
		}

		app.UI.Append(win)
	}

	lightPass := app.Render.Get("light").(*engine.LightPass)
	bufferWindow("Diffuse", geoPass.Buffer.Diffuse, 10, 10, false)
	bufferWindow("Occlusion", lightPass.SSAO.Output, 10, 210, true)
	bufferWindow("Shadowmap", lightPass.Shadows.Output, 10, 410, true)

	paletteWindow := func(x, y float32, palette render.Palette) {
		win := app.UI.NewRect(winColor, x, y, 100, 180, -15)
		label := app.UI.NewText("Palette", textColor, 4, 8*20-4, -16)
		win.Append(label)

		perRow := 5
		for i, color := range palette {
			row := i / perRow
			col := i % perRow
			c := app.UI.NewRect(color, float32(col*20), float32(row*20), 20, 20, -17)
			win.Append(c)
		}

		app.UI.Append(win)
	}

	paletteWindow(300, 20, render.DefaultPalette)

	versiontext := fmt.Sprintf("goworld | %s", time.Now())
	watermark := app.UI.NewText(versiontext, render.Color4(1, 1, 1, 1), 400, 400, 0)
	app.UI.Append(watermark)

	paletteIdx := 5
	selected := game.NewColorVoxel(render.DefaultPalette[paletteIdx])

	sampleWorld := func() (mgl.Vec3, bool) {
		depth, depthExists := geoPass.Buffer.SampleDepth(int(engine.Mouse.X), int(engine.Mouse.Y))
		if !depthExists {
			return mgl.Vec3{}, false
		}
		return camera.Unproject(mgl.Vec3{
			engine.Mouse.X / float32(geoPass.Buffer.Depth.Width),
			engine.Mouse.Y / float32(geoPass.Buffer.Depth.Height),
			depth,
		}), true
	}
	sampleNormal := func() (mgl.Vec3, bool) {
		viewNormal, exists := geoPass.Buffer.SampleNormal(int(engine.Mouse.X), int(engine.Mouse.Y))
		if exists {
			viewInv := camera.View.Inv()
			worldNormal := viewInv.Mul4x1(viewNormal.Vec4(0)).Vec3()
			return worldNormal, true
		}
		return viewNormal, false
	}

	/* Render loop */
	app.UpdateFunc = func(dt float32) {
		versiontext = fmt.Sprintf("goworld | %s", time.Now())
		//watermark.Set(versiontext)

		world, worldExists := sampleWorld()
		if !worldExists {
			return
		}

		normal, normalExists := sampleNormal()
		if !normalExists {
			return
		}

		cx := int(world.X()) / csize
		cz := int(world.Z()) / csize
		if cx < 0 || cz < 0 || cx >= ccount || cz >= ccount {
			fmt.Println("target out of bounds")
			return
		}

		if engine.KeyReleased(engine.KeyG) {
			fmt.Println("raycast")
			w.Raycast(1000, app.Scene.Camera.Position, app.Scene.Camera.Forward)
		}

		if engine.KeyReleased(engine.KeyF) {
			paletteIdx++
			selected = game.NewColorVoxel(render.DefaultPalette[paletteIdx%len(render.DefaultPalette)])
		}

		if engine.KeyReleased(engine.KeyR) {
			paletteIdx--
			selected = game.NewColorVoxel(render.DefaultPalette[paletteIdx%len(render.DefaultPalette)])
		}

		// place voxel
		if engine.MouseDownPress(engine.MouseButton2) {
			fmt.Println("place at", world)
			target := world.Add(normal.Mul(0.5))
			chunks[cx][cz].Set(int(target[0])%csize, int(target[1])%csize, int(target[2])%csize, selected)
			chunks[cx][cz].Compute()
		}

		// remove voxel
		if engine.KeyPressed(engine.KeyC) {
			fmt.Println("delete from", world)
			target := world.Sub(normal.Mul(0.5))
			chunks[cx][cz].Set(int(target[0])%csize, int(target[1])%csize, int(target[2])%csize, nil)
			chunks[cx][cz].Compute()
		}

		// eyedropper
		if engine.KeyPressed(engine.KeyI) {
			target := world.Sub(normal.Mul(0.5))
			selected = chunks[cx][cz].At(int(target[0])%csize, int(target[1])%csize, int(target[2])%csize)
		}
	}

	fmt.Println("ok")
	app.Run()
}

// ChunkFunc is a chunk function :)
//type ChunkFunc func(*game.Chunk, ChunkFuncParams)

func generateChunk(chk *game.ColorChunk, ox int, oy int, oz int) {
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
	cloud := &game.ColorVoxel{
		R: 255,
		G: 255,
		B: 255,
	}

	/* Fill chunk with voxels */
	size := chk.Size

	rockNoise := math.NewNoise(chk.Seed+10000, 1.0/40.0)
	grassNoise := math.NewNoise(chk.Seed+10002, 1.0/28.0)
	cloudNoise := math.NewNoise(chk.Seed+24511626, 1/40.0)

	grassHeight := 8

	for z := 0; z < size; z++ {
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				gh := int(9 * grassNoise.Sample(x+ox, oy, z+oz))
				rh := int(44 * rockNoise.Sample(x+ox, oy, z+oz))
				ch := int(8*cloudNoise.Sample(x+ox, y+oy, z+oz)) + 8

				var vtype *game.ColorVoxel
				if y < grassHeight {
					vtype = rock2
				}

				if y == grassHeight {
					vtype = grass
				}
				if y < grassHeight+gh && y > grassHeight {
					vtype = grass
				}
				if y < rh {
					vtype = rock
				}

				if ch > 12 && y > 98-ch && y < 100+ch {
					vtype = cloud
				}

				chk.Set(x, y, z, vtype)
			}
		}
	}
}
