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

func (m *matcache) Instantiate(ref Ref, callback func(Standard)) {
}

func (m *matcache) Delete(mat Standard) {
}

func (m *matcache) Submit() {
}

func (m *matcache) Destroy() {
}
