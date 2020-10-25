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

	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	// "github.com/johanhenriksson/goworld/geometry"
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

	// create a camera
	camera := engine.CreateCamera(&render.ScreenBuffer, vec3.New(1, 22, 1), 55.0, 0.1, 600.0)
	camera.Clear = render.Color4(0.973, 0.945, 0.876, 1.0) // light gray
	camera.Rotation.X = 22
	camera.Rotation.Y = 135

	// scene & lighting setup
	scene := engine.NewScene()
	scene.Camera = camera
	scene.Lights = []engine.Light{
		{ // directional light
			Intensity:  1.2,
			Color:      vec3.New(0.9*0.973, 0.9*0.945, 0.9*0.776),
			Type:       engine.DirectionalLight,
			Projection: mat4.Orthographic(-71, 120, -20, 140, -10, 140),
			Position:   vec3.New(-2, 2, -1),
			Shadows:    false,
		},
	}

	// create chunk
	world := game.NewWorld(31481234, 16)
	chunk := world.AddChunk(0, 0)

	// first person controls
	player := game.NewPlayer(camera, func(player *game.Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Flying = true

	// create editor
	geometryPass := app.Pipeline.Get("geometry").(*engine.GeometryPass)
	edit := editor.NewEditor(chunk, camera, geometryPass.Buffer)
	scene.Add(edit)

	// buffer debug windows
	uim.Attach(editor.DebugBufferWindows(app))

	// gizmo := geometry.NewGizmo(vec3.New(3, 9, 3))
	// scene.Add(gizmo)

	// particles := engine.NewParticleSystem(vec3.New(3, 9, 3))
	// scene.Add(particles)

	// render
	app.Draw = func(wnd *engine.Window, dt float32) {
		app.Pipeline.Draw(scene)
	}

	// update loop
	app.Update = func(dt float32) {
		scene.Update(dt)

		// movement etc
		player.Update(dt)
	}

	fmt.Println("Ok")
	app.Run()
}
