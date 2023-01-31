package renderer

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/graph"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/upload"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Draw(args render.Args, scene object.T)
	GBuffer() pass.GeometryBuffer
	Recreate()
	Screenshot()
	Destroy()
}

type rgraph struct {
	target  vulkan.Target
	gbuffer pass.GeometryBuffer
	graph   graph.T
}

func NewGraph(target vulkan.Target) T {
	r := &rgraph{
		target: target,
	}
	r.Recreate()
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

func (r *rgraph) Draw(args render.Args, scene object.T) {
	// render passes can be partially parallelized by dividing them into two parts,
	// recording and submission. queue submits must happen in order, so that semaphores
	// behave as expected. however, the actual recording of the command buffer can run
	// concurrently.
	//
	// to allow this, MeshCache and TextureCache must also be made thread safe, since
	// they currently work in a blocking manner.
	w := r.target.Worker(args.Context.Index)
	r.graph.Draw(w, args, scene)
}

func (r *rgraph) Recreate() {
	r.Destroy()
	// realloc descriptor pool
	r.target.Pool().Recreate()

	g := graph.New(r.target.Device())

	shadows := pass.NewShadowPass(r.target)
	shadowNode := g.Node(shadows)

	deferred := pass.NewGeometryPass(r.target, shadows)
	deferredNode := g.Node(deferred)
	deferredNode.After(shadowNode, vk.PipelineStageTopOfPipeBit)

	r.gbuffer = deferred.GBuffer()

	forward := g.Node(pass.NewForwardPass(r.target, r.gbuffer))
	forward.After(deferredNode, vk.PipelineStageTopOfPipeBit)

	gbufferCopy := g.Node(pass.NewGBufferCopyPass(r.gbuffer))
	gbufferCopy.After(forward, vk.PipelineStageTopOfPipeBit)

	output := g.Node(pass.NewOutputPass(r.target, r.gbuffer))
	output.After(forward, vk.PipelineStageTopOfPipeBit)

	lines := g.Node(pass.NewLinePass(r.target, r.gbuffer))
	lines.After(output, vk.PipelineStageTopOfPipeBit)

	gui := g.Node(pass.NewGuiPass(r.target))
	gui.After(lines, vk.PipelineStageTopOfPipeBit)

	r.graph = g
}

func (r *rgraph) GBuffer() pass.GeometryBuffer {
	return r.gbuffer
}

func (r *rgraph) Destroy() {
	r.target.Device().WaitIdle()

	if r.graph != nil {
		r.graph.Destroy()
	}
	r.gbuffer = nil
}
