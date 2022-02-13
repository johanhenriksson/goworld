package buffer

import (
	vk "github.com/vulkan-go/vulkan"
)

type Vertex[K any] interface {
	Buffer([]K)
}

type vbuffer[K any] struct {
	T
	binding    vk.VertexInputBindingDescription
	attributes []vk.VertexInputAttributeDescription
}

func NewVertex[K any]() Vertex[K] {

	// this is actually for the whole pipeline?
	// info := vk.PipelineVertexInputStateCreateInfo{
	// 	SType: vk.StructureTypePipelineVertexInputStateCreateInfo,
	// }

	return &vbuffer[K]{}
}

func (b *vbuffer[K]) Buffer(vertices []K) {

}

type VertexInputs interface {
}

type pipeline struct {
	// render pass
	// vertex inputs
	// descriptor sets
	// shader modules
	// input assembly state
	// rasterization
	// depth & stencil state
	// color blend attachment
	// viewport, scissor
}
