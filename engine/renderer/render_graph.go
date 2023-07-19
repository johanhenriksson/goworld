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
	Draw(scene object.Object, time, delta float32)
	GBuffer() pass.GeometryBuffer
	Recreate()
	Screenshot()
	Destroy()
}

type rgraph struct {
	graph.T
	app vulkan.App

	// temporary reference
	// get rid of it once the editor no longer needs it
	gbuffer pass.GeometryBuffer
}

func NewGraph(app vulkan.App) T {
	r := &rgraph{
		app: app,
	}
	r.T = graph.New(app, func(g graph.T, target pass.RenderTarget, gbuffer pass.GeometryBuffer) {
		r.gbuffer = gbuffer

		shadows := pass.NewShadowPass(app)
		shadowNode := g.Node(shadows)

		deferred := g.Node(pass.NewDeferredPass(app, target, gbuffer, shadows))
		deferred.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, target, gbuffer))
		forward.After(deferred, core1_0.PipelineStageTopOfPipe)

		// at this point we are done writing to the gbuffer, so we may copy it.
		gbufferCopy := g.Node(pass.NewGBufferCopyPass(gbuffer))
		gbufferCopy.After(forward, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(app, target))
		lines.After(forward, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, target))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		output := g.Node(pass.NewOutputPass(app, target))
		output.After(gui, core1_0.PipelineStageTopOfPipe)

		// editor forward
		// editor lines
	})
	return r
}

func (r *rgraph) Screenshot() {
	idx := 0
	r.app.Device().WaitIdle()
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
