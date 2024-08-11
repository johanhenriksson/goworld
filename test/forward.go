package test

import (
	. "github.com/johanhenriksson/goworld/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/app"
	"github.com/johanhenriksson/goworld/engine/graph"
	"github.com/johanhenriksson/goworld/engine/pass"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func ForwardGraph(app engine.App, target engine.Target) engine.Renderer {
	return graph.New(app, target, func(g *graph.Graph, output engine.Target) []graph.Resource {
		size := output.Size()

		// allocate main depth buffer
		depth := engine.NewDepthTarget(app.Device(), "main-depth", size)

		// main off-screen color buffer
		offscreen := engine.NewColorTarget(app.Device(), "main-color", image.FormatRGBA8Unorm, size)

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), size)
		if err != nil {
			panic(err)
		}

		// depth pre-pass
		depthPass := g.Node(pass.NewDepthPass(app, depth))

		shadows := pass.NewShadowPass(app, output)
		shadowNode := g.Node(shadows)

		forward := g.Node(pass.NewForwardPass(app, offscreen, depth, shadows))
		forward.After(depthPass, core1_0.PipelineStageEarlyFragmentTests)
		forward.After(shadowNode, core1_0.PipelineStageFragmentShader)

		outputPass := g.Node(pass.NewOutputPass(app, output, offscreen))
		outputPass.After(forward, core1_0.PipelineStageTopOfPipe)

		return []graph.Resource{
			depth,
			offscreen,
			gbuffer,
		}
	})
}

var _ = Describe("forward renderer", Label("e2e"), func() {
	It("renders correctly", func() {
		img := app.Frame(
			app.Args{
				Width:    512,
				Height:   512,
				Title:    "goworld",
				Renderer: ForwardGraph,
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
					Size: vec2.New(5, 5),
				})).
					Parent(scene).
					Texture(texture.Diffuse, color.White).
					Create()

				object.Builder(cube.NewObject(cube.Args{
					Size: 1,
					Mat:  material.StandardForward(),
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
		Expect(img).To(ApproxImage("forward.png"))
	})
})
