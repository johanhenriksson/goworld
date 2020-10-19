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

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

func main() {
	fmt.Println("goworld")

	app := engine.NewApplication("goworld", 1400, 1000)
	uim := ui.NewManager(app)
	app.Pipeline.Append("ui", uim)

	scene := engine.NewScene()

	/* grab a reference to the geometry render pass */
	geoPass := app.Pipeline.Get("geometry").(*engine.GeometryPass)

	// create a camera
	camera := engine.CreateCamera(&render.ScreenBuffer, vec3.New(1, 22, 1), 55.0, 0.1, 600.0)
	camera.Rotation.X = 22
	camera.Rotation.Y = 135
	camera.Clear = render.Color4(0.141, 0.128, 0.118, 1.0) // dark gray
	camera.Clear = render.Color4(0.368, 0.611, 0.800, 1.0) // blue
	camera.Clear = render.Color4(0.973, 0.945, 0.876, 1.0) // light gray

	scene.Camera = camera
	scene.Lights = []engine.Light{
		{ // directional light
			Intensity:  0.8,
			Color:      vec3.New(0.9*0.973, 0.9*0.945, 0.9*0.776),
			Type:       engine.DirectionalLight,
			Projection: mat4.Orthographic(-71, 120, -20, 140, -10, 140),
			Position:   vec3.New(-2, 2, -1),
			Shadows:    false,
		},
		{ // light
			Attenuation: engine.Attenuation{
				Constant:  1.00,
				Linear:    0.09,
				Quadratic: 0.32,
			},
			Color:     vec3.New(0.517, 0.506, 0.447),
			Intensity: 1.0,
			Range:     70,
			Type:      engine.PointLight,
			Position:  vec3.New(16, 30, 16),
		},
	}

	csize := 16
	ccount := 1

	world := game.NewWorld(31481234, csize)

	fmt.Print("Generating chunks... ")
	chunks := make([][]*game.ChunkMesh, ccount)
	for cx := 0; cx < ccount; cx++ {
		chunks[cx] = make([]*game.ChunkMesh, ccount)
		for cz := 0; cz < ccount; cz++ {
			chunk := world.AddChunk(cx, cz)
			mesh := game.NewChunkMesh(chunk)
			scene.Add(mesh)

			chunks[cx][cz] = mesh
			fmt.Printf("(%d,%d) ", cx, cz)
		}
	}
	fmt.Println("World generation complete")

	player := game.NewPlayer(camera, func(player *game.Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Flying = true

	// game.NewPlacementGrid(chunks[0][0])

	// buffer display windows
	// uim.Attach(editor.DebugBufferWindows(app))

	// palette globals
	paletteIdx := 5
	selected := game.NewVoxel(render.DefaultPalette[paletteIdx])

	// paletteWnd := editor.PaletteWindow(render.DefaultPalette, func(newPaletteIdx int) {
	// 	paletteIdx = newPaletteIdx
	// 	selected = game.NewVoxel(render.DefaultPalette[paletteIdx])
	// })
	// paletteWnd.SetPosition(vec2.New(280, 10))
	// paletteWnd.Flow(vec2.New(200, 400))
	// uim.Attach(paletteWnd)

	// sample world position at current mouse coords
	sampleWorld := func() (vec3.T, bool) {
		depth, depthExists := geoPass.Buffer.SampleDepth(mouse.Position)
		if !depthExists {
			return vec3.Zero, false
		}
		return camera.Unproject(vec3.Extend(
			mouse.Position.Div(geoPass.Buffer.Depth.Size()),
			depth,
		)), true
	}

	// sample world normal at current mouse coords
	sampleNormal := func() (vec3.T, bool) {
		viewNormal, exists := geoPass.Buffer.SampleNormal(mouse.Position)
		if exists {
			viewInv := camera.View.Invert()
			worldNormal := viewInv.TransformDir(viewNormal)
			return worldNormal, true
		}
		return viewNormal, false
	}

	app.Draw = func(wnd *engine.Window, dt float32) {
		app.Pipeline.Draw(scene)
	}

	/* Render loop */
	app.Update = func(dt float32) {
		scene.Update(dt)

		// movement etc
		player.Update(dt)

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
