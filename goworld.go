package main

/*
 * Copyright (C) 2016-2020 Johan Henriksson
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
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

var winColor = render.Color4(0.15, 0.15, 0.15, 0.85)
var textColor = render.Color4(1, 1, 1, 1)

var windowStyle = ui.Style{
	"color":   ui.Color(winColor),
	"radius":  ui.Float(3),
	"padding": ui.Float(5),
}

func main() {
	fmt.Println("goworld")

	app := engine.NewApplication("voxels", 1400, 1000)
	uim := ui.NewManager(app)
	app.Render.Append("ui", uim)

	/* grab a reference to the geometry render pass */
	geoPass := app.Render.Get("geometry").(*engine.GeometryPass)

	// create a camera
	camera := engine.CreateCamera(&render.ScreenBuffer, vec3.T{1, 22, 1}, 55.0, 0.1, 600.0)
	camera.Rotation.X = 22
	camera.Rotation.Y = 135
	camera.Clear = render.Color4(0.141, 0.128, 0.118, 1.0) // dark gray
	camera.Clear = render.Color4(0.973, 0.945, 0.776, 1.0) // light gray
	camera.Clear = render.Color4(0.368, 0.611, 0.800, 1.0) // blue

	app.Scene.Camera = camera
	app.Scene.Lights = []engine.Light{
		{ // directional light
			Intensity:  0.8,
			Color:      vec3.T{0.9 * 0.973, 0.9 * 0.945, 0.9 * 0.776},
			Type:       engine.DirectionalLight,
			Projection: mat4.Orthographic(-71, 120, -20, 140, -10, 140),
			Position:   vec3.T{-2, 2, -1},
			Shadows:    false,
		},
		{ // light
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:     vec3.T{0.517, 0.506, 0.447},
			Intensity: 1.0,
			Range:     70,
			Type:      engine.PointLight,
			Position:  vec3.T{16, 30, 16},
		},
		{ // text highlight
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:     vec3.T{0.517, 0.506, 0.447},
			Intensity: 8.0,
			Range:     30,
			Type:      engine.PointLight,
			Position:  vec3.T{30, 35, 52},
		},
	}

	csize := 16
	ccount := 8

	world := game.NewWorld(31481234, csize)

	fmt.Print("Generating chunks... ")
	chunks := make([][]*game.ChunkMesh, ccount)
	for cx := 0; cx < ccount; cx++ {
		chunks[cx] = make([]*game.ChunkMesh, ccount)
		for cz := 0; cz < ccount; cz++ {

			obj := engine.NewObject(vec3.NewI(cx, 0, cz).ScaleI(csize))
			chunk := world.AddChunk(cx, cz)
			mesh := game.NewChunkMesh(obj, chunk)
			mesh.Compute()
			app.Scene.Add(obj)

			chunks[cx][cz] = mesh
			fmt.Printf("(%d,%d) ", cx, cz)
		}
	}
	fmt.Println("World generation complete")

	// test model
	// building := engine.NewObject(4.5, 9.04, 8.5)
	// building.Scale = mgl.Vec3{0.1, 0.1, 0.1}
	// palette := assets.GetMaterialCached("uv_palette")
	// geometry.NewObjModel(building, palette, "models/building.obj")
	// app.Scene.Add(building)

	// this composition system sucks
	//game.NewPlacementGrid(chunks[0])

	// buffer display windows
	lightPass := app.Render.Get("light").(*engine.LightPass)
	bufferWindows := ui.NewRect(ui.Style{"spacing": ui.Float(10)},
		newBufferWindow("Diffuse", geoPass.Buffer.Diffuse, false),
		newBufferWindow("Normal", geoPass.Buffer.Normal, false),
		newBufferWindow("Occlusion", lightPass.SSAO.Gaussian.Output, true),
		newBufferWindow("Shadowmap", lightPass.Shadows.Output, true))
	bufferWindows.SetPosition(10, 10)
	bufferWindows.Flow(ui.Size{500, 1000})
	uim.Attach(bufferWindows)

	// palette globals
	paletteIdx := 5
	selected := game.NewVoxel(render.DefaultPalette[paletteIdx])

	paletteWnd := newPaletteWindow(render.DefaultPalette, func(newPaletteIdx int) {
		paletteIdx = newPaletteIdx
		selected = game.NewVoxel(render.DefaultPalette[paletteIdx])
	})
	paletteWnd.SetPosition(280, 10)
	paletteWnd.Flow(ui.Size{200, 400})
	uim.Attach(paletteWnd)

	// watermark / fps text
	versiontext := fmt.Sprintf("goworld")
	watermark := ui.NewText(versiontext, ui.Style{"color": ui.Color(render.White)})
	watermark.SetPosition(10, float32(app.Window.Height-30))
	uim.Attach(watermark)

	// uv_checker := assets.GetTexture("textures/uv_checker.png")
	// uv_checker.Border = 50
	// br := ui.NewRect(ui.Style{
	// 	"radius":  ui.Float(25),
	// 	"padding": ui.Float(10),
	// 	"color":   ui.Color(render.White),
	// 	"image":   ui.Texture(uv_checker),
	// })
	// br.SetPosition(500, 300)
	// br.Resize(ui.Size{200, 200})
	// uim.Attach(br)

	// sample world position at current mouse coords
	sampleWorld := func() (vec3.T, bool) {
		depth, depthExists := geoPass.Buffer.SampleDepth(int(mouse.X), int(mouse.Y))
		if !depthExists {
			return vec3.Zero, false
		}
		return camera.Unproject(vec3.T{
			mouse.X / float32(geoPass.Buffer.Depth.Width),
			mouse.Y / float32(geoPass.Buffer.Depth.Height),
			depth,
		}), true
	}

	// sample world normal at current mouse coords
	sampleNormal := func() (vec3.T, bool) {
		viewNormal, exists := geoPass.Buffer.SampleNormal(int(mouse.X), int(mouse.Y))
		if exists {
			viewInv := camera.View.Inv()
			worldNormal := viewInv.MulVec4(vec4.Extend(viewNormal, 0)).XYZ()
			return worldNormal, true
		}
		return viewNormal, false
	}

	// physics constants
	gravity := float32(53)
	speed := float32(60)
	airspeed := float32(33)
	jumpvel := 0.25 * gravity
	friction := vec3.T{0.91, 1, 0.91}
	airfriction := vec3.T{0.955, 1, 0.955}
	camOffset := vec3.T{0, 1.75, 0}
	fly := false

	// player physics state
	position := camera.Position.Sub(camOffset)
	velocity := vec3.T{0, 0, 0}
	grounded := true

	/* Render loop */
	app.UpdateFunc = func(dt float32) {
		versiontext = fmt.Sprintf("goworld | %s | %.0f fps", time.Now().Format("2006-01-02 15:04"), app.Window.FPS)
		watermark.Set(versiontext)

		/*** movement **************************************/

		move := vec3.T{0, 0, 0}
		moving := false
		if keys.Down(keys.W) && !keys.Down(keys.S) {
			move.Z += 1.0
			moving = true
		}
		if keys.Down(keys.S) && !keys.Down(keys.W) {
			move.Z -= 1.0
			moving = true
		}
		if keys.Down(keys.A) && !keys.Down(keys.D) {
			move.X -= 1.0
			moving = true
		}
		if keys.Down(keys.D) && !keys.Down(keys.A) {
			move.X += 1.0
			moving = true
		}
		if fly && keys.Down(keys.Q) && !keys.Down(keys.E) {
			move.Y -= 1.0
			moving = true
		}
		if fly && keys.Down(keys.E) && !keys.Down(keys.Q) {
			move.Y += 1.0
			moving = true
		}
		if keys.Pressed(keys.V) {
			fly = !fly
		}

		if moving {
			right := camera.Transform.Right.Scaled(move.X)
			forward := camera.Transform.Forward.Scaled(move.Z)
			up := vec3.T{0, move.Y, 0}

			move = right.Add(forward)
			move.Y = 0 // remove y component
			if fly {
				move = move.Add(up)
			}
			move.Normalize()
		}
		if grounded || fly {
			move.Scale(speed)
		} else {
			move.Scale(airspeed)
		}

		if keys.Down(keys.LeftShift) {
			move.Scale(2)
		}

		// apply movement
		velocity = velocity.Add(move.Scaled(dt))

		// friction
		if grounded {
			velocity = velocity.Mul(friction)
		} else {
			velocity = velocity.Mul(airfriction)
		}

		// gravity
		if !fly {
			velocity.Y -= gravity * dt
		} else {
			// apply Y friction while flying
			velocity.Y *= airfriction.X
		}

		step := velocity.Scaled(dt)

		// apply movement in Y
		position.Y += step.Y
		step.Y = 0

		// ground collision
		height := world.HeightAt(position)
		if position.Y < height {
			position.Y = height
			velocity.Y = 0
			grounded = true
		} else {
			grounded = false
		}

		// jumping
		if grounded && keys.Down(keys.Space) {
			velocity.Y += jumpvel
		}

		// x collision
		xstep := position.Add(vec3.T{step.X, 0, 0})
		if world.HeightAt(xstep) > position.Y {
			step.X = 0
		}

		// z collision
		zstep := position.Add(vec3.T{0, 0, step.Z})
		if world.HeightAt(zstep) > position.Y {
			step.Z = 0
		}

		// add horizontal movement
		position = position.Add(step)

		// update camera position
		camera.Position = position.Add(camOffset)

		/*** end movement **************************************/

		worldPos, worldExists := sampleWorld()
		if !worldExists {
			return
		}

		normal, normalExists := sampleNormal()
		if !normalExists {
			return
		}

		cx := int(worldPos.X) / csize
		cz := int(worldPos.Z) / csize
		if cx < 0 || cz < 0 || cx >= ccount || cz >= ccount {
			return
		}
		chunk := chunks[cx][cz]

		if keys.Released(keys.R) {
			// replace voxel
			fmt.Println("Replace at", worldPos)
			target := worldPos.Sub(normal.Scaled(0.5))
			world.Set(int(target.X), int(target.Y), int(target.Z), selected)

			// recompute mesh
			chunk.Light.Calculate()
			chunk.Compute()

			// write to disk
			go chunk.Write("chunks")
		}

		// place voxel
		if mouse.Pressed(mouse.Button2) {
			fmt.Println("Place at", worldPos)
			target := worldPos.Add(normal.Scaled(0.5))
			world.Set(int(target.X), int(target.Y), int(target.Z), selected)

			// recompute mesh
			chunk.Light.Calculate()
			chunk.Compute()

			// write to disk
			go chunk.Write("chunks")
		}

		// remove voxel
		if keys.Pressed(keys.C) {
			fmt.Println("Delete from", worldPos)
			target := worldPos.Sub(normal.Scaled(0.5))
			world.Set(int(target.X), int(target.Y), int(target.Z), game.EmptyVoxel)

			// recompute mesh
			chunk.Light.Calculate()
			chunk.Compute()

			// write to disk
			go chunk.Write("chunks")
		}

		// eyedropper
		if keys.Pressed(keys.F) {
			fmt.Println("Sample", worldPos)
			target := worldPos.Sub(normal.Scaled(0.5))
			selected = world.Voxel(int(target.X), int(target.Y), int(target.Z))
		}
	}

	fmt.Println("Ok")
	app.Run()
}

func newPaletteWindow(palette render.Palette, onClickItem func(int)) ui.Component {
	cols := 5
	gridStyle := ui.Style{"layout": ui.String("column"), "spacing": ui.Float(2)}
	rowStyle := ui.Style{"layout": ui.String("row"), "spacing": ui.Float(2)}
	rows := make([]ui.Component, 0, len(palette)/cols+1)
	row := make([]ui.Component, 0, cols)

	for i := 1; i <= len(palette); i++ {
		itemIdx := i - 1
		color := palette[itemIdx]

		swatch := ui.NewRect(ui.Style{"color": ui.Color(color), "layout": ui.String("fixed")})
		swatch.Resize(ui.Size{20, 20})
		swatch.OnClick(func(ev ui.MouseEvent) {
			if ev.Button == mouse.Button1 {
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

func newBufferWindow(title string, texture *render.Texture, depth bool) ui.Component {
	var img ui.Component
	if depth {
		img = ui.NewDepthImage(texture, 240, 160, false)
	} else {
		img = ui.NewImage(texture, 240, 160, false, ui.NoStyle)
	}

	return ui.NewRect(windowStyle,
		ui.NewText(title, ui.NoStyle),
		img)
}

// ChunkFunc is a chunk function :)
//type ChunkFunc func(*game.Chunk, ChunkFuncParams)
