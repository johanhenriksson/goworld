package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
)

type VKRenderer struct {
	Pre      *engine.PrePass
	Geometry *GeometryPass
	Output   *OutputPass

	backend vulkan.T
	meshes  cache.Meshes
}

func NewRenderer(backend vulkan.T) engine.Renderer {
	meshes := cache.NewVkCache(backend)
	return &VKRenderer{
		Pre:      &engine.PrePass{},
		Geometry: NewGeometryPass(backend, meshes),
		Output:   NewOutputPass(backend, meshes),
		backend:  backend,
		meshes:   meshes,
	}
}

func (r *VKRenderer) Draw(args render.Args, scene object.T) {
	r.Pre.Draw(args, scene)
	r.Geometry.Draw(args, scene)
}

func (r *VKRenderer) Destroy() {
	r.backend.Device().WaitIdle()

	r.Geometry.Destroy()
	r.meshes.Destroy()
}
