package cache

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type samplers struct {
	textures TextureCache
	desc     *descriptor.SamplerArray
	reverse  map[texture.Ref]*SamplerHandle
	next     int
	blank    texture.T
}

type SamplerHandle struct {
	ID      int
	Texture texture.T
}

type SamplerCache interface {
	Fetch(ref texture.Ref) *SamplerHandle
	UpdateDescriptors()
}

func NewSamplerCache(textures TextureCache, desc *descriptor.SamplerArray) SamplerCache {
	samplers := &samplers{
		textures: textures,
		desc:     desc,
		reverse:  make(map[texture.Ref]*SamplerHandle, 1000),
		next:     0,
		blank:    textures.FetchSync(color.White),
	}

	// ensure id 0 is always blank
	samplers.assignHandle(color.White)

	return samplers
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
	var exists bool
	handle.Texture, exists = s.textures.Fetch(ref)
	if exists {
		return handle
	}
	return nil
}

func (s *samplers) UpdateDescriptors() {
	textures := make([]texture.T, s.desc.Count)
	for i := range textures {
		textures[i] = s.blank
	}
	for _, handle := range s.reverse {
		if handle.Texture != nil {
			textures[handle.ID] = handle.Texture
		}
	}
	s.desc.SetRange(textures, 0)
}
