package graph

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/pass"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Instantiates the default render graph
func Default(app engine.App, target engine.Target) T {
	return New(app, target, func(g T, output engine.Target) []Resource {
		size := output.Size()

		//
		// screen buffers
		//

		// allocate main depth buffer
		depth := engine.NewDepthTarget(app.Device(), "main-depth", size)

		// main off-screen color buffer
		hdrBuffer := engine.NewColorTarget(app.Device(), "main-color", core1_0.FormatR16G16B16A16SignedFloat, size)

		// create geometry buffer
		gbuffer, err := pass.NewGbuffer(app.Device(), size)
		if err != nil {
			panic(err)
		}

		// allocate SSAO output buffer
		ssaoFormat := core1_0.FormatR16SignedFloat
		ssaoOutput := engine.NewColorTarget(app.Device(), "ssao-output", ssaoFormat, engine.TargetSize{
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
		depthPass := g.Node(pass.NewDepthPass(app, depth))

		// deferred geometry
		// - wait for depth pass before fragment tests
		deferredGeometry := g.Node(pass.NewDeferredGeometryPass(app, depth, gbuffer))
		deferredGeometry.After(depthPass, core1_0.PipelineStageEarlyFragmentTests)

		// ssao pass
		// - wait for geometry before executing fragment shader
		ssao := g.Node(pass.NewAmbientOcclusionPass(app, ssaoOutput, gbuffer))
		ssao.After(deferredGeometry, core1_0.PipelineStageFragmentShader)

		// ssao blur pass
		// - wait for ssao pass before executing fragment shader
		blurOutput := engine.NewColorTarget(app.Device(), "blur-output", ssaoOutput.SurfaceFormat(), ssaoOutput.Size())
		blur := g.Node(pass.NewBlurPass(app, blurOutput, ssaoOutput))
		blur.After(ssao, core1_0.PipelineStageFragmentShader)

		// deferred lighting
		// - wait for geometry and ssao blur before executing fragment shader
		deferredLighting := g.Node(pass.NewDeferredLightingPass(app, hdrBuffer, gbuffer, shadows, blurOutput))
		deferredLighting.After(shadowNode, core1_0.PipelineStageFragmentShader)
		deferredLighting.After(blur, core1_0.PipelineStageFragmentShader)

		// forward pass
		// - wait for deferred lighting before executing fragment shader
		forward := g.Node(pass.NewForwardPass(app, hdrBuffer, depth, shadows))
		forward.After(deferredLighting, core1_0.PipelineStageFragmentShader)

		//
		// final image composition
		//

		// post process pass
		composition := engine.NewColorTarget(app.Device(), "composition", core1_0.FormatR8G8B8A8UnsignedNormalized, hdrBuffer.Size())
		post := g.Node(pass.NewPostProcessPass(app, composition, hdrBuffer))
		post.After(forward, core1_0.PipelineStageFragmentShader)

		lines := g.Node(pass.NewLinePass(app, composition, depth))
		lines.After(post, core1_0.PipelineStageFragmentShader)

		gui := g.Node(pass.NewGuiPass(app, composition))
		gui.After(lines, core1_0.PipelineStageFragmentShader)

		outputPass := g.Node(pass.NewOutputPass(app, output, composition))
		outputPass.After(gui, core1_0.PipelineStageFragmentShader)

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
