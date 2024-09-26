package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/texture"
)

const MainSubpass = renderpass.Name("main")

type BasicDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

func AssignMeshTextures(samplers cache.SamplerCache, msh mesh.Mesh, slots []texture.Slot) uniform.TextureIds {
	if len(slots) > uniform.MaxTextures {
		panic("too many textures")
	}
	textureIds := uniform.TextureIds{}
	for id, slot := range slots {
		ref := msh.Texture(slot)
		if ref != nil {
			handle, exists := samplers.TryFetch(ref)
			if exists {
				textureIds[id] = uniform.TextureId(handle.ID)
			}
		}
	}
	return textureIds
}
