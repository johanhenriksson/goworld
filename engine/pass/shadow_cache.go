package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/cache"
)

type ShadowCache struct {
	samplers cache.SamplerCache
	lookup   ShadowmapLookupFn
	shared   bool
}

var _ light.ShadowmapStore = &ShadowCache{}

func NewShadowCache(samplers cache.SamplerCache, lookup ShadowmapLookupFn) *ShadowCache {
	return &ShadowCache{
		samplers: samplers,
		lookup:   lookup,
		shared:   true,
	}
}

func (s *ShadowCache) Lookup(lit light.T, cascade int) (int, bool) {
	if shadowtex := s.lookup(lit, cascade); shadowtex != nil {
		handle := s.samplers.Assign(shadowtex)
		return handle.ID, true
	}
	// no shadowmap available
	return 0, false
}

// Flush the underlying sampler cache
func (s *ShadowCache) Flush() {
	s.samplers.Flush()
}
