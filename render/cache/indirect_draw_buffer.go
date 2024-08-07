package cache

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type IndirectDrawBuffer struct {
	command.Buffer

	vertexBuffer buffer.T
	indexBuffer  buffer.T
	vertexOffset int
	indexOffset  int
	indexType    core1_0.IndexType

	drawCmds buffer.Array[DrawIndirectCommand]
	count    int
}

func NewIndirectDrawBuffer(device device.T, size int) *IndirectDrawBuffer {
	cmds := buffer.NewArray[DrawIndirectCommand](device, buffer.Args{
		Key:    "DrawIndirect",
		Size:   size,
		Usage:  core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer,
		Memory: core1_0.MemoryPropertyDeviceLocal | core1_0.MemoryPropertyHostVisible,
	})
	return &IndirectDrawBuffer{
		drawCmds: cmds,
	}
}

func (i *IndirectDrawBuffer) Flush() {
	i.drawCmds.Flush()
}

func (i *IndirectDrawBuffer) Destroy() {
	i.drawCmds.Destroy()
}

func (i *IndirectDrawBuffer) BeginDrawIndirect(cmds command.Buffer) {
	if i.Buffer != nil {
		panic("draw indirect has already begun")
	}
	i.Buffer = cmds
	i.vertexBuffer = nil
	i.indexBuffer = nil
	i.count = 0
}

func (i *IndirectDrawBuffer) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	if i.Buffer == nil {
		panic("indirect draw has not begun")
	}
	if i.vertexBuffer == nil {
		i.Buffer.CmdBindVertexBuffer(vtx, offset)
		i.vertexBuffer = vtx
		i.vertexOffset = offset
	} else if i.vertexBuffer != vtx {
		panic("vertex buffer already bound")
	}
	if i.vertexOffset != offset {
		panic("vertex offset mismatch")
	}
}

func (i *IndirectDrawBuffer) CmdBindIndexBuffers(idx buffer.T, offset int, kind core1_0.IndexType) {
	if i.Buffer == nil {
		panic("indirect draw has not begun")
	}
	if i.indexBuffer == nil {
		i.Buffer.CmdBindIndexBuffers(idx, offset, kind)
		i.indexBuffer = idx
		i.indexOffset = offset
		i.indexType = kind
	} else if i.indexBuffer != idx {
		panic("index buffer already bound")
	}
	if i.indexType != kind {
		panic("index type mismatch")
	}
	if i.indexOffset != offset {
		panic("index offset mismatch")
	}
}

func (i *IndirectDrawBuffer) CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int) {
	panic("CmdDraw is not allowed during indirect draw")
}

func (i *IndirectDrawBuffer) CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int) {
	if i.Buffer == nil {
		panic("indirect draw has not begun")
	}
	cmd := DrawIndirectCommand{
		IndexCount:    uint32(indexCount),
		InstanceCount: uint32(instanceCount),
		FirstIndex:    uint32(firstIndex),
		VertexOffset:  int32(vertexOffset),
		FirstInstance: uint32(firstInstance),
	}
	i.drawCmds.Write(i.count, &cmd)
	i.count++
}

func (i *IndirectDrawBuffer) EndDrawIndirect() {
	if i.Buffer == nil {
		panic("indirect draw has not begun")
	}
	if i.vertexBuffer == nil {
		panic("no vertex buffer bound")
	}
	if i.indexBuffer == nil {
		panic("no index buffer bound")
	}
	i.drawCmds.Flush()
	stride := int(unsafe.Sizeof(DrawIndirectCommand{}))
	i.Buffer.Ptr().CmdDrawIndexedIndirect(i.drawCmds.Ptr(), 0, i.count, stride)
	i.Buffer = nil
}

type DrawIndirectCommand struct {
	IndexCount    uint32
	InstanceCount uint32
	FirstIndex    uint32
	VertexOffset  int32
	FirstInstance uint32
}
