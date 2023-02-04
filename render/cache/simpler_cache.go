package cache

import (
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type samplers struct {
	textures TextureCache
	desc     *descriptor.SamplerArray
	mapping  map[int]texture.T
	reverse  map[texture.Ref]int
	next     int
}

type SamplerCache interface {
	Fetch(ref texture.Ref) int
}

func NewSamplerCache(textures TextureCache, desc *descriptor.SamplerArray) SamplerCache {
	return &samplers{
		textures: textures,
		desc:     desc,
		mapping:  make(map[int]texture.T, 1000),
		reverse:  make(map[texture.Ref]int, 1000),
		next:     1,
	}
}

func (s *samplers) assignHandle(ref texture.Ref) int {
	if handle, exists := s.reverse[ref]; exists {
		return handle
	}
	handle := s.next
	if handle >= s.desc.Count {
		panic("out of handles")
	}
	s.reverse[ref] = handle
	s.next++
	return handle
}

func (s *samplers) Fetch(ref texture.Ref) int {
	handle := s.assignHandle(ref)
	tex := s.textures.Fetch(ref)
	if tex != nil {
		s.desc.Set(handle, tex)
		return handle
	}
	return 0
}
