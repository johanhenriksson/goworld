package command

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type IndirectDrawBuffer struct {
	commands    buffer.Array[DrawIndirectIndexed]
	nextIndex   int
	batchOffset int
}

func NewIndirectDrawBuffer(device device.T, key string, size int) *IndirectDrawBuffer {
	cmds := buffer.NewArray[DrawIndirectIndexed](device, buffer.Args{
		Key:  key,
		Size: size,
		Usage: core1_0.BufferUsageStorageBuffer | core1_0.BufferUsageIndirectBuffer |
			core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: core1_0.MemoryPropertyDeviceLocal | core1_0.MemoryPropertyHostVisible,
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

func (i *IndirectDrawBuffer) DrawIndexed(indexCount, firstIndex, vertexOffset, firstInstance, instanceCount int) {
	if indexCount == 0 {
		return
	}
	if instanceCount == 0 {
		return
	}
	i.commands.Set(i.nextIndex, DrawIndirectIndexed{
		IndexCount:    uint32(indexCount),
		InstanceCount: uint32(instanceCount),
		FirstIndex:    uint32(firstIndex),
		VertexOffset:  int32(vertexOffset),
		FirstInstance: uint32(firstInstance),
	})
	i.nextIndex++
}

func (i *IndirectDrawBuffer) EndDrawIndirect(cmd Buffer) {
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
