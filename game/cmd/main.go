package main

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game/player"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

func main() {
	engine.Run(engine.Args{
		Width:  1600,
		Height: 1200,
		Title:  "goworld",
	},
		editor.Scene(func(scene object.Object) {

			// rot := float32(45)
			// box := cube.New(cube.Args{
			// 	Size: 1,
			// 	Mat:  material.StandardForward(),
			// })
			// box.SetTexture(texture.Diffuse, color.RGB(0.8, 0.2, 0.3))
			// object.Builder(object.Empty("Box")).
			// 	Attach(box).
			// 	Attach(script.New(func(scene, self object.Component, dt float32) {
			// 		rot += dt * 60
			// 		self.Parent().Transform().SetRotation(quat.Euler(0, rot, 0))
			// 	})).
			// 	Position(vec3.New(0, 0.5, 0)).
			// 	Parent(scene).
			// 	Create()
			//
			// ground := plane.New(plane.Args{
			// 	Size: vec2.New(10, 10),
			// 	Mat:  material.StandardForward(),
			// })
			// ground.SetTexture(texture.Diffuse, color.White)
			// object.Builder(object.Empty("Ground")).
			// 	Attach(ground).
			// 	Parent(scene).
			// 	Create()
			//
			// object.Builder(camera.NewObject(camera.Args{
			// 	Fov:   60,
			// 	Near:  0.01,
			// 	Far:   50,
			// 	Clear: color.Green,
			// })).
			// 	Position(vec3.New(0, 4, -4)).
			// 	Rotation(quat.Euler(45, 0, 0)).
			// 	Parent(scene).
			// 	Create()
			//
			// // directional light
			// object.Attach(
			// 	scene,
			// 	object.Builder(object.Empty("Sun")).
			// 		Attach(light.NewDirectional(light.DirectionalArgs{
			// 			Intensity: 100.3,
			// 			Color:     color.RGB(1, 1, 1),
			// 			Shadows:   true,
			// 			Cascades:  4,
			// 		})).
			// 		Position(vec3.New(1, 2, 3)).
			// 		Rotation(quat.Euler(rot, 0, 0)).
			// 		Attach(script.New(func(scene, self object.Component, dt float32) {
			// 			self.Parent().Transform().SetRotation(quat.Euler(30, -45, 0))
			// 		})).
			// 		Create())
			//
			// return
			// todo: should on the scene object by default
			object.Attach(scene, physics.NewWorld())

			// physics boxes
			for x := 0; x < 3; x++ {
				for z := 0; z < 3; z++ {
					box := cube.New(cube.Args{
						Size: 1,
						Mat:  material.TransparentForward(),
					})
					box.SetTexture(texture.Diffuse, random.Choice(color.DefaultPalette).WithAlpha(0.5))

					object.Builder(object.Empty("Box")).
						Position(vec3.New(20+3*float32(x), 30, 15+3*float32(z))).
						Attach(physics.NewRigidBody(5)).
						Attach(physics.NewBox(vec3.One)).
						Attach(box).
						Parent(scene).
						Create()
				}
			}

			// character
			char := player.New()
			char.Transform().SetPosition(vec3.New(5, 32, 5))
			object.Attach(scene, char)

			// terrain
			m := terrain.NewMap("default", 32)
			object.Attach(scene, terrain.NewWorld(m, 128))

			// add water
			object.Attach(scene, object.Builder(terrain.NewWater(2, 16)).
				Position(vec3.New(0, -1, 0)).
				Create())

			// directional light
			object.Attach(
				scene,
				object.Builder(object.Empty("Sun")).
					Attach(light.NewDirectional(light.DirectionalArgs{
						Intensity: 1.3,
						Color:     color.RGB(1, 1, 1),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(1, 2, 3)).
					Rotation(quat.Euler(30, -45, 0)).
					Create())

			object.Attach(
				scene,
				object.Builder(object.Empty("Light")).
					Position(vec3.New(10, 15, 10)).
					Attach(light.NewPoint(light.PointArgs{
						Range:     5,
						Color:     color.Blue,
						Intensity: 5,
					})).
					Create())
		}),
	)
}
