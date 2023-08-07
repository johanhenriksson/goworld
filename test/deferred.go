package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/graph"
	"github.com/johanhenriksson/goworld/engine/pass"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func DeferredGraph(app vulkan.App, target vulkan.Target) graph.T {
	return graph.New(app, target, func(g graph.T, output vulkan.Target) []graph.Resource {
		size := output.Size()

		// allocate main depth buffer
		depth := vulkan.NewDepthTarget(app.Device(), "main-depth", size)

		// main off-screen color buffer
		offscreen := vulkan.NewColorTarget(app.Device(), "main-color", image.FormatRGBA8Unorm, size)

		empty := vulkan.NewColorTarget(app.Device(), "main-color", image.FormatRGBA8Unorm, vulkan.TargetSize{
			Width: 1, Height: 1, Frames: 1, Scale: 1,
		})

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), size)
		if err != nil {
			panic(err)
		}

		shadows := pass.NewShadowPass(app, output)
		shadowNode := g.Node(shadows)

		// deferred geometry
		deferredGeometry := g.Node(pass.NewDeferredGeometryPass(app, depth, gbuffer))

		// deferred lighting
		deferredLighting := g.Node(pass.NewDeferredLightingPass(app, offscreen, gbuffer, shadows, empty))
		deferredLighting.After(shadowNode, core1_0.PipelineStageTopOfPipe)
		deferredLighting.After(deferredGeometry, core1_0.PipelineStageTopOfPipe)

		outputPass := g.Node(pass.NewOutputPass(app, output, offscreen))
		outputPass.After(deferredLighting, core1_0.PipelineStageTopOfPipe)

		return []graph.Resource{
			depth,
			offscreen,
			gbuffer,
		}
	})
}

var _ = Describe("deferred renderer", Label("e2e"), func() {
	It("renders correctly", func() {
		img := engine.Frame(engine.Args{
			Width:    512,
			Height:   512,
			Title:    "goworld",
			Renderer: DeferredGraph,
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
					Size: 5,
					Mat:  material.StandardDeferred(),
				})).
					Texture(texture.Diffuse, color.White).
					Parent(scene).
					Create()

				object.Builder(cube.NewObject(cube.Args{
					Size: 1,
					Mat:  material.StandardDeferred(),
				})).
					Position(vec3.New(0, 0.5, 0)).
					Texture(texture.Diffuse, texture.Checker).
					Parent(scene).
					Create()

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
		Expect(img).To(ApproxImage("deferred.png"))
	})
})
