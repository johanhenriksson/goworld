package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Mesh interface {
	Draw(command.Buffer, int)
	Destroy()
}

type vkMesh struct {
	key      string
	elements int
	idxType  core1_0.IndexType
	vertices buffer.T
	indices  buffer.T
}

func (m *vkMesh) Draw(cmd command.Buffer, index int) {
	if m.elements <= 0 {
		// nothing to draw
		return
	}

	cmd.CmdBindVertexBuffer(m.vertices, 0)
	cmd.CmdBindIndexBuffers(m.indices, 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.elements, 1, 0, 0, index)
}

func (m *vkMesh) Destroy() {
	if m.vertices != nil {
		m.vertices.Destroy()
		m.vertices = nil
	}
	if m.indices != nil {
		m.indices.Destroy()
		m.indices = nil
	}
}
