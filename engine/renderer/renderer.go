package renderer

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type T interface {
	Draw(args render.Args, scene object.T)
	Buffers() pass.BufferOutput
	Destroy()
}

type vkrenderer struct {
	Pre      *pass.PrePass
	Shadows  pass.ShadowPass
	Geometry *pass.GeometryPass
	Output   *pass.OutputPass
	Lines    *pass.LinePass
	GUI      *pass.GuiPass

	backend  vulkan.T
	meshes   cache.MeshCache
	textures cache.TextureCache
}

func New(backend vulkan.T, geometryPasses, shadowPasses []pass.DeferredSubpass) T {
	r := &vkrenderer{
		backend:  backend,
		meshes:   cache.NewMeshCache(backend),
		textures: cache.NewTextureCache(backend),
	}

	r.Pre = &pass.PrePass{}
	r.Shadows = pass.NewShadowPass(backend, r.meshes, shadowPasses)
	r.Geometry = pass.NewGeometryPass(backend, r.meshes, r.textures, r.Shadows, geometryPasses)
	r.Output = pass.NewOutputPass(backend, r.meshes, r.textures, r.Geometry)
	r.Lines = pass.NewLinePass(backend, r.meshes, r.Output, r.Geometry)
	r.GUI = pass.NewGuiPass(backend, r.Lines, r.meshes)

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

	r.meshes.Tick()
	r.textures.Tick()

	r.Pre.Draw(args, scene)
	r.Shadows.Draw(args, scene)
	r.Geometry.Draw(args, scene)
	r.Output.Draw(args, scene)
	r.Lines.Draw(args, scene)
	r.GUI.Draw(args, scene)
}

func (r *vkrenderer) Buffers() pass.BufferOutput {
	return r.Geometry.GeometryBuffer
}

func (r *vkrenderer) Destroy() {
	r.backend.Device().WaitIdle()

	r.GUI.Destroy()
	r.Lines.Destroy()
	r.Shadows.Destroy()
	r.Geometry.Destroy()
	r.Output.Destroy()
	r.meshes.Destroy()
	r.textures.Destroy()
}
