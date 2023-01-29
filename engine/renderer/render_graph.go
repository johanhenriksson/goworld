package renderer

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/graph"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type rgraph struct {
	target         vulkan.Target
	pool           descriptor.Pool
	deferred       pass.Deferred
	geometryPasses []pass.DeferredSubpass
	shadowPasses   []pass.DeferredSubpass
	graph          graph.T
}

func NewGraph(target vulkan.Target, geometryPasses, shadowPasses []pass.DeferredSubpass) T {
	r := &rgraph{
		target:         target,
		geometryPasses: geometryPasses,
		shadowPasses:   shadowPasses,
	}
	r.Recreate()
	return r
}

func (r *rgraph) Draw(args render.Args, scene object.T) {
	// render passes can be partially parallelized by dividing them into two parts,
	// recording and submission. queue submits must happen in order, so that semaphores
	// behave as expected. however, the actual recording of the command buffer can run
	// concurrently.
	//
	// to allow this, MeshCache and TextureCache must also be made thread safe, since
	// they currently work in a blocking manner.
	w := r.target.Worker(0)
	r.graph.Draw(w, args, scene)
}

func (r *rgraph) Recreate() {
	r.Destroy()

	r.pool = descriptor.NewPool(r.target.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 1000,
		},
		{
			Type:            vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 1000,
		},
		{
			Type:            vk.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 10000,
		},
		{
			Type:            vk.DescriptorTypeInputAttachment,
			DescriptorCount: 100,
		},
	})

	g := graph.New(r.target.Device())

	shadows := pass.NewShadowPass(r.target, r.pool, r.shadowPasses, nil)
	shadowNode := g.Node(shadows)

	deferred := pass.NewGeometryPass(r.target, r.pool, shadows, r.geometryPasses)
	deferredNode := g.Node(deferred)
	shadowNode.After(deferredNode, vk.PipelineStageTopOfPipeBit)

	forward := pass.NewForwardPass(r.target, r.pool, deferred.GBuffer(), deferred)
	forwardNode := g.Node(forward)
	forwardNode.After(deferredNode, vk.PipelineStageTopOfPipeBit)

	output := pass.NewOutputPass(r.target, r.pool, deferred.GBuffer(), deferred)
	outputNode := g.Node(output)
	outputNode.After(forwardNode, vk.PipelineStageTopOfPipeBit)

	lines := pass.NewLinePass(r.target, r.pool, deferred.GBuffer(), output)
	lineNode := g.Node(lines)
	lineNode.After(outputNode, vk.PipelineStageTopOfPipeBit)
	// preNode.Before(lineNode)

	gui := pass.NewGuiPass(r.target, r.pool, lines)
	guiNode := g.Node(gui)
	guiNode.After(lineNode, vk.PipelineStageTopOfPipeBit)
	// lineNode.Before(guiNode)

	r.deferred = deferred
	r.graph = g
}

func (r *rgraph) GBuffer() pass.GeometryBuffer {
	return r.deferred.GBuffer()
}

func (r *rgraph) Destroy() {
	r.target.Device().WaitIdle()

	if r.pool != nil {
		r.pool.Destroy()
		r.pool = nil
	}

	if r.graph != nil {
		r.graph.Destroy()
	}
	r.deferred = nil
}
