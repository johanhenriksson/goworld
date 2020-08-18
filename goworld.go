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

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"

	mgl "github.com/go-gl/mathgl/mgl32"
)

var winColor = render.Color4(0.15, 0.15, 0.15, 1)
var textColor = render.Color4(1, 1, 1, 1)

var windowStyle = ui.Style{
	"background": ui.Color(winColor),
	"radius":     ui.Float(3),
	"padding":    ui.Float(5),
}

func main() {
	fmt.Println("goworld")

	app := engine.NewApplication("voxels", 1200, 800)
	uim := ui.NewManager(app)
	app.Render.Append("ui", uim)

	rect := ui.NewRect(windowStyle,
		ui.NewRect(ui.Style{"layout": ui.String("row"), "spacing": ui.Float(100)},
			ui.NewText("Hello Really Long Line", ui.NoStyle),
			ui.NewText("Please", ui.NoStyle)),
		ui.NewRect(ui.NoStyle,
			ui.NewText("Please", ui.NoStyle),
			ui.NewText("Hello Really Long Line", ui.NoStyle)))
	uim.Attach(rect)
	rect.SetPosition(400, 400)
	rect.DesiredSize(200, 1000)

	/* grab a reference to the geometry render pass */
	geoPass := app.Render.Get("geometry").(*engine.GeometryPass)

	// create a camera
	camera := engine.CreateCamera(&render.ScreenBuffer, -10, 22, -10, 55.0, 0.1, 600.0)
	camera.Rotation[0] = 22
	camera.Rotation[1] = 135
	camera.Clear = render.Color4(0.141, 0.128, 0.118, 1.0) // dark gray
	camera.Clear = render.Color4(0, 0, 0, 1)
	camera.Clear = render.Color4(0.368, 0.611, 0.800, 1.0) // blue
	//camera.Clear = render.Color{0.973, 0.945, 0.776, 1.0} // light gray

	app.Scene.Camera = camera
	app.Scene.Lights = []engine.Light{
		{ // directional light
			Intensity:  0.8,
			Color:      mgl.Vec3{0.9 * 0.973, 0.9 * 0.945, 0.9 * 0.776},
			Type:       engine.DirectionalLight,
			Projection: mgl.Ortho(-200, 300, -30, 250, -200, 760),
			Position:   mgl.Vec3{-3, 2, -2},
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

	csize := 16
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
			app.Scene.Add(obj)

			chunks[cx][cz] = chk
			fmt.Printf("(%d,%d) ", cx, cz)
		}
	}
	fmt.Println("done")

	// test cube
	building := app.Scene.NewObject(4.5, 9.04, 8.5)
	building.Scale = mgl.Vec3{0.1, 0.1, 0.1}
	palette := assets.GetMaterialCached("uv_palette")
	geometry.NewObjModel(building, palette, "models/building.obj")

	app.Scene.Add(building)

	// this composition system sucks
	//game.NewPlacementGrid(chunks[0])

	// buffer display windows
	lightPass := app.Render.Get("light").(*engine.LightPass)
	bufferWindows := ui.NewRect(ui.Style{"spacing": ui.Float(10)},
		newBufferWindow("Diffuse", geoPass.Buffer.Diffuse, 10, 10, false),
		newBufferWindow("Occlusion", lightPass.SSAO.Gaussian.Output, 10, 215, true),
		newBufferWindow("Shadowmap", lightPass.Shadows.Output, 10, 420, true))
	bufferWindows.SetPosition(10, 10)
	bufferWindows.DesiredSize(500, 1000)
	uim.Attach(bufferWindows)

	// palette globals
	paletteIdx := 5
	selected := game.NewColorVoxel(render.DefaultPalette[paletteIdx])

	paletteWnd := newPaletteWindow(render.DefaultPalette, func(newPaletteIdx int) {
		paletteIdx = newPaletteIdx
		selected = game.NewColorVoxel(render.DefaultPalette[paletteIdx])
	})
	paletteWnd.SetPosition(280, 10)
	paletteWnd.DesiredSize(200, 400)
	uim.Attach(paletteWnd)

	// watermark / fps text
	versiontext := fmt.Sprintf("goworld")
	watermark := ui.NewText(versiontext, ui.Style{"color": ui.Color(render.White)})
	watermark.SetPosition(10, float32(app.Window.Height-30))
	watermark.Texture.Save("test.png")
	uim.Attach(watermark)

	// sample world position at current mouse coords
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

	// sample world normal at current mouse coords
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
		versiontext = fmt.Sprintf("goworld | %s | %.0f fps", time.Now().Format("2006-01-02 15:04"), app.Window.FPS)
		watermark.Set(versiontext)

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
			return
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
			chunks[cx][cz].Set(int(target[0])%csize, int(target[1])%csize, int(target[2])%csize, game.EmptyColorVoxel)
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

func newPaletteWindow(palette render.Palette, onClickItem func(int)) ui.Component {
	cols := 5
	spacing := ui.Float(2)
	gridStyle := ui.Style{"spacing": spacing}
	rowStyle := ui.Style{"layout": ui.String("row"), "spacing": spacing}
	rows := make([]ui.Component, 0, len(palette)/cols+1)
	row := make([]ui.Component, 0, cols)

	for i := 1; i <= len(palette); i++ {
		itemIdx := i - 1
		color := palette[itemIdx]

		swatch := ui.NewRect(ui.Style{"background": ui.Color(color), "layout": ui.String("fixed")})
		swatch.SetSize(20, 20)
		swatch.OnClick(func(ev ui.MouseEvent) {
			if ev.Button == engine.MouseButton1 {
				onClickItem(itemIdx)
			}
		})

		row = append(row, swatch)

		if i%cols == 0 {
			rows = append(rows, ui.NewRect(rowStyle, row...))
			row = make([]ui.Component, 0, cols)
		}
	}

	return ui.NewRect(windowStyle,
		ui.NewText("Palette", ui.NoStyle),
		ui.NewRect(gridStyle, rows...))
}

func newBufferWindow(title string, texture *render.Texture, x, y float32, depth bool) ui.Component {
	var img ui.Component
	if depth {
		img = ui.NewDepthImage(texture, 240, 160, false)
	} else {
		img = ui.NewImage(texture, 240, 160, false)
	}

	return ui.NewRect(windowStyle,
		ui.NewText(title, ui.NoStyle),
		img)
}

// ChunkFunc is a chunk function :)
//type ChunkFunc func(*game.Chunk, ChunkFuncParams)

func generateChunk(chk *game.ColorChunk, ox int, oy int, oz int) {
	/* Define voxels */
	rock2 := game.ColorVoxel{
		R: 137,
		G: 131,
		B: 119,
	}
	rock := game.ColorVoxel{
		R: 173,
		G: 169,
		B: 158,
	}
	grass := game.ColorVoxel{
		R: 72,
		G: 140,
		B: 54,
	}
	cloud := game.ColorVoxel{
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

				var vtype game.ColorVoxel
				if y < grassHeight {
					vtype = rock2
				}

				if y == grassHeight {
					vtype = grass
				}
				if y <= grassHeight+gh && y > grassHeight {
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
