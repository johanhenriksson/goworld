package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Mesh interface {
	Draw(cmd command.Buffer, index int)
	DrawInstanced(cmd command.Buffer, startIndex, count int)
}

// Represents a mesh stored in dedicated vertex/index buffers
type cachedMesh struct {
	key        string
	indexCount int
	idxType    core1_0.IndexType
	vertices   buffer.T
	indices    buffer.T
}

func (m *cachedMesh) Draw(cmd command.Buffer, index int) {
	m.DrawInstanced(cmd, index, 1)
}

func (m *cachedMesh) DrawInstanced(cmd command.Buffer, firstInstance, instanceCount int) {
	if m.indexCount <= 0 {
		// nothing to draw
		return
	}
	if instanceCount <= 0 {
		// nothing to draw
		return
	}

	cmd.CmdBindVertexBuffer(m.vertices, 0)
	cmd.CmdBindIndexBuffers(m.indices, 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.indexCount, instanceCount, 0, 0, firstInstance)
}
