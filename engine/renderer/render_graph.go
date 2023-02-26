package renderer

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/graph"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render/upload"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type T interface {
	Draw(scene object.T)
	GBuffer() pass.GeometryBuffer
	Recreate()
	Screenshot()
	Destroy()
}

type rgraph struct {
	graph.T
	app     vulkan.App
	gbuffer pass.GeometryBuffer
	target  pass.RenderTarget
}

func NewGraph(app vulkan.App) T {
	r := &rgraph{
		app: app,
	}
	r.T = graph.New(app, func(g graph.T) {
		if r.target != nil {
			r.target.Destroy()
			r.target = nil
		}
		if r.gbuffer != nil {
			r.gbuffer.Destroy()
			r.gbuffer = nil
		}
		// this is an awkward place to instantiate render target / gbuffer
		// todo: maybe move it into the graph?
		// putting it here does not even let us destroy it properly
		var err error
		r.target, err = pass.NewRenderTarget(app.Device(), app.Width(), app.Height(), core1_0.FormatR8G8B8A8UnsignedNormalized, app.Device().GetDepthFormat())
		if err != nil {
			panic(err)
		}

		r.gbuffer, err = pass.NewGbuffer(app.Device(), app.Width(), app.Height())
		if err != nil {
			panic(err)
		}

		shadows := pass.NewShadowPass(app)
		shadowNode := g.Node(shadows)

		deferred := g.Node(pass.NewDeferredPass(app, r.target, r.gbuffer, shadows))
		deferred.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, r.target, r.gbuffer))
		forward.After(deferred, core1_0.PipelineStageTopOfPipe)

		// at this point we are done writing to the gbuffer, so we may copy it.
		gbufferCopy := g.Node(pass.NewGBufferCopyPass(r.gbuffer))
		gbufferCopy.After(forward, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(app, r.target))
		lines.After(forward, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, r.target))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		output := g.Node(pass.NewOutputPass(app, r.target))
		output.After(gui, core1_0.PipelineStageTopOfPipe)

		// editor forward
		// editor lines
	})
	return r
}

func (r *rgraph) Screenshot() {
	idx := 0
	ss, err := upload.DownloadImage(r.app.Device(), r.app.Worker(idx), r.app.Surfaces()[idx])
	if err != nil {
		panic(err)
	}
	filename := fmt.Sprintf("Screenshot-%s.png", time.Now().Format("2006-01-02_15-04-05"))
	if err := upload.SavePng(ss, filename); err != nil {
		panic(err)
	}
}

func (r *rgraph) GBuffer() pass.GeometryBuffer {
	return r.gbuffer
}

func (r *rgraph) Destroy() {
	if r.T != nil {
		r.T.Destroy()
	}
	r.gbuffer = nil
}
