package cache

import (
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type SamplerCache T[texture.Ref, int]

type samplers struct {
	textures *textures
	desc     *descriptor.SamplerArray
	mapping  map[int]texture.T
	next     int
}

func NewSamplerCache(backend vulkan.T, desc *descriptor.SamplerArray) SamplerCache {
	return New[texture.Ref, int](&samplers{
		textures: &textures{
			backend: backend,
			worker:  backend.Transferer(),
		},
		desc:    desc,
		mapping: make(map[int]texture.T, 100),
		next:    0,
	})
}

func (s *samplers) ItemName() string {
	return "Sampler"
}

func (s *samplers) Instantiate(ref texture.Ref) int {
	tex := s.textures.Instantiate(ref)
	id := s.next
	s.next++
	s.mapping[id] = tex
	s.desc.Set(id, tex)
	return id
}

func (s *samplers) Update(id int, ref texture.Ref) {
	tex := s.mapping[id]
	s.textures.Update(tex, ref)
}

func (s *samplers) Delete(id int) {
	tex := s.mapping[id]
	s.textures.Delete(tex)
	// todo: unset uniform
	// return id to pool?
}

func (s *samplers) Destroy() {
	s.textures.Destroy()
}
