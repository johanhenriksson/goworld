package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type samplers struct {
	textures    TextureCache
	desc        *descriptor.SamplerArray
	reverse     map[string]*SamplerHandle
	free        map[int]bool
	descriptors []texture.T
	next        int

	// the max age must be shorter than the max life of the texture cache.
	// if using a per-frame sampler cache, then the max life time should be
	// at most (texture max life) / (number of swapchain frames)
	maxAge int

	// blank keeps a reference to a blank (white) texture
	blank texture.T
}

type SamplerHandle struct {
	ID      int
	Texture texture.T
	age     int
}

type SamplerCache interface {
	Fetch(ref texture.Ref) *SamplerHandle
	TryFetch(ref texture.Ref) (*SamplerHandle, bool)

	// Writes descriptor updates to the backing Sampler Array.
	UpdateDescriptors()
}

func NewSamplerCache(textures TextureCache, desc *descriptor.SamplerArray) SamplerCache {
	samplers := &samplers{
		textures:    textures,
		desc:        desc,
		reverse:     make(map[string]*SamplerHandle, 1000),
		free:        make(map[int]bool, 100),
		descriptors: make([]texture.T, desc.Count),
		next:        0,
		maxAge:      textures.MaxAge() / 4,
		blank:       textures.Fetch(color.White),
	}

	// ensure id 0 is always blank
	samplers.assignHandle(color.White)

	return samplers
}

func (s *samplers) nextID() int {
	// check free list
	for handle := range s.free {
		delete(s.free, handle)
		return handle
	}

	// allocate new handle
	id := s.next
	if id >= s.desc.Count {
		panic("out of handles")
	}
	s.next++
	return id
}

func (s *samplers) assignHandle(ref texture.Ref) *SamplerHandle {
	if handle, exists := s.reverse[ref.Key()]; exists {
		// reset the age of the existing handle, if we have one
		handle.age = 0
		return handle
	}
	handle := &SamplerHandle{
		ID:  s.nextID(),
		age: 0,
	}
	s.reverse[ref.Key()] = handle
	return handle
}

func (s *samplers) TryFetch(ref texture.Ref) (*SamplerHandle, bool) {
	handle := s.assignHandle(ref)
	var exists bool
	if handle.Texture, exists = s.textures.TryFetch(ref); exists {
		return handle, true
	}
	return nil, false
}

func (s *samplers) Fetch(ref texture.Ref) *SamplerHandle {
	handle := s.assignHandle(ref)
	handle.Texture = s.textures.Fetch(ref)
	return handle
}

func (s *samplers) UpdateDescriptors() {
	for ref, handle := range s.reverse {
		// increase the age of the handle and check for eviction
		handle.age++
		if handle.age > s.maxAge {
			log.Println("release handle", handle.ID, "from", handle.Texture.Key())
			delete(s.reverse, ref)
			s.free[handle.ID] = true

			// overwrite descriptor with blank texture
			handle.Texture = s.blank
			// s.descriptors[handle.ID] = nil
			// s.desc.Set(handle.ID, s.blank)
			// continue
		}

		tex := handle.Texture
		if tex == nil {
			continue
		}

		if s.descriptors[handle.ID] == tex {
			// texture hasnt changed, nothing to do.
			continue
		}

		// texture has changed! update descriptor
		s.descriptors[handle.ID] = tex
		s.desc.Set(handle.ID, tex)
	}
}
