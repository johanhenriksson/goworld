package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
)

type VKRenderer struct {
	Pre      *engine.PrePass
	Shadows  ShadowPass
	Geometry *GeometryPass
	Output   *OutputPass
	Lines    *LinePass

	backend  vulkan.T
	meshes   MeshCache
	textures TextureCache
}

type Pass interface {
	Draw(args render.Args, scene object.T)
	Completed() sync.Semaphore
	Destroy()
}

func NewRenderer(backend vulkan.T) engine.Renderer {
	r := &VKRenderer{
		backend:  backend,
		meshes:   NewMeshCache(backend),
		textures: NewTextureCache(backend),
	}

	r.Pre = &engine.PrePass{}
	r.Shadows = NewShadowPass(backend, r.meshes)
	r.Geometry = NewGeometryPass(backend, r.meshes, r.textures, r.Shadows)
	r.Output = NewOutputPass(backend, r.meshes, r.textures, r.Geometry)
	r.Lines = NewLinePass(backend, r.meshes, r.Output, r.Geometry)

	return r
}

func (r *VKRenderer) Draw(args render.Args, scene object.T) {
	r.Pre.Draw(args, scene)
	r.Shadows.Draw(args, scene)
	r.Geometry.Draw(args, scene)
	r.Output.Draw(args, scene)
	r.Lines.Draw(args, scene)
}

func (r *VKRenderer) Buffers() engine.BufferOutput {
	return r.Geometry.GeometryBuffer
}

func (r *VKRenderer) Destroy() {
	r.backend.Device().WaitIdle()

	r.Lines.Destroy()
	r.Shadows.Destroy()
	r.Geometry.Destroy()
	r.Output.Destroy()
	r.meshes.Destroy()
	r.textures.Destroy()
}
