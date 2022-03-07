package command

import (
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/color"
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
	CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int)
	CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int)
	CmdBeginRenderPass(pass pipeline.Pass, framebuffer framebuffer.T, clear color.T)
	CmdEndRenderPass()
	CmdSetViewport(x, y, w, h int)
	CmdSetScissor(x, y, w, h int)
	CmdPushConstant(layout pipeline.Layout, stages vk.ShaderStageFlags, offset int, value any)
	CmdImageBarrier(srcMask, dstMask vk.PipelineStageFlags, image image.T, oldLayout, newLayout vk.ImageLayout)
	CmdCopyBufferToImage(source buffer.T, dst image.T, layout vk.ImageLayout)
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
		layout.Ptr(), 0, uint32(len(sets)),
		util.Map(sets, func(i int, s descriptor.Set) vk.DescriptorSet { return s.Ptr() }),
		0, nil)
}

func (b *buf) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	vk.CmdBindVertexBuffers(b.Ptr(), 0, 1, []vk.Buffer{vtx.Ptr()}, []vk.DeviceSize{vk.DeviceSize(offset)})
}

func (b *buf) CmdBindIndexBuffers(idx buffer.T, offset int, kind vk.IndexType) {
	vk.CmdBindIndexBuffer(b.Ptr(), idx.Ptr(), vk.DeviceSize(offset), kind)
}

func (b *buf) CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int) {
	vk.CmdDraw(b.Ptr(), uint32(vertexCount), uint32(instanceCount), uint32(firstVertex), uint32(firstInstance))
}

func (b *buf) CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int) {
	vk.CmdDrawIndexed(b.Ptr(), uint32(indexCount), uint32(instanceCount), uint32(firstIndex), int32(vertexOffset), uint32(firstInstance))
}

func (b *buf) CmdBeginRenderPass(pass pipeline.Pass, framebuffer framebuffer.T, clear color.T) {
	w, h := framebuffer.Size()

	clearValues := make([]vk.ClearValue, 2)
	clearValues[1].SetDepthStencil(1, 0)
	clearValues[0].SetColor([]float32{clear.R, clear.G, clear.B, clear.A})

	vk.CmdBeginRenderPass(b.Ptr(), &vk.RenderPassBeginInfo{
		SType:       vk.StructureTypeRenderPassBeginInfo,
		RenderPass:  pass.Ptr(),
		Framebuffer: framebuffer.Ptr(),
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{},
			Extent: vk.Extent2D{
				Width:  uint32(w),
				Height: uint32(h),
			},
		},
		ClearValueCount: 2,
		PClearValues:    clearValues,
	}, vk.SubpassContentsInline)
}

func (b *buf) CmdEndRenderPass() {
	vk.CmdEndRenderPass(b.ptr)
}

func (b *buf) CmdSetViewport(x, y, w, h int) {
	vk.CmdSetViewport(b.Ptr(), 0, 1, []vk.Viewport{
		{
			X:        float32(x),
			Y:        float32(y),
			Width:    float32(w),
			Height:   float32(h),
			MinDepth: 0,
			MaxDepth: 1,
		},
	})
}

func (b *buf) CmdSetScissor(x, y, w, h int) {
	vk.CmdSetScissor(b.Ptr(), 0, 1, []vk.Rect2D{
		{
			Offset: vk.Offset2D{
				X: int32(x),
				Y: int32(y),
			},
			Extent: vk.Extent2D{
				Width:  uint32(w),
				Height: uint32(h),
			},
		},
	})
}

func (b *buf) CmdPushConstant(layout pipeline.Layout, stages vk.ShaderStageFlags, offset int, value any) {
	ptr := reflect.ValueOf(value).UnsafePointer()
	size := reflect.ValueOf(value).Elem().Type().Size()
	vk.CmdPushConstants(b.ptr, layout.Ptr(), stages, uint32(offset), uint32(size), unsafe.Pointer(ptr))
}

func (b *buf) CmdImageBarrier(srcMask, dstMask vk.PipelineStageFlags, image image.T, oldLayout, newLayout vk.ImageLayout) {
	vk.CmdPipelineBarrier(b.ptr, srcMask, dstMask, vk.DependencyFlags(0), 0, nil, 0, nil, 1, []vk.ImageMemoryBarrier{
		{
			SType:     vk.StructureTypeImageMemoryBarrier,
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Image:     image.Ptr(),
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask: vk.ImageAspectFlags(vk.ImageAspectColorBit),
				LayerCount: 1,
				LevelCount: 1,
			},
		},
	})
}

func (b *buf) CmdCopyBufferToImage(source buffer.T, dst image.T, layout vk.ImageLayout) {
	vk.CmdCopyBufferToImage(b.ptr, source.Ptr(), dst.Ptr(), layout, 1, []vk.BufferImageCopy{
		{
			ImageSubresource: vk.ImageSubresourceLayers{
				AspectMask: vk.ImageAspectFlags(vk.ImageAspectColorBit),
				LayerCount: 1,
			},
			ImageExtent: vk.Extent3D{
				Width:  uint32(dst.Width()),
				Height: uint32(dst.Height()),
				Depth:  1,
			},
		},
	})
}
