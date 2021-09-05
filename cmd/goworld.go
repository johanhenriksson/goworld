package main

/*
 * Copyright (C) 2016-2021 Johan Henriksson
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

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/lsfn/ode"

	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

func main() {
	fmt.Println("goworld")

	w := ode.NewWorld()
	w.SetGravity(ode.Vector3{0, -9.82, 0})

	scene := scene.New()

	wnd, err := window.New(window.Args{
		Title:        "goworld2",
		Width:        1600,
		Height:       900,
		InputHandler: scene,
	})
	if err != nil {
		panic(err)
	}

	fwidth, fheight := wnd.BufferSize()
	render.ScreenBuffer.Width, render.ScreenBuffer.Height = fwidth, fheight
	fmt.Printf("buffer %+v\n", render.ScreenBuffer)
	aspect := float32(fwidth) / float32(fheight)

	// app := engine.NewApplication("goworld", 1400, 1000)
	renderer := engine.NewRenderer()

	uim := ui.NewManager(1600, 900)
	renderer.Append("ui", uim)

	// create a cam
	// cam := engine.CreateCamera(&render.ScreenBuffer, vec3.New(1, 22, 1), 55.0, 0.1, 600.0)
	// cam.Clear = render.Hex("#eddaab")
	cam := camera.New(aspect, 55.0, 0.1, 600)
	scene.SetCamera(cam)

	gizmo := mover.New(mover.Args{})
	gizmo.Transform().SetPosition(vec3.New(-1, 0, -1))
	scene.Adopt(gizmo)

	// create chunk
	world := game.NewWorld(31481234, 16)
	chunk := world.AddChunk(0, 0)

	// first person controls
	player := game.NewPlayer(vec3.New(1, 22, 1), cam, func(player *game.Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Flying = true
	player.Eye.Transform().SetRotation(vec3.New(22, 135, 0))
	scene.Adopt(player)

	// create editor
	edit := editor.NewEditor(chunk, cam, renderer.Geometry.Buffer)
	scene.Adopt(edit)
	uim.Attach(edit.Palette)

	// buffer debug windows
	uim.Attach(editor.DebugBufferWindows(renderer))

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
