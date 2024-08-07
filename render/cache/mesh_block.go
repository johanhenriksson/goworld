package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Represents a mesh stored in a sub-region of shared vertex/index buffers
type meshBlock struct {
	idxType  core1_0.IndexType
	vertices buffer.Block
	indices  buffer.Block

	indexCount   int
	firstIndex   int
	vertexOffset int
}

func (m *meshBlock) Draw(cmd command.Buffer, index int) {
	m.DrawInstanced(cmd, index, 1)
}

func (m *meshBlock) DrawInstanced(cmd command.Buffer, firstInstance, instanceCount int) {
	if m.indexCount <= 0 {
		// nothing to draw
		return
	}
	if instanceCount <= 0 {
		// nothing to draw
		return
	}

	cmd.CmdBindVertexBuffer(m.vertices.Buffer(), 0)
	cmd.CmdBindIndexBuffers(m.indices.Buffer(), 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.indexCount, instanceCount, m.firstIndex, m.vertexOffset, firstInstance)
}
