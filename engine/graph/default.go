package graph

import (
	"github.com/johanhenriksson/goworld/engine/pass"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Instantiates the default render graph
func Default(app vulkan.App, target vulkan.Target) T {
	return New(app, target, func(g T, output vulkan.Target) []Resource {
		size := output.Size()

		//
		// screen buffers
		//

		// allocate main depth buffer
		depth := vulkan.NewDepthTarget(app.Device(), "main-depth", size)

		// main off-screen color buffer
		hdrBuffer := vulkan.NewColorTarget(app.Device(), "main-color", core1_0.FormatR16G16B16A16SignedFloat, size)

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), size)
		if err != nil {
			panic(err)
		}

		// allocate SSAO output buffer
		ssaoFormat := core1_0.FormatR16SignedFloat
		ssaoOutput := vulkan.NewColorTarget(app.Device(), "ssao-output", ssaoFormat, vulkan.TargetSize{
			Width:  size.Width / 2,
			Height: size.Height / 2,
			Frames: size.Frames,
			Scale:  size.Scale,
		})

		//
		// main render pass
		//

		shadows := pass.NewShadowPass(app, output)
		shadowNode := g.Node(shadows)

		// depth pre-pass
		depthPass := g.Node(pass.NewDepthPass(app, depth, gbuffer))

		// deferred geometry
		deferredGeometry := g.Node(pass.NewDeferredGeometryPass(app, depth, gbuffer))
		deferredGeometry.After(depthPass, core1_0.PipelineStageTopOfPipe)

		// ssao pass
		ssao := g.Node(pass.NewAmbientOcclusionPass(app, ssaoOutput, gbuffer))
		ssao.After(deferredGeometry, core1_0.PipelineStageTopOfPipe)

		// ssao blur pass
		blurOutput := vulkan.NewColorTarget(app.Device(), "blur-output", ssaoOutput.SurfaceFormat(), ssaoOutput.Size())
		blur := g.Node(pass.NewBlurPass(app, blurOutput, ssaoOutput))
		blur.After(ssao, core1_0.PipelineStageTopOfPipe)

		// deferred lighting
		deferredLighting := g.Node(pass.NewDeferredLightingPass(app, hdrBuffer, gbuffer, shadows, blurOutput))
		deferredLighting.After(shadowNode, core1_0.PipelineStageTopOfPipe)
		deferredLighting.After(blur, core1_0.PipelineStageTopOfPipe)

		// forward pass
		forward := g.Node(pass.NewForwardPass(app, hdrBuffer, depth, shadows))
		forward.After(deferredLighting, core1_0.PipelineStageTopOfPipe)

		//
		// final image composition
		//

		// post process pass
		composition := vulkan.NewColorTarget(app.Device(), "composition", hdrBuffer.SurfaceFormat(), hdrBuffer.Size())
		post := g.Node(pass.NewPostProcessPass(app, composition, hdrBuffer))
		post.After(forward, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(app, composition, depth))
		lines.After(post, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, composition))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		outputPass := g.Node(pass.NewOutputPass(app, output, composition))
		outputPass.After(gui, core1_0.PipelineStageTopOfPipe)

		return []Resource{
			depth,
			hdrBuffer,
			gbuffer,
			ssaoOutput,
			blurOutput,
			composition,
		}
	})
}
