package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

type MaterialCache cache.T[*material.Def, []Material]

// Material implements render logic for a specific material.
type Material interface {
	ID() material.ID
	Textures() []texture.Slot
	Destroy()

	// Begin is called prior to recording, once per frame.
	// Its purpose is to clear object buffers and set up per-frame data such as cameras & lighting.
	Begin(uniform.Camera)

	// End is called after all draw groups have been processed.
	// It runs once per frame and is primarily responsible for flushing uniform buffers.
	End()

	// Bind is called just prior to recording draw calls.
	// Its called once for each group, and may be called multiple times each frame.
	// The primary use for it is to bind the material prior to drawing each group.
	Bind(command.Recorder)

	// Draw is called for each mesh in the group.
	// Its purpose is to set up per-draw data such as object transforms and textures
	// as well as issuing the draw call.
	Draw(command.Recorder, mesh.Mesh)

	// Unbind is called after recording all draw calls for the group.
	Unbind(command.Recorder)
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
