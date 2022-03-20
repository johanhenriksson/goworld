package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
)

type VKRenderer struct {
	Pre      *engine.PrePass
	Shadows  ShadowPass
	Geometry *GeometryPass
	Output   *OutputPass

	backend  vulkan.T
	meshes   cache.Meshes
	textures cache.Textures
}

type Pass interface {
	Draw(args render.Args, scene object.T)
	Completed() sync.Semaphore
	Destroy()
}

func NewRenderer(backend vulkan.T) engine.Renderer {
	r := &VKRenderer{
		backend:  backend,
		meshes:   cache.NewVkCache(backend),
		textures: cache.NewVkTextures(backend),
	}

	r.Pre = &engine.PrePass{}
	r.Shadows = NewShadowPass(backend, r.meshes)
	r.Geometry = NewGeometryPass(backend, r.meshes, r.Shadows)
	r.Output = NewOutputPass(backend, r.meshes, r.textures, r.Geometry, r.Shadows)

	return r
}

func (r *VKRenderer) Draw(args render.Args, scene object.T) {
	r.Pre.Draw(args, scene)
	r.Shadows.Draw(args, scene)
	r.Geometry.Draw(args, scene)
	r.Output.Draw(args, scene)
}

func (r *VKRenderer) Destroy() {
	r.backend.Device().WaitIdle()

	r.Shadows.Destroy()
	r.Geometry.Destroy()
	r.Output.Destroy()
	r.meshes.Destroy()
	r.textures.Destroy()
}
