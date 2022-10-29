package material

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type Cache cache.T[Ref, Standard]

type VkMesh interface {
	Draw(command.Buffer, int)
	Destroy()
}

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

func (m *matcache) ItemName() string {
	return "Material"
}

func (m *matcache) Instantiate(ref Ref) Standard {
	return nil
}

func (m *matcache) Update(cached Standard, ref Ref) {
}

func (m *matcache) Delete(mat Standard) {
}

func (m *matcache) Destroy() {
}
