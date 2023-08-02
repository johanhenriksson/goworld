package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

var _ = Describe("forward renderer", Label("e2e"), func() {
	It("renders correctly", func() {
		img := engine.Frame(engine.Args{
			Width:  512,
			Height: 512,
			Title:  "goworld",
		},
			func(scene object.Object) {
				object.Builder(object.Empty("Camera")).
					Rotation(quat.Euler(30, 45, 0)).
					Position(vec3.New(0, 1, 0)).
					Attach(
						object.Builder(object.Empty("Eye")).
							Attach(camera.New(camera.Args{
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

				object.Builder(plane.NewObject(plane.Args{
					Size:  5,
					Color: color.White,
				})).
					Parent(scene).
					Create()

				c := object.Builder(cube.NewObject(cube.Args{
					Size: 1,
					Mat:  material.StandardForward(),
				})).
					Position(vec3.New(0, 0.5, 0)).
					Parent(scene).
					Create()
				c.Mesh.SetTexture("diffuse", texture.Checker)

				// directional light
				rot := float32(45)
				object.Attach(
					scene,
					object.Builder(object.Empty("Sun")).
						Attach(light.NewDirectional(light.DirectionalArgs{
							Intensity: 1,
							Color:     color.RGB(1, 1, 1),
							Shadows:   true,
							Cascades:  4,
						})).
						Position(vec3.New(1, 2, 3)).
						Rotation(quat.Euler(rot, 0, 0)).
						Create())
			},
		)
		Expect(img).To(ApproxImage("forward.png"))
	})
})
