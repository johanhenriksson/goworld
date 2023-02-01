package command

import (
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Buffer interface {
	device.Resource[core1_0.CommandBuffer]

	Reset()
	Begin()
	End()

	CmdCopyBuffer(src, dst buffer.T, regions ...core1_0.BufferCopy)
	CmdBindGraphicsPipeline(pipe pipeline.T)
	CmdBindGraphicsDescriptor(sets descriptor.Set)
	CmdBindVertexBuffer(vtx buffer.T, offset int)
	CmdBindIndexBuffers(idx buffer.T, offset int, kind core1_0.IndexType)
	CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int)
	CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int)
	CmdBeginRenderPass(pass renderpass.T, framebuffer framebuffer.T)
	CmdNextSubpass()
	CmdEndRenderPass()
	CmdSetViewport(x, y, w, h int)
	CmdSetScissor(x, y, w, h int)
	CmdPushConstant(stages core1_0.ShaderStageFlags, offset int, value any)
	CmdImageBarrier(srcMask, dstMask core1_0.PipelineStageFlags, image image.T, oldLayout, newLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags)
	CmdCopyBufferToImage(source buffer.T, dst image.T, layout core1_0.ImageLayout)
	CmdCopyImageToBuffer(src image.T, srcLayout core1_0.ImageLayout, aspect core1_0.ImageAspectFlags, dst buffer.T)
	CmdConvertImage(src image.T, srcLayout core1_0.ImageLayout, dst image.T, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags)
	CmdCopyImage(src image.T, srcLayout core1_0.ImageLayout, dst image.T, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags)
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
	if _, err := b.ptr.Begin(core1_0.CommandBufferBeginInfo{}); err != nil {
		panic(err)
	}
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
	b.ptr.CmdBindPipeline(core1_0.PipelineBindPointGraphics, pipe.Ptr())
	b.pipeline = pipe
}

func (b *buf) CmdBindGraphicsDescriptor(set descriptor.Set) {
	if b.pipeline == nil {
		panic("bind graphics pipeline first")
	}
	b.ptr.CmdBindDescriptorSets(core1_0.PipelineBindPointGraphics, b.pipeline.Layout().Ptr(), 0, []core1_0.DescriptorSet{set.Ptr()}, nil)
}

func (b *buf) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	binding := bufferBinding{buffer: vtx.Ptr(), offset: offset}
	if b.vertex == binding {
		return
	}
	b.ptr.CmdBindVertexBuffers(0, []core1_0.Buffer{vtx.Ptr()}, []int{offset})
	b.vertex = binding
}

func (b *buf) CmdBindIndexBuffers(idx buffer.T, offset int, kind core1_0.IndexType) {
	binding := bufferBinding{buffer: idx.Ptr(), offset: offset, indexType: kind}
	if b.index == binding {
		return
	}
	b.ptr.CmdBindIndexBuffer(idx.Ptr(), offset, kind)
	b.index = binding
}

func (b *buf) CmdDraw(vertexCount, instanceCount, firstVertex, firstInstance int) {
	b.ptr.CmdDraw(vertexCount, instanceCount, uint32(firstVertex), uint32(firstInstance))
}

func (b *buf) CmdDrawIndexed(indexCount, instanceCount, firstIndex, vertexOffset, firstInstance int) {
	b.ptr.CmdDrawIndexed(indexCount, instanceCount, uint32(firstIndex), vertexOffset, uint32(firstInstance))
}

func (b *buf) CmdBeginRenderPass(pass renderpass.T, framebuffer framebuffer.T) {
	clear := pass.Clear()
	w, h := framebuffer.Size()

	b.ptr.CmdBeginRenderPass(core1_0.SubpassContentsInline, core1_0.RenderPassBeginInfo{
		RenderPass:  pass.Ptr(),
		Framebuffer: framebuffer.Ptr(),
		RenderArea: core1_0.Rect2D{
			Offset: core1_0.Offset2D{},
			Extent: core1_0.Extent2D{
				Width:  w,
				Height: h,
			},
		},
		ClearValues: clear,
	})

	b.CmdSetViewport(0, 0, w, h)
	b.CmdSetScissor(0, 0, w, h)
}

func (b *buf) CmdEndRenderPass() {
	b.ptr.CmdEndRenderPass()
}

func (b *buf) CmdNextSubpass() {
	b.ptr.CmdNextSubpass(core1_0.SubpassContentsInline)
}

func (b *buf) CmdSetViewport(x, y, w, h int) {
	b.ptr.CmdSetViewport([]core1_0.Viewport{
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
	b.ptr.CmdSetScissor([]core1_0.Rect2D{
		{
			Offset: core1_0.Offset2D{
				X: x,
				Y: y,
			},
			Extent: core1_0.Extent2D{
				Width:  w,
				Height: h,
			},
		},
	})
}

func (b *buf) CmdPushConstant(stages core1_0.ShaderStageFlags, offset int, value any) {
	if b.pipeline == nil {
		panic("bind graphics pipeline first")
	}
	// this is awkward
	// ptr := reflect.ValueOf(value).UnsafePointer()
	size := reflect.ValueOf(value).Elem().Type().Size()
	ptr := reflect.ValueOf(value).UnsafePointer()
	valueBytes := make([]byte, size)

	device.Memcpy(unsafe.Pointer(&valueBytes[0]), ptr, int(size))
	b.ptr.CmdPushConstants(b.pipeline.Layout().Ptr(), stages, offset, valueBytes)
}

func (b *buf) CmdImageBarrier(srcMask, dstMask core1_0.PipelineStageFlags, image image.T, oldLayout, newLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags) {
	b.ptr.CmdPipelineBarrier(core1_0.PipelineStageFlags(srcMask), core1_0.PipelineStageFlags(dstMask), core1_0.DependencyFlags(0), nil, nil, []core1_0.ImageMemoryBarrier{
		{
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Image:     image.Ptr(),
			SubresourceRange: core1_0.ImageSubresourceRange{
				AspectMask: core1_0.ImageAspectFlags(aspects),
				LayerCount: 1,
				LevelCount: 1,
			},
			SrcAccessMask: core1_0.AccessMemoryRead | core1_0.AccessMemoryWrite,
			DstAccessMask: core1_0.AccessMemoryRead | core1_0.AccessMemoryWrite,
		},
	})
}

func (b *buf) CmdCopyBufferToImage(source buffer.T, dst image.T, layout core1_0.ImageLayout) {
	b.ptr.CmdCopyBufferToImage(source.Ptr(), dst.Ptr(), layout, []core1_0.BufferImageCopy{
		{
			ImageSubresource: core1_0.ImageSubresourceLayers{
				AspectMask: core1_0.ImageAspectColor,
				LayerCount: 1,
			},
			ImageExtent: core1_0.Extent3D{
				Width:  dst.Width(),
				Height: dst.Height(),
				Depth:  1,
			},
		},
	})
}

func (b *buf) CmdCopyImageToBuffer(src image.T, srcLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags, dst buffer.T) {
	b.ptr.CmdCopyImageToBuffer(src.Ptr(), srcLayout, dst.Ptr(), []core1_0.BufferImageCopy{
		{
			ImageSubresource: core1_0.ImageSubresourceLayers{
				AspectMask: core1_0.ImageAspectFlags(aspects),
				LayerCount: 1,
			},
			ImageExtent: core1_0.Extent3D{
				Width:  src.Width(),
				Height: src.Height(),
				Depth:  1,
			},
		},
	})
}

func (b *buf) CmdConvertImage(src image.T, srcLayout core1_0.ImageLayout, dst image.T, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags) {
	b.ptr.CmdBlitImage(src.Ptr(), srcLayout, dst.Ptr(), dstLayout, []core1_0.ImageBlit{
		{
			SrcOffsets: [2]core1_0.Offset3D{
				{X: 0, Y: 0, Z: 0},
				{X: src.Width(), Y: src.Height(), Z: 1},
			},
			SrcSubresource: core1_0.ImageSubresourceLayers{
				AspectMask:     core1_0.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			DstOffsets: [2]core1_0.Offset3D{
				{X: 0, Y: 0, Z: 0},
				{X: dst.Width(), Y: dst.Height(), Z: 1},
			},
			DstSubresource: core1_0.ImageSubresourceLayers{
				AspectMask:     core1_0.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		},
	}, core1_0.FilterNearest)
}

func (b *buf) CmdCopyImage(src image.T, srcLayout core1_0.ImageLayout, dst image.T, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags) {
	b.ptr.CmdCopyImage(src.Ptr(), srcLayout, dst.Ptr(), dstLayout, []core1_0.ImageCopy{
		{
			SrcOffset: core1_0.Offset3D{X: 0, Y: 0, Z: 0},
			SrcSubresource: core1_0.ImageSubresourceLayers{
				AspectMask:     core1_0.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			DstOffset: core1_0.Offset3D{X: 0, Y: 0, Z: 0},
			DstSubresource: core1_0.ImageSubresourceLayers{
				AspectMask:     core1_0.ImageAspectFlags(aspects),
				MipLevel:       0,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
			Extent: core1_0.Extent3D{
				Width:  dst.Width(),
				Height: dst.Height(),
				Depth:  1,
			},
		},
	})
}
