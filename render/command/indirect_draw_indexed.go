package command

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DrawIndexed struct {
	IndexCount     uint32
	InstanceCount  uint32
	IndexOffset    uint32
	VertexOffset   int32
	InstanceOffset uint32
}

type DrawIndexedBuffer interface {
	CmdDrawIndexed(DrawIndexed)
}

type IndirectDrawIndexedBuffer struct {
	commands    *buffer.Array[DrawIndexed]
	nextIndex   int
	batchOffset int
}

var _ DrawIndexedBuffer = (*IndirectDrawIndexedBuffer)(nil)

func NewIndirectDrawIndexedBuffer(dev *device.Device, key string, size int) *IndirectDrawIndexedBuffer {
	cmds := buffer.NewArray[DrawIndexed](dev, buffer.Args{
		Key:  key,
		Size: size,
		Usage: core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer |
			core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: device.MemoryTypeShared,
	})
	return &IndirectDrawIndexedBuffer{
		commands: cmds,
	}
}

func (i *IndirectDrawIndexedBuffer) Reset() {
	i.nextIndex = 0
	i.batchOffset = 0
}

func (i *IndirectDrawIndexedBuffer) BeginDrawIndirect() {
	i.batchOffset = i.nextIndex
}

func (i *IndirectDrawIndexedBuffer) CmdDrawIndexed(cmd DrawIndexed) {
	if cmd.IndexCount == 0 {
		return
	}
	if cmd.InstanceCount == 0 {
		return
	}
	i.commands.Set(i.nextIndex, cmd)
	i.nextIndex++
}

func (i *IndirectDrawIndexedBuffer) EndDrawIndirect(cmd *Buffer) {
	batchCount := i.nextIndex - i.batchOffset
	if batchCount == 0 {
		return
	}

	cmd.CmdDrawIndexedIndirect(i.commands, i.batchOffset*i.commands.Stride(), batchCount, i.commands.Stride())
}

func (i *IndirectDrawIndexedBuffer) Flush() {
	i.commands.Flush()
}

func (i *IndirectDrawIndexedBuffer) Destroy() {
	i.commands.Destroy()
}
