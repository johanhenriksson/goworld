package renderer

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Draw(args render.Args, scene object.T)
	Deferred() pass.DeferredPass
	Recreate()
	Destroy()

	SamplePosition(cursor vec2.T) (vec3.T, bool)
	SampleNormal(cursor vec2.T) (vec3.T, bool)
}

type vkrenderer struct {
	target         vulkan.Target
	pool           descriptor.Pool
	deferred       pass.DeferredPass
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

func (r *vkrenderer) Deferred() pass.DeferredPass {
	return r.deferred
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

	pre := &pass.PrePass{}
	shadows := pass.NewShadowPass(r.target, r.pool, r.shadowPasses)
	geometry := pass.NewGeometryPass(r.target, r.pool, shadows, r.geometryPasses)
	forward := pass.NewForwardPass(r.target, r.pool, geometry.GeometryBuffer, geometry.Completed())
	output := pass.NewOutputPass(r.target, r.pool, geometry, forward.Completed())
	lines := pass.NewLinePass(r.target, r.pool, output, geometry, output.Completed())
	gui := pass.NewGuiPass(r.target, r.pool, lines)

	r.deferred = geometry
	r.passes = []pass.Pass{
		pre,
		shadows,
		geometry,
		forward,
		output,
		lines,
		gui,
	}
}

func (r *vkrenderer) SamplePosition(cursor vec2.T) (vec3.T, bool) {
	if r.deferred == nil {
		return vec3.Zero, false
	}
	return r.deferred.SamplePosition(cursor)
}

func (r *vkrenderer) SampleNormal(cursor vec2.T) (vec3.T, bool) {
	if r.deferred == nil {
		return vec3.Zero, false
	}
	return r.deferred.SampleNormal(cursor)
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
}
