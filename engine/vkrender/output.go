package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type OutputPass struct {
	backend vulkan.T
	meshes  cache.Meshes
	quad    cache.GpuMesh
}

func NewOutputPass(backend vulkan.T, meshes cache.Meshes) *OutputPass {
	p := &OutputPass{
		backend: backend,
		meshes:  meshes,
	}

	quadvtx := vertex.NewTriangles("screen_quad", []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})

	p.quad = p.meshes.Fetch(quadvtx, nil)

	return p
}

func (p *OutputPass) Draw(args render.Args, scene object.T) {
}

func (p *OutputPass) Destroy() {

}
