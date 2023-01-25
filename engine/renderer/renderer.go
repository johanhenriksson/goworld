package renderer

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Draw(args render.Args, scene object.T)
	GBuffer() pass.GeometryBuffer
	Recreate()
	Destroy()
}

type vkrenderer struct {
	target         vulkan.Target
	pool           descriptor.Pool
	deferred       pass.Deferred
	geometryPasses []pass.DeferredSubpass
	shadowPasses   []pass.DeferredSubpass
	passes         []pass.Pass
}

func New(target vulkan.Target, geometryPasses, shadowPasses []pass.DeferredSubpass) T {
	r := &vkrenderer{
		target:         target,
		geometryPasses: geometryPasses,
		shadowPasses:   shadowPasses,
	}
	r.Recreate()
	return r
}

func (r *vkrenderer) Draw(args render.Args, scene object.T) {
	// render passes can be partially parallelized by dividing them into two parts,
	// recording and submission. queue submits must happen in order, so that semaphores
	// behave as expected. however, the actual recording of the command buffer can run
	// concurrently.
	//
	// to allow this, MeshCache and TextureCache must also be made thread safe, since
	// they currently work in a blocking manner.
	for _, pass := range r.passes {
		pass.Draw(args, scene)
	}
}

func (r *vkrenderer) Recreate() {
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

	pre := pass.NewPrePass(r.target)
	deferred := pass.NewDeferredPass(r.target, r.pool, r.geometryPasses, r.shadowPasses)
	forward := pass.NewForwardPass(r.target, r.pool, deferred.GBuffer(), deferred)
	output := pass.NewOutputPass(r.target, r.pool, deferred.GBuffer(), forward)
	lines := pass.NewLinePass(r.target, r.pool, deferred.GBuffer(), output)
	gui := pass.NewGuiPass(r.target, r.pool, lines)
	post := pass.NewPostPass(r.target, gui)

	r.deferred = deferred
	r.passes = []pass.Pass{
		pre,
		deferred,
		forward,
		output,
		lines,
		gui,
		post,
	}
}

func (r *vkrenderer) GBuffer() pass.GeometryBuffer {
	return r.deferred.GBuffer()
}

func (r *vkrenderer) Destroy() {
	r.target.Device().WaitIdle()

	if r.pool != nil {
		r.pool.Destroy()
		r.pool = nil
	}

	for _, pass := range r.passes {
		pass.Destroy()
	}
	r.passes = nil
	r.deferred = nil
}
