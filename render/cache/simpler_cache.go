package cache

import (
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type samplers struct {
	textures TextureCache
	desc     *descriptor.SamplerArray
	mapping  map[*SamplerHandle]texture.T
	reverse  map[texture.Ref]*SamplerHandle
	next     int
}

type SamplerHandle struct {
	ID      int
	Texture texture.T
}

type SamplerCache interface {
	Fetch(ref texture.Ref) *SamplerHandle
}

func NewSamplerCache(textures TextureCache, desc *descriptor.SamplerArray) SamplerCache {
	return &samplers{
		textures: textures,
		desc:     desc,
		mapping:  make(map[*SamplerHandle]texture.T, 1000),
		reverse:  make(map[texture.Ref]*SamplerHandle, 1000),
		next:     1,
	}
}

func (s *samplers) assignHandle(ref texture.Ref) *SamplerHandle {
	if handle, exists := s.reverse[ref]; exists {
		return handle
	}
	id := s.next
	if id >= s.desc.Count {
		panic("out of handles")
	}
	handle := &SamplerHandle{
		ID: id,
	}
	s.reverse[ref] = handle
	s.next++
	return handle
}

func (s *samplers) Fetch(ref texture.Ref) *SamplerHandle {
	handle := s.assignHandle(ref)
	handle.Texture = s.textures.Fetch(ref)
	if handle.Texture != nil {
		s.desc.Set(handle.ID, handle.Texture)
		return handle
	}
	return nil
}
