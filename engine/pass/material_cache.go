package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type MaterialCache[T any] struct {
	cache  map[material.ID][]T
	app    vulkan.App
	frames int
	maker  MaterialMaker[T]
}

func NewMaterialCache[T any](app vulkan.App, frames int, maker MaterialMaker[T]) *MaterialCache[T] {
	ms := &MaterialCache[T]{
		app:    app,
		frames: frames,
		cache:  map[material.ID][]T{},
		maker:  maker,
	}
	return ms
}

func (m *MaterialCache[T]) Get(msh mesh.Mesh, frame int) (T, bool) {
	matId := msh.MaterialID()
	mat, exists := m.cache[matId]
	if !exists {
		// initialize material
		var ready bool
		mat, ready = m.Load(msh.Material())
		if !ready {
			// not ready yet
			var empty T
			return empty, false
		}
	}
	return mat[frame], true
}

func (m *MaterialCache[T]) Exists(matId material.ID) bool {
	_, exists := m.cache[matId]
	return exists
}

func (m *MaterialCache[T]) Destroy() {
	for _, mat := range m.cache {
		m.maker.Destroy(mat[0])
	}
	m.cache = nil
}

func (m *MaterialCache[T]) Load(def *material.Def) ([]T, bool) {

	mat := m.maker.Instantiate(def, m.frames)
	if len(mat) == 0 {
		// not ready yet
		return nil, false
	}

	id := material.Hash(def)
	m.cache[id] = mat
	return mat, true
}

type MaterialMaker[T any] interface {
	Instantiate(mat *material.Def, count int) []T
	Destroy(T)
	// Prepare(func(*MeshGroup[T]))
	Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[T], lights []light.T)
}

func AssignMeshTextures(samplers cache.SamplerCache, msh mesh.Mesh, slots []texture.Slot) [4]uint32 {
	textureIds := [4]uint32{}
	for id, slot := range slots {
		ref := msh.Texture(slot)
		if ref != nil {
			handle, exists := samplers.TryFetch(ref)
			if exists {
				textureIds[id] = uint32(handle.ID)
			}
		}
	}
	return textureIds
}
