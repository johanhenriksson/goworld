package cache

import (
	"time"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type SamplerCache struct {
	textures    TextureCache
	reverse     map[string]*SamplerHandle
	free        map[int]bool
	descriptors texture.Array
	next        int
	size        int

	// the max age must be shorter than the max life of the texture cache.
	maxAge   time.Duration
	lastTick time.Time

	// blank keeps a reference to a blank (white) texture
	blank *texture.Texture
}

type SamplerHandle struct {
	ID      int
	Texture *texture.Texture
	age     time.Duration
}

func NewSamplerCache(textures TextureCache, size int) *SamplerCache {
	samplers := &SamplerCache{
		textures:    textures,
		reverse:     make(map[string]*SamplerHandle, size),
		free:        make(map[int]bool, size),
		descriptors: make(texture.Array, size),
		next:        0,
		size:        size,
		maxAge:      textures.MaxAge(),
		lastTick:    time.Now(),
		blank:       textures.Fetch(color.White),
	}

	// ensure id 0 is always blank/white
	samplers.assignHandle(color.White)

	// initialize descriptors with blank texture
	// todo: maybe use an error texture
	for i := range samplers.descriptors {
		samplers.descriptors[i] = samplers.blank
	}

	return samplers
}

func (s *SamplerCache) MaxAge() time.Duration {
	return s.maxAge
}

func (s *SamplerCache) Size() int {
	return s.size
}

func (s *SamplerCache) nextID() int {
	// check free list
	for handle := range s.free {
		delete(s.free, handle)
		return handle
	}

	// allocate new handle
	id := s.next
	if id >= s.size {
		panic("out of handles")
	}
	s.next++
	return id
}

type Keyed interface {
	Key() string
}

func (s *SamplerCache) assignHandle(ref Keyed) *SamplerHandle {
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

func (s *SamplerCache) TryFetch(ref assets.Texture) (*SamplerHandle, bool) {
	handle := s.assignHandle(ref)
	var exists bool
	handle.Texture, exists = s.textures.TryFetch(ref)
	if !exists {
		return nil, false
	}

	s.descriptors[handle.ID] = handle.Texture
	return handle, true
}

func (s *SamplerCache) Fetch(ref assets.Texture) *SamplerHandle {
	handle := s.assignHandle(ref)
	handle.Texture = s.textures.Fetch(ref)
	s.descriptors[handle.ID] = handle.Texture
	return handle
}

func (s *SamplerCache) Assign(tex *texture.Texture) *SamplerHandle {
	handle := s.assignHandle(tex)
	handle.Texture = tex
	s.descriptors[handle.ID] = tex
	return handle
}

func (s *SamplerCache) Flush(samplers *descriptor.SamplerArray) {
	s.Tick()

	samplers.SetRange(0, s.descriptors[:s.next])
}

func (s *SamplerCache) Tick() {
	now := time.Now()
	delta := now.Sub(s.lastTick)
	s.lastTick = now

	for ref, handle := range s.reverse {
		if handle.ID == 0 {
			// never evict the blank texture
			continue
		}

		// increase the age of the handle and check for eviction
		handle.age += delta
		if handle.age > s.maxAge {
			delete(s.reverse, ref)
			s.free[handle.ID] = true

			// overwrite descriptor with blank texture
			handle.Texture = s.blank
			s.descriptors[handle.ID] = s.blank
		}
	}
}
