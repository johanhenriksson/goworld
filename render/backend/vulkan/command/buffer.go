package command

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"

	vk "github.com/vulkan-go/vulkan"
)

type Buffer interface {
	device.Resource[vk.CommandBuffer]

	SubmitSync(vk.Queue)
	Begin()
	End()

	CopyBuffer(src, dst buffer.T, regions ...vk.BufferCopy)
}

type buf struct {
	ptr    vk.CommandBuffer
	pool   vk.CommandPool
	device device.T
}

func newBuffer(device device.T, pool vk.CommandPool, ptr vk.CommandBuffer) Buffer {
	return &buf{
		ptr:    ptr,
		pool:   pool,
		device: device,
	}
}

func (b *buf) Ptr() vk.CommandBuffer {
	return b.ptr
}

func (b *buf) Destroy() {
	vk.FreeCommandBuffers(b.device.Ptr(), b.pool, 1, []vk.CommandBuffer{b.ptr})
	b.ptr = nil
}

func (b *buf) SubmitSync(queue vk.Queue) {
	fence := sync.NewFence(b.device, false)
	defer fence.Destroy()

	info := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   0,
		SignalSemaphoreCount: 0,
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{b.ptr},
		PWaitDstStageMask:    []vk.PipelineStageFlags{},
	}
	vk.QueueSubmit(queue, 1, []vk.SubmitInfo{info}, fence.Ptr())

	fence.Wait()
}

func (b *buf) Begin() {
	info := vk.CommandBufferBeginInfo{
		SType: vk.StructureTypeCommandBufferBeginInfo,
	}
	vk.BeginCommandBuffer(b.ptr, &info)
}

func (b *buf) End() {
	vk.EndCommandBuffer(b.ptr)
}

func (b *buf) CopyBuffer(src, dst buffer.T, regions ...vk.BufferCopy) {
	vk.CmdCopyBuffer(b.ptr, src.Ptr(), dst.Ptr(), uint32(len(regions)), regions)
}
