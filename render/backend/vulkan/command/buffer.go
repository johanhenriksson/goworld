package command

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Buffer interface {
	device.Resource[vk.CommandBuffer]

	SubmitSync(vk.Queue)
	Reset()
	Begin()
	End()

	CmdCopyBuffer(src, dst buffer.T, regions ...vk.BufferCopy)
	CmdBindGraphicsPipeline(pipe pipeline.T)
	CmdBindGraphicsDescriptors(layout pipeline.Layout, sets []descriptor.Set)
	CmdBindVertexBuffer(vtx buffer.T, offset int)
	CmdBindIndexBuffers(idx buffer.T, offset int, kind vk.IndexType)
	CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int)
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

func (b *buf) Reset() {
	vk.ResetCommandBuffer(b.ptr, vk.CommandBufferResetFlags(vk.CommandBufferResetReleaseResourcesBit))
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

func (b *buf) CmdCopyBuffer(src, dst buffer.T, regions ...vk.BufferCopy) {
	if len(regions) == 0 {
		regions = []vk.BufferCopy{
			{
				SrcOffset: 0,
				DstOffset: 0,
				Size:      vk.DeviceSize(src.Size()),
			},
		}
	}
	if src.Ptr() == nil || dst.Ptr() == nil {
		panic("copy to/from null buffer")
	}
	vk.CmdCopyBuffer(b.ptr, src.Ptr(), dst.Ptr(), uint32(len(regions)), regions)
}

func (b *buf) CmdBindGraphicsPipeline(pipe pipeline.T) {
	vk.CmdBindPipeline(b.Ptr(), vk.PipelineBindPointGraphics, pipe.Ptr())
}

func (b *buf) CmdBindGraphicsDescriptors(layout pipeline.Layout, sets []descriptor.Set) {
	vk.CmdBindDescriptorSets(
		b.Ptr(),
		vk.PipelineBindPointGraphics,
		layout.Ptr(), 0, 1,
		util.Map(sets, func(i int, s descriptor.Set) vk.DescriptorSet { return s.Ptr() }),
		0, nil)
}

func (b *buf) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	vk.CmdBindVertexBuffers(b.Ptr(), 0, 1, []vk.Buffer{vtx.Ptr()}, []vk.DeviceSize{vk.DeviceSize(offset)})
}

func (b *buf) CmdBindIndexBuffers(idx buffer.T, offset int, kind vk.IndexType) {
	vk.CmdBindIndexBuffer(b.Ptr(), idx.Ptr(), vk.DeviceSize(offset), kind)
}

func (b *buf) CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int) {
	vk.CmdDrawIndexed(b.Ptr(), uint32(indexCount), uint32(instanceCount), uint32(firstIndex), int32(vertexOffset), uint32(firstInstance))
}
