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
	Geometry() pass.Deferred
	GBuffer() pass.GeometryBuffer
	Recreate()
	Screenshot()
	Destroy()
}

type rgraph struct {
	graph.T
	app      vulkan.App
	gbuffer  pass.GeometryBuffer
	geometry pass.Deferred
}

func NewGraph(app vulkan.App) T {
	renderTarget, err := pass.NewRenderTarget(app.Device(), app.Width(), app.Height(), core1_0.FormatR8G8B8A8UnsignedNormalized, app.Device().GetDepthFormat())
	if err != nil {
		panic(err)
	}

	gbuffer, err := pass.NewGbuffer(app, renderTarget)
	if err != nil {
		panic(err)
	}

	r := &rgraph{
		app:     app,
		gbuffer: gbuffer,
	}
	r.T = graph.New(app, func(g graph.T) {
		shadows := pass.NewShadowPass(app)
		shadowNode := g.Node(shadows)

		deferred := pass.NewGeometryPass(app, gbuffer, shadows)
		deferredNode := g.Node(deferred)
		deferredNode.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, gbuffer))
		forward.After(deferredNode, core1_0.PipelineStageTopOfPipe)

		gbufferCopy := g.Node(pass.NewGBufferCopyPass(gbuffer))
		gbufferCopy.After(forward, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(r.app, gbuffer))
		lines.After(forward, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(r.app, gbuffer))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		output := g.Node(pass.NewOutputPass(r.app, gbuffer))
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

func (r *rgraph) Geometry() pass.Deferred {
	return r.geometry
}

func (r *rgraph) Destroy() {
	if r.T != nil {
		r.T.Destroy()
	}
	r.gbuffer = nil
}
