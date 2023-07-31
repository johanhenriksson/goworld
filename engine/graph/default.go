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
		depth, err := vulkan.NewDepthTarget(app.Device(), "main-depth", target.Width(), target.Height(), target.Frames(), target.Scale())

		// create off-screen render buffer
		offscreen, err := vulkan.NewColorTarget(app.Device(), "offscreen",
			output.Width(), output.Height(), output.Frames(), output.Scale(),
			image.FormatRGBA8Unorm)
		if err != nil {
			panic(err)
		}

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), output.Width(), output.Height(), output.Frames())
		if err != nil {
			panic(err)
		}

		shadows := pass.NewShadowPass(app, output)
		shadowNode := g.Node(shadows)

		deferred := g.Node(pass.NewDeferredPass(app, offscreen, depth, gbuffer, shadows))
		deferred.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, offscreen, depth, gbuffer))
		forward.After(deferred, core1_0.PipelineStageTopOfPipe)

		ssaoPass := pass.NewAmbientOcclusionPass(app, output, gbuffer)
		ssao := g.Node(ssaoPass)
		ssao.After(forward, core1_0.PipelineStageTopOfPipe)

		blurPass := pass.NewBlurPass(app, ssaoPass.Target)
		blur := g.Node(blurPass)
		blur.After(ssao, core1_0.PipelineStageTopOfPipe)

		postPass := pass.NewPostProcessPass(app, offscreen, blurPass.Target())
		post := g.Node(postPass)
		post.After(blur, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(app, postPass.Target(), depth))
		lines.After(post, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, postPass.Target()))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		outputPass := g.Node(pass.NewOutputPass(app, output, postPass.Target()))
		outputPass.After(gui, core1_0.PipelineStageTopOfPipe)

		// editor forward
		// editor lines

		return []Resource{
			depth,
			offscreen,
			gbuffer,
		}
	})
}
