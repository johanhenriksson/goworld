package command

import (
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/vkngwrapper/core/v2/core1_0"

	vk "github.com/vulkan-go/vulkan"
)

type Buffer interface {
	device.Resource[core1_0.CommandBuffer]

	Reset()
	Begin()
	End()

	CmdCopyBuffer(src, dst buffer.T, regions ...vk.BufferCopy)
	CmdBindGraphicsPipeline(pipe pipeline.T)
	CmdBindGraphicsDescriptor(sets descriptor.Set)
	CmdBindVertexBuffer(vtx buffer.T, offset int)
	CmdBindIndexBuffers(idx buffer.T, offset int, kind vk.IndexType)
	CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int)
	CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int)
	CmdBeginRenderPass(pass renderpass.T, framebuffer framebuffer.T)
	CmdNextSubpass()
	CmdEndRenderPass()
	CmdSetViewport(x, y, w, h int)
	CmdSetScissor(x, y, w, h int)
	CmdPushConstant(stages vk.ShaderStageFlagBits, offset int, value any)
	CmdImageBarrier(srcMask, dstMask vk.PipelineStageFlagBits, image image.T, oldLayout, newLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits)
	CmdCopyBufferToImage(source buffer.T, dst image.T, layout vk.ImageLayout)
	CmdCopyImageToBuffer(src image.T, srcLayout vk.ImageLayout, aspect vk.ImageAspectFlagBits, dst buffer.T)
	CmdConvertImage(src image.T, srcLayout vk.ImageLayout, dst image.T, dstLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits)
	CmdCopyImage(src image.T, srcLayout vk.ImageLayout, dst image.T, dstLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits)
}

type buf struct {
	ptr    core1_0.CommandBuffer
	pool   core1_0.CommandPool
	device device.T

	// cached bindings
	pipeline pipeline.T
	vertex   bufferBinding
	index    bufferBinding
}

type bufferBinding struct {
	buffer    core1_0.Buffer
	offset    int
	indexType core1_0.IndexType
}

func newBuffer(device device.T, pool core1_0.CommandPool, ptr core1_0.CommandBuffer) Buffer {
	return &buf{
		ptr:    ptr,
		pool:   pool,
		device: device,
	}
}

func (b *buf) Ptr() core1_0.CommandBuffer {
	return b.ptr
}

func (b *buf) Destroy() {
	b.ptr.Free()
	b.ptr = nil
}

func (b *buf) Reset() {
	b.ptr.Reset(core1_0.CommandBufferResetReleaseResources)
}

func (b *buf) Begin() {
	b.ptr.Begin(core1_0.CommandBufferBeginInfo{})
}

func (b *buf) End() {
	b.ptr.End()
}

func (b *buf) CmdCopyBuffer(src, dst buffer.T, regions ...core1_0.BufferCopy) {
	if len(regions) == 0 {
		regions = []core1_0.BufferCopy{
			{
				SrcOffset: 0,
				DstOffset: 0,
				Size:      src.Size(),
			},
		}
	}
	if src.Ptr() == nil || dst.Ptr() == nil {
		panic("copy to/from null buffer")
	}
	b.ptr.CmdCopyBuffer(src.Ptr(), dst.Ptr(), regions)
}

func (b *buf) CmdBindGraphicsPipeline(pipe pipeline.T) {
	if b.pipeline != nil && b.pipeline.Ptr() == pipe.Ptr() {
		return
	}
	vk.CmdBindPipeline(b.Ptr(), vk.PipelineBindPointGraphics, pipe.Ptr())
	b.pipeline = pipe
}

func (b *buf) CmdBindGraphicsDescriptor(set descriptor.Set) {
	if b.pipeline == nil {
		panic("bind graphics pipeline first")
	}
	vk.CmdBindDescriptorSets(
		b.Ptr(),
		vk.PipelineBindPointGraphics,
		b.pipeline.Layout().Ptr(), 0, 1,
		[]vk.DescriptorSet{set.Ptr()},
		0, nil)
}

func (b *buf) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	binding := bufferBinding{buffer: vtx.Ptr(), offset: offset}
	if b.vertex == binding {
		return
	}
	vk.CmdBindVertexBuffers(b.Ptr(), 0, 1, []vk.Buffer{vtx.Ptr()}, []vk.DeviceSize{vk.DeviceSize(offset)})
	b.vertex = binding
}

func (b *buf) CmdBindIndexBuffers(idx buffer.T, offset int, kind vk.IndexType) {
	binding := bufferBinding{buffer: idx.Ptr(), offset: offset, indexType: kind}
	if b.index == binding {
		return
	}
	vk.CmdBindIndexBuffer(b.Ptr(), idx.Ptr(), vk.DeviceSize(offset), kind)
	b.index = binding
}

func (b *buf) CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int) {
	vk.CmdDraw(b.Ptr(), uint32(vertexCount), uint32(instanceCount), uint32(firstVertex), uint32(firstInstance))
}

func (b *buf) CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int) {
	vk.CmdDrawIndexed(b.Ptr(), uint32(indexCount), uint32(instanceCount), uint32(firstIndex), int32(vertexOffset), uint32(firstInstance))
}

func (b *buf) CmdBeginRenderPass(pass renderpass.T, framebuffer framebuffer.T) {
	clear := pass.Clear()
	w, h := framebuffer.Size()

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
		ClearValueCount: uint32(len(clear)),
		PClearValues:    clear,
	}, vk.SubpassContentsInline)

	b.CmdSetViewport(0, 0, w, h)
	b.CmdSetScissor(0, 0, w, h)
}

func (b *buf) CmdEndRenderPass() {
	vk.CmdEndRenderPass(b.ptr)
}

func (b *buf) CmdNextSubpass() {
	vk.CmdNextSubpass(b.ptr, vk.SubpassContentsInline)
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

func (b *buf) CmdPushConstant(stages vk.ShaderStageFlagBits, offset int, value any) {
	if b.pipeline == nil {
		panic("bind graphics pipeline first")
	}
	ptr := reflect.ValueOf(value).UnsafePointer()
	size := reflect.ValueOf(value).Elem().Type().Size()
	vk.CmdPushConstants(
		b.ptr,
		b.pipeline.Layout().Ptr(),
		vk.ShaderStageFlags(stages),
		uint32(offset), uint32(size),
		ptr)
}

func (b *buf) CmdImageBarrier(srcMask, dstMask vk.PipelineStageFlagBits, image image.T, oldLayout, newLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits) {
	vk.CmdPipelineBarrier(b.ptr, vk.PipelineStageFlags(srcMask), vk.PipelineStageFlags(dstMask), vk.DependencyFlags(0), 0, nil, 0, nil, 1, []vk.ImageMemoryBarrier{
		{
			SType:     vk.StructureTypeImageMemoryBarrier,
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Image:     image.Ptr(),
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask: vk.ImageAspectFlags(aspects),
				LayerCount: 1,
				LevelCount: 1,
			},
			SrcAccessMask: vk.AccessFlags(vk.AccessMemoryReadBit | vk.AccessMemoryWriteBit),
			DstAccessMask: vk.AccessFlags(vk.AccessMemoryReadBit | vk.AccessMemoryWriteBit),
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

func (b *buf) CmdCopyImageToBuffer(src image.T, srcLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits, dst buffer.T) {
	vk.CmdCopyImageToBuffer(b.ptr, src.Ptr(), srcLayout, dst.Ptr(), 1, []vk.BufferImageCopy{
		{
			ImageSubresource: vk.ImageSubresourceLayers{
				AspectMask: vk.ImageAspectFlags(aspects),
				LayerCount: 1,
			},
			ImageExtent: vk.Extent3D{
				Width:  uint32(src.Width()),
				Height: uint32(src.Height()),
				Depth:  1,
			},
		},
	})
}

func (b *buf) CmdConvertImage(src image.T, srcLayout vk.ImageLayout, dst image.T, dstLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits) {
	vk.CmdBlitImage(b.ptr, src.Ptr(), srcLayout, dst.Ptr(), dstLayout, 1, []vk.ImageBlit{
		{
			SrcOffsets: [2]vk.Offset3D{
				{X: 0, Y: 0, Z: 0},
				{X: int32(src.Width()), Y: int32(src.Height()), Z: 1},
			},
			SrcSubresource: vk.ImageSubresourceLayers{
				AspectMask:     vk.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			DstOffsets: [2]vk.Offset3D{
				{X: 0, Y: 0, Z: 0},
				{X: int32(dst.Width()), Y: int32(dst.Height()), Z: 1},
			},
			DstSubresource: vk.ImageSubresourceLayers{
				AspectMask:     vk.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		},
	}, vk.FilterNearest)
}

func (b *buf) CmdCopyImage(src image.T, srcLayout vk.ImageLayout, dst image.T, dstLayout vk.ImageLayout, aspects vk.ImageAspectFlagBits) {
	vk.CmdCopyImage(b.ptr, src.Ptr(), srcLayout, dst.Ptr(), dstLayout, 1, []vk.ImageCopy{
		{
			SrcOffset: vk.Offset3D{X: 0, Y: 0, Z: 0},
			SrcSubresource: vk.ImageSubresourceLayers{
				AspectMask:     vk.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			DstOffset: vk.Offset3D{X: 0, Y: 0, Z: 0},
			DstSubresource: vk.ImageSubresourceLayers{
				AspectMask:     vk.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			Extent: vk.Extent3D{
				Width:  uint32(dst.Width()),
				Height: uint32(dst.Height()),
				Depth:  1,
			},
		},
	})
}
