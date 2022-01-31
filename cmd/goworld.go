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
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/ui"

	"github.com/lsfn/ode"
)

func main() {
	fmt.Println("goworld")

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

	// app := engine.NewApplication("goworld", 1400, 1000)
	renderer := engine.NewRenderer(wnd)

	uim := ui.NewManager(1600, 1200)
	renderer.Append("ui", uim)

	scene.Attach(light.NewDirectional(light.DirectionalArgs{
		Intensity: 1.8,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(1, -1, 1),
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

	guim := gui.New()

	for !wnd.ShouldClose() {
		scene.Update(0.030)

		renderer.Draw(scene)
		guim.DrawPass()

		wnd.SwapBuffers()
	}
}
