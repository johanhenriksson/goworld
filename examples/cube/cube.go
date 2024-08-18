package main

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/engine/app"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

func main() {
	app.Run(
		app.Args{
			Width:  1200,
			Height: 800,
			Title:  "goworld: cube",
		},
		func(pool object.Pool, scene object.Object) {
			rot := float32(45)
			box := cube.New(pool, cube.Args{
				Size: 1,
				Mat:  material.StandardDeferred(),
			})
			box.SetTexture(texture.Diffuse, random.Choice(color.DefaultPalette))

			object.Builder(object.Empty(pool, "Cube")).
				Position(vec3.T{Y: 0.5}).
				Attach(box).
				Attach(script.New(pool, func(scene, self object.Component, dt float32) {
					rot += dt * 360.0 / 6
					self.Parent().Transform().SetRotation(quat.Euler(0, rot, 0))
				})).
				Parent(scene).
				Create()

			ground := plane.New(pool, plane.Args{
				Size: vec2.New(10, 10),
				Mat:  material.StandardDeferred(),
			})
			ground.SetTexture(texture.Diffuse, color.White)
			object.Builder(object.Empty(pool, "Ground")).
				Attach(ground).
				Parent(scene).
				Create()

			// directional light
			object.Attach(
				scene,
				object.Builder(object.Empty(pool, "Sun")).
					Attach(light.NewDirectional(pool, light.DirectionalArgs{
						Intensity: 1.3,
						Color:     color.RGB(1, 1, 1),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(1, 2, 3)).
					Rotation(quat.Euler(45, 0, 0)).
					Create())

			object.Builder(object.Empty(pool, "Camera")).
				Rotation(quat.Euler(30, 45, 0)).
				Position(vec3.New(0, 0.5, 0)).
				Attach(
					object.Builder(object.Empty(pool, "Eye")).
						Attach(camera.New(pool, camera.Args{
							Fov:   60,
							Near:  0.1,
							Far:   100,
							Clear: color.White,
						})).
						Position(vec3.New(0, 0, -2)).
						Create(),
				).
				Parent(scene).
				Create()
		},
	)
}
