package cache

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"
)

type SamplerCache T[texture.Ref, int]

type samplers struct {
	textures Backend[texture.Ref, texture.T]
	desc     *descriptor.SamplerArray
	mapping  map[int]texture.T
	next     int
}

func NewSamplerCache(device device.T, worker command.Worker, desc *descriptor.SamplerArray) SamplerCache {
	return New[texture.Ref, int](&samplers{
		textures: &textures{
			device: device,
			worker: worker,
		},
		desc:    desc,
		mapping: make(map[int]texture.T, 100),
		next:    0,
	})
}

func (s *samplers) Name() string {
	return "Sampler"
}

func (s *samplers) Instantiate(ref texture.Ref, callback func(int)) {
	s.textures.Instantiate(ref, func(tex texture.T) {
		id := s.next
		s.next++
		s.mapping[id] = tex
		s.desc.Set(id, tex)
		callback(id)
	})
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
