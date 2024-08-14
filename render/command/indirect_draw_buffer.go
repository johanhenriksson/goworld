package command

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type IndirectDrawBuffer struct {
	commands    *buffer.Array[DrawIndexed]
	nextIndex   int
	batchOffset int
}

var _ DrawIndexedBuffer = (*IndirectDrawBuffer)(nil)

func NewIndirectDrawBuffer(dev *device.Device, key string, size int) *IndirectDrawBuffer {
	cmds := buffer.NewArray[DrawIndexed](dev, buffer.Args{
		Key:  key,
		Size: size,
		Usage: core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer |
			core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: device.MemoryTypeShared,
	})
	return &IndirectDrawBuffer{
		commands: cmds,
	}
}

func (i *IndirectDrawBuffer) Reset() {
	i.nextIndex = 0
	i.batchOffset = 0
}

func (i *IndirectDrawBuffer) BeginDrawIndirect() {
	i.batchOffset = i.nextIndex
}

func (i *IndirectDrawBuffer) CmdDrawIndexed(cmd DrawIndexed) {
	if cmd.IndexCount == 0 {
		return
	}
	if cmd.InstanceCount == 0 {
		return
	}
	i.commands.Set(i.nextIndex, cmd)
	i.nextIndex++
}

func (i *IndirectDrawBuffer) EndDrawIndirect(cmd *Buffer) {
	batchCount := i.nextIndex - i.batchOffset
	if batchCount == 0 {
		return
	}

	cmd.CmdDrawIndexedIndirect(i.commands, i.batchOffset*i.commands.Stride(), batchCount, i.commands.Stride())
}

func (i *IndirectDrawBuffer) Flush() {
	i.commands.Flush()
}

func (i *IndirectDrawBuffer) Destroy() {
	i.commands.Destroy()
}
