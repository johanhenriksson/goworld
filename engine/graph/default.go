package graph

import (
	"github.com/johanhenriksson/goworld/engine/pass"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Instantiates the default render graph
func Default(app vulkan.App, target vulkan.Target) T {
	return New(app, target, func(g T, output vulkan.Target) []Resource {
		size := output.Size()

		//
		// main render pass
		//

		// allocate main depth buffer
		depth := vulkan.NewDepthTarget(app.Device(), "main-depth", size)

		// main off-screen color buffer
		offscreen := vulkan.NewColorTarget(app.Device(), "main-color", image.FormatRGBA8Unorm, size)

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), size)
		if err != nil {
			panic(err)
		}

		shadows := pass.NewShadowPass(app, output)
		shadowNode := g.Node(shadows)

		deferred := g.Node(pass.NewDeferredPass(app, offscreen, depth, gbuffer, shadows))
		deferred.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, offscreen, depth, gbuffer))
		forward.After(deferred, core1_0.PipelineStageTopOfPipe)

		//
		// post processing
		//

		// allocate SSAO output buffer
		ssaoFormat := core1_0.FormatR16SignedFloat
		ssaoOutput := vulkan.NewColorTarget(app.Device(), "ssao-output", ssaoFormat, vulkan.TargetSize{
			Width:  size.Width / 2,
			Height: size.Height / 2,
			Frames: size.Frames,
			Scale:  size.Scale,
		})

		// create SSAO pass
		ssao := g.Node(pass.NewAmbientOcclusionPass(app, ssaoOutput, gbuffer))
		ssao.After(forward, core1_0.PipelineStageTopOfPipe)

		// SSAO blur pass
		blurOutput := vulkan.NewColorTarget(app.Device(), "blur-output", ssaoOutput.SurfaceFormat(), ssaoOutput.Size())
		blur := g.Node(pass.NewBlurPass(app, blurOutput, ssaoOutput))
		blur.After(ssao, core1_0.PipelineStageTopOfPipe)

		// post process pass
		composition := vulkan.NewColorTarget(app.Device(), "composition", offscreen.SurfaceFormat(), offscreen.Size())
		post := g.Node(pass.NewPostProcessPass(app, composition, offscreen, blurOutput))
		post.After(blur, core1_0.PipelineStageTopOfPipe)

		//
		// final image composition
		//

		lines := g.Node(pass.NewLinePass(app, composition, depth))
		lines.After(post, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, composition))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		outputPass := g.Node(pass.NewOutputPass(app, output, composition))
		outputPass.After(gui, core1_0.PipelineStageTopOfPipe)

		return []Resource{
			depth,
			offscreen,
			gbuffer,
			ssaoOutput,
			blurOutput,
			composition,
		}
	})
}
