package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Represents a mesh stored in a (sub)region of vertex/index buffers
type GpuMesh struct {
	Key string

	IndexType core1_0.IndexType
	Vertices  buffer.Block
	Indices   buffer.Block

	IndexCount   int
	IndexOffset  int
	VertexOffset int
}

func (m *GpuMesh) Bind(cmd *command.Buffer) {
	if m.IndexCount <= 0 {
		// nothing to draw
		// todo: this can happen if the mesh is not ready?
		return
	}
	cmd.CmdBindVertexBuffer(m.Vertices.Buffer(), 0)
	cmd.CmdBindIndexBuffers(m.Indices.Buffer(), 0, m.IndexType)
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
