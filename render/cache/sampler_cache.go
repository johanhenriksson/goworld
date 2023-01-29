package cache

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"
)

type SamplerCache T[texture.Ref, int]

type samplers struct {
	textures *textures
	desc     *descriptor.SamplerArray
	mapping  map[int]texture.T
	next     int
}

func NewSamplerCache(dev device.T, transferer command.Worker, desc *descriptor.SamplerArray) SamplerCache {
	return New[texture.Ref, int](&samplers{
		textures: &textures{
			device: dev,
			worker: transferer,
		},
		desc:    desc,
		mapping: make(map[int]texture.T, 100),
		next:    0,
	})
}

func (s *samplers) Name() string {
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

func (s *samplers) Update(id int, ref texture.Ref) int {
	tex := s.textures.Instantiate(ref)
	s.desc.Set(id, tex)

	// deallocate old texture
	s.textures.Delete(s.mapping[id])

	s.mapping[id] = tex
	return id
}

func (s *samplers) Delete(id int) {
	tex := s.mapping[id]
	s.textures.Delete(tex)
	delete(s.mapping, id)
	// return id to pool?
}

func (s *samplers) Destroy() {
	s.textures.Destroy()
}