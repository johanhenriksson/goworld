package cache

import (
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Represents a mesh stored in a (sub)region of vertex/index buffers
type GpuMesh struct {
	key string

	indexType core1_0.IndexType
	vertices  buffer.Block
	indices   buffer.Block

	indexCount   int
	indexOffset  int
	vertexOffset int

	bounds shape.Sphere
}

func (m *GpuMesh) Key() string          { return m.key }
func (m *GpuMesh) Version() int         { return 1 }
func (m *GpuMesh) Bounds() shape.Sphere { return m.bounds }

func (m *GpuMesh) Bind(cmd *command.Buffer) {
	if m.indexCount <= 0 {
		// nothing to draw
		// todo: this can happen if the mesh is not ready?
		return
	}
	cmd.CmdBindVertexBuffer(m.vertices.Buffer(), 0)
	cmd.CmdBindIndexBuffers(m.indices.Buffer(), 0, m.indexType)
}

func (m *GpuMesh) Draw(cmd command.DrawIndexedBuffer, instanceOffset int) {
	m.DrawInstanced(cmd, instanceOffset, 1)
}

func (m *GpuMesh) DrawInstanced(cmd command.DrawIndexedBuffer, instanceOffset, instanceCount int) {
	if m.indexCount <= 0 {
		// nothing to draw
		return
	}
	if instanceCount <= 0 {
		// nothing to draw
		return
	}

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(command.DrawIndexed{
		IndexCount:     uint32(m.indexCount),
		InstanceCount:  uint32(instanceCount),
		IndexOffset:    uint32(m.indexOffset),
		VertexOffset:   int32(m.vertexOffset),
		InstanceOffset: uint32(instanceOffset),
	})
}
