package command

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Draw struct {
	VertexCount    uint32
	InstanceCount  uint32
	VertexOffset   int32
	InstanceOffset uint32
}

type DrawBuffer interface {
	CmdDraw(Draw)
}

type IndirectDrawBuffer struct {
	Commands    *buffer.Array[Draw]
	nextIndex   int
	batchOffset int
}

var _ DrawBuffer = (*IndirectDrawBuffer)(nil)

func NewIndirectDrawBuffer(dev *device.Device, key string, size int) *IndirectDrawBuffer {
	cmds := buffer.NewArray[Draw](dev, buffer.Args{
		Key:  key,
		Size: size,
		Usage: core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer |
			core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: device.MemoryTypeShared,
	})
	return &IndirectDrawBuffer{
		Commands: cmds,
	}
}

func (i *IndirectDrawBuffer) Count() int {
	return i.nextIndex
}

func (i *IndirectDrawBuffer) Reset() {
	i.nextIndex = 0
	i.batchOffset = 0
}

func (i *IndirectDrawBuffer) BeginDrawIndirect() {
	i.batchOffset = i.nextIndex
}

func (i *IndirectDrawBuffer) CmdDraw(cmd Draw) {
	if cmd.InstanceCount == 0 {
		return
	}
	i.Commands.Set(i.nextIndex, cmd)
	i.nextIndex++
}

func (i *IndirectDrawBuffer) EndDrawIndirect(cmd *Buffer) {
	batchCount := i.nextIndex - i.batchOffset
	if batchCount == 0 {
		return
	}

	cmd.CmdDrawIndirect(i.Commands, i.batchOffset*i.Commands.Stride(), batchCount, i.Commands.Stride())
}

func (i *IndirectDrawBuffer) Flush() {
	i.Commands.Flush()
}

func (i *IndirectDrawBuffer) Destroy() {
	i.Commands.Destroy()
}
