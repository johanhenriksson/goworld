package vkrender

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
)

type PBRMaterial interface {
	material.T[*GeometryDescriptors]
}

type MaterialRef interface {
	cache.Item
	VertexType() any
}

type MaterialCache cache.T[MaterialRef, PBRMaterial]

// mesh cache backend
type vkmaterials struct {
	backend vulkan.T
	worker  command.Worker
}

func NewMaterialCache(backend vulkan.T) MaterialCache {
	return cache.New[MaterialRef, PBRMaterial](&vkmaterials{
		backend: backend,
		worker:  backend.Transferer(),
	})
}

func (t *vkmaterials) ItemName() string {
	return "Material"
}

func (t *vkmaterials) Instantiate(ref MaterialRef) PBRMaterial {
	return nil
}

func (m *vkmaterials) Update(mat PBRMaterial, ref MaterialRef) {
}

func (m *vkmaterials) Delete(mat PBRMaterial) {
}

func (m *vkmaterials) Destroy() {
}
