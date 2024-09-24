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
	Vertices  buffer.Block
	Indices   buffer.Block

	IndexCount   int
	IndexOffset  int
	VertexOffset int

	bounds shape.Sphere
}

func (m *GpuMesh) Key() string          { return m.key }
func (m *GpuMesh) Version() int         { return 1 }
func (m *GpuMesh) Bounds() shape.Sphere { return m.bounds }

func (m *GpuMesh) Bind(cmd *command.Buffer) {
	if m.IndexCount <= 0 {
		// nothing to draw
		// todo: this can happen if the mesh is not ready?
		return
	}
	cmd.CmdBindVertexBuffer(m.Vertices.Buffer(), 0)
	cmd.CmdBindIndexBuffers(m.Indices.Buffer(), 0, m.indexType)
}

func (m *GpuMesh) Draw(cmd command.DrawIndexedBuffer, instanceOffset int) {
	m.DrawInstanced(cmd, instanceOffset, 1)
}

func (m *GpuMesh) DrawInstanced(cmd command.DrawIndexedBuffer, instanceOffset, instanceCount int) {
	if m.IndexCount <= 0 {
		// nothing to draw
		return
	}
	if instanceCount <= 0 {
		// nothing to draw
		return
	}

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(command.DrawIndexed{
		IndexCount:     uint32(m.IndexCount),
		InstanceCount:  uint32(instanceCount),
		IndexOffset:    uint32(m.IndexOffset),
		VertexOffset:   int32(m.VertexOffset),
		InstanceOffset: uint32(instanceOffset),
	})
}
