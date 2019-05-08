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
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v2.1/gl"
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
	//lightPass := app.Render.Get("light").(*engine.LightPass)

	/* create a camera */

	width, height := app.Window.GetBufferSize()
	camera := engine.CreateCamera(100, 90, -20, float32(width), float32(height), 65.0, 0.1, 1500.0)
	camera.Rotation[0] = 38
	camera.Rotation[1] = 230
	camera.Clear = render.Color4(0.141, 0.128, 0.118, 1.0)
	camera.Clear = render.Color4(0, 0, 0, 1)
	camera.Clear = render.Color4(0.368, 0.611, 0.800, 1.0)
	//camera.Clear = render.Color{0.973, 0.945, 0.776, 1.0}

	app.Scene.Camera = camera
	app.Scene.Lights = []engine.Light{
		{ // directional light
			Color:      mgl.Vec3{0.5 * 0.973, 0.5 * 0.945, 0.5 * 0.776},
			Type:       engine.DirectionalLight,
			Projection: mgl.Ortho(-320, 580, -30, 600, -320, 760),
			Position:   mgl.Vec3{-2, 1, -1},
		},
		{ // centered point light
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:    mgl.Vec3{0.517, 0.506, 0.447},
			Range:    80,
			Type:     engine.PointLight,
			Position: mgl.Vec3{26, 36, 32},
		},
		{ // centered point light
			Attenuation: engine.Attenuation{
				Constant:  0.50,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:    mgl.Vec3{0.517, 0.506, 0.447},
			Range:    80,
			Type:     engine.PointLight,
			Position: mgl.Vec3{-16, 36, 32},
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
	/*
		bufferWindow := func(title string, texture *render.Texture, x, y float32, depth bool) {
			winColor := render.Color{0.15, 0.15, 0.15, 0.8}
			textColor := render.Color{1, 1, 1, 1}

			win := app.UI.NewRect(winColor, x, y, 250, 280, -10)
			label := app.UI.NewText(title, textColor, 0, 0, -21)
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

			app.UI.Append(win)
		}

		bufferWindow("Normal", geoPass.Buffer.Normal, 30, 30, false)
			bufferWindow("Occlusion", lightPass.SSAO.Gaussian.Output, 30, 340, false)
			bufferWindow("Shadowmap", lightPass.Shadows.Output, 30, 650, true)
	*/

	versiontext := fmt.Sprintf("goworld | %s", time.Now())
	watermark := app.UI.NewText(versiontext, render.Color4(1, 1, 1, 1), WIDTH-300, 0, 0)
	app.UI.Append(watermark)

	paletteIdx := 5
	selected := game.NewColorVoxel(render.DefaultPalette[paletteIdx])

	sampleNormal := func(x, y float32) mgl.Vec3 {
		geoPass.Buffer.Bind()
		normalEncoded := geoPass.Buffer.Sample(gl.COLOR_ATTACHMENT1, int(x), int(HEIGHT-y))
		if normalEncoded.R != 0 || normalEncoded.G != 0 || normalEncoded.B != 0 {
			viewNormal := normalEncoded.Vec3().Mul(2).Sub(mgl.Vec3{1, 1, 1}).Normalize() // normals [-1,1]

			viewInv := camera.View.Inv()
			worldNormal := viewInv.Mul4x1(mgl.Vec4{viewNormal[0], viewNormal[1], viewNormal[2], 0}).Vec3()
			return worldNormal
		}
		return normalEncoded.Vec3() // zero
	}

	/* Render loop */
	app.UpdateFunc = func(dt float32) {
		versiontext = fmt.Sprintf("goworld | %s", time.Now())
		watermark.Set(versiontext)

		geoPass.Buffer.Bind()
		world := camera.Unproject(engine.Mouse.X, engine.Mouse.Y)
		normal := sampleNormal(engine.Mouse.X, engine.Mouse.Y)
		if normal[0] == 0 && normal[1] == 0 && normal[2] == 0 {
			return
		}

		cx := int(world.X()) / csize
		cz := int(world.Z()) / csize

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
		if engine.MouseDownPress(1) {
			target := world.Add(normal.Mul(0.5))
			chunks[cx][cz].Set(int(target[0])%csize, int(target[1])%csize, int(target[2])%csize, selected)
			chunks[cx][cz].Compute()
		}

		// remove voxel
		if engine.KeyPressed(engine.KeyC) {
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

	rockNoise := math.NewNoise(10000, 1.0/40.0)
	grassNoise := math.NewNoise(10002, 1.0/28.0)
	cloudNoise := math.NewNoise(24511626, 1/40.0)

	grassHeight := 8

	for z := 0; z < size; z++ {
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				gh := int(9 * grassNoise.Sample(x+ox, oy, z+oz))
				rh := int(44 * rockNoise.Sample(x+ox, oy, z+oz))
				ch := int(8*cloudNoise.Sample(x+ox, y+oy, z+oz)) + 8

				var vtype *game.ColorVoxel = nil
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
