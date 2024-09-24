package command

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type IndirectIndexedDrawBuffer struct {
	commands    *buffer.Array[DrawIndexed]
	nextIndex   int
	batchOffset int
}

var _ DrawIndexedBuffer = (*IndirectIndexedDrawBuffer)(nil)

func NewIndirectIndexedDrawBuffer(dev *device.Device, key string, size int) *IndirectIndexedDrawBuffer {
	cmds := buffer.NewArray[DrawIndexed](dev, buffer.Args{
		Key:  key,
		Size: size,
		Usage: core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer |
			core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: device.MemoryTypeShared,
	})
	return &IndirectIndexedDrawBuffer{
		commands: cmds,
	}
}

func (i *IndirectIndexedDrawBuffer) Reset() {
	i.nextIndex = 0
	i.batchOffset = 0
}

func (i *IndirectIndexedDrawBuffer) BeginDrawIndirect() {
	i.batchOffset = i.nextIndex
}

func (i *IndirectIndexedDrawBuffer) CmdDrawIndexed(cmd DrawIndexed) {
	if cmd.IndexCount == 0 {
		return
	}
	if cmd.InstanceCount == 0 {
		return
	}
	i.commands.Set(i.nextIndex, cmd)
	i.nextIndex++
}

func (i *IndirectIndexedDrawBuffer) EndDrawIndirect(cmd *Buffer) {
	batchCount := i.nextIndex - i.batchOffset
	if batchCount == 0 {
		return
	}

	cmd.CmdDrawIndexedIndirect(i.commands, i.batchOffset*i.commands.Stride(), batchCount, i.commands.Stride())
}

func (i *IndirectIndexedDrawBuffer) Flush() {
	i.commands.Flush()
}

func (i *IndirectIndexedDrawBuffer) Destroy() {
	i.commands.Destroy()
}
