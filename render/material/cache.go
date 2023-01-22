package material

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type Cache cache.T[Ref, Standard]

// mesh cache backend
type matcache struct {
	backend vulkan.T
	worker  command.Worker
}

func NewCache(backend vulkan.T) Cache {
	return cache.New[Ref, Standard](&matcache{
		backend: backend,
		worker:  backend.Transferer(),
	})
}

func (m *matcache) Name() string {
	return "Material"
}

func (m *matcache) Instantiate(ref Ref) Standard {
	return nil
}

func (m *matcache) Update(cached Standard, ref Ref) Standard {
	return cached
}

func (m *matcache) Delete(mat Standard) {
}

func (m *matcache) Destroy() {
}
