package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type IndirectDrawBuffer struct {
	commands buffer.Array[command.DrawIndirectIndexed]
	count    int
}

func NewIndirectDrawBuffer(device device.T, key string, size int) *IndirectDrawBuffer {
	cmds := buffer.NewArray[command.DrawIndirectIndexed](device, buffer.Args{
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

func (i *IndirectDrawBuffer) BeginDrawIndirect() {
	i.count = 0
}

func (i *IndirectDrawBuffer) DrawIndexed(indexCount, firstIndex, vertexOffset, firstInstance, instanceCount int) {
	i.commands.Set(i.count, command.DrawIndirectIndexed{
		IndexCount:    uint32(indexCount),
		InstanceCount: uint32(instanceCount),
		FirstIndex:    uint32(firstIndex),
		VertexOffset:  int32(vertexOffset),
		FirstInstance: uint32(firstInstance),
	})
	i.count++
}

func (i *IndirectDrawBuffer) EndDrawIndirect(cmd command.Buffer) {
	// flush?
	cmd.CmdDrawIndexedIndirect(i.commands, 0, i.count, i.commands.Stride())
}

func (i *IndirectDrawBuffer) Flush() {
	i.commands.Flush()
}

func (i *IndirectDrawBuffer) Destroy() {
	i.commands.Destroy()
}
