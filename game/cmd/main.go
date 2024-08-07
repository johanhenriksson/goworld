package main

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/script"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/player"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/geometry/cone"
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
			// todo: should on the scene object by default
			object.Attach(scene, physics.NewWorld())

			generator := chunk.ExampleWorldgen(4, 16)
			// chonks := chunk.NewWorld(16, generator, 128)
			// object.Attach(scene, chonks)

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

			object.Builder(object.Empty("Box")).
				Position(vec3.New(14, 14, 14)).
				Attach(physics.NewRigidBody(0)).
				Attach(chunk.NewMesh(chunk.Generate(generator, 4, 0, 0))).
				Attach(physics.NewMesh()).
				Parent(scene).
				Create()

			object.Builder(object.Empty("ParentyParent")).
				Position(vec3.New(10, 15, 10)).
				Attach(physics.NewRigidBody(1)).
				Attach(physics.NewCompound()).
				Attach(physics.NewSphere(1)).
				Attach(
					object.Builder(cone.NewObject(cone.Args{
						Segments: 4,
						Height:   1,
						Radius:   0.25,
						Color:    color.White,
					})).
						Position(vec3.New(-2, 0, 0)).
						Create(),
				).
				Attach(
					object.Builder(cone.NewObject(cone.Args{
						Segments: 4,
						Height:   1,
						Radius:   0.25,
						Color:    color.White,
					})).
						Position(vec3.New(2, 0, 0)).
						Create(),
				).
				Parent(scene).
				Create()

			// character
			char := player.New()
			char.Transform().SetPosition(vec3.New(5, 32, 5))
			object.Attach(scene, char)

			// terrain
			m := terrain.NewMap("default", 32)
			object.Attach(scene, terrain.NewWorld(m, 200))
			object.Attach(scene, terrain.NewWorld(m, 1000))

			// add water
			object.Attach(scene, object.Builder(terrain.NewWater(512, 3000)).
				Position(vec3.New(0, -1, 0)).
				Create())

			// directional light
			rot := float32(45)
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
					Rotation(quat.Euler(rot, 0, 0)).
					Attach(script.New(func(scene, self object.Component, dt float32) {
						rot -= dt * 360.0 / 86400.0
						self.Parent().Transform().SetRotation(quat.Euler(rot, 0, 0))
					})).
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
