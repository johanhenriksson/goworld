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
	target  vulkan.Target
	gbuffer pass.GeometryBuffer
}

func NewGraph(target vulkan.Target) T {
	r := &rgraph{
		target:  target,
		gbuffer: nil,
	}
	r.T = graph.New(target, func(g graph.T) {
		shadows := pass.NewShadowPass(target)
		shadowNode := g.Node(shadows)

		deferred := pass.NewGeometryPass(target, shadows)
		deferredNode := g.Node(deferred)
		deferredNode.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		// store gbuffer reference
		r.gbuffer = deferred.GBuffer()

		forward := g.Node(pass.NewForwardPass(target, r.gbuffer))
		forward.After(deferredNode, core1_0.PipelineStageTopOfPipe)

		gbufferCopy := g.Node(pass.NewGBufferCopyPass(r.gbuffer))
		gbufferCopy.After(forward, core1_0.PipelineStageTopOfPipe)

		output := g.Node(pass.NewOutputPass(r.target, r.gbuffer))
		output.After(forward, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(r.target, r.gbuffer))
		lines.After(output, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(r.target))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)
	})
	return r
}

func (r *rgraph) Screenshot() {
	idx := 0
	ss, err := upload.DownloadImage(r.target.Device(), r.target.Worker(idx), r.target.Surfaces()[idx])
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
