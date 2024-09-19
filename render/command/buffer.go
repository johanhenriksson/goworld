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

type Buffer struct {
	ptr    core1_0.CommandBuffer
	pool   core1_0.CommandPool
	device *device.Device

	// cached bindings
	pipeline *pipeline.Pipeline
	vertex   bufferBinding
	index    bufferBinding
	scissor  core1_0.Rect2D
	viewport core1_0.Viewport
}

type bufferBinding struct {
	buffer    core1_0.Buffer
	offset    int
	indexType core1_0.IndexType
}

func newBuffer(device *device.Device, pool core1_0.CommandPool, ptr core1_0.CommandBuffer) *Buffer {
	return &Buffer{
		ptr:    ptr,
		pool:   pool,
		device: device,
	}
}

func (b *Buffer) Ptr() core1_0.CommandBuffer {
	return b.ptr
}

func (b *Buffer) Destroy() {
	b.ptr.Free()
	b.ptr = nil
}

func (b *Buffer) Reset() {
	b.ptr.Reset(core1_0.CommandBufferResetReleaseResources)
}

func (b *Buffer) Begin() {
	if _, err := b.ptr.Begin(core1_0.CommandBufferBeginInfo{}); err != nil {
		panic(err)
	}
}

func (b *Buffer) End() {
	b.ptr.End()
}

func (b *Buffer) CmdCopyBuffer(src, dst buffer.T, regions ...core1_0.BufferCopy) {
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

func (b *Buffer) CmdBindGraphicsPipeline(pipe *pipeline.Pipeline) {
	// if b.pipeline != nil && b.pipeline.Ptr() == pipe.Ptr() {
	// 	return
	// }
	b.ptr.CmdBindPipeline(core1_0.PipelineBindPointGraphics, pipe.Ptr())
	b.pipeline = pipe
}

func (b *Buffer) CmdBindGraphicsDescriptor(index int, set descriptor.Set) {
	if b.pipeline == nil {
		panic("bind graphics pipeline first")
	}
	b.ptr.CmdBindDescriptorSets(core1_0.PipelineBindPointGraphics, b.pipeline.Layout().Ptr(), index, []core1_0.DescriptorSet{set.Ptr()}, nil)
}

func (b *Buffer) CmdBindVertexBuffer(vtx buffer.T, offset int) {
	binding := bufferBinding{buffer: vtx.Ptr(), offset: offset}
	if b.vertex == binding {
		return
	}
	b.ptr.CmdBindVertexBuffers(0, []core1_0.Buffer{vtx.Ptr()}, []int{offset})
	b.vertex = binding
}

func (b *Buffer) CmdBindIndexBuffers(idx buffer.T, offset int, kind core1_0.IndexType) {
	binding := bufferBinding{buffer: idx.Ptr(), offset: offset, indexType: kind}
	if b.index == binding {
		return
	}
	b.ptr.CmdBindIndexBuffer(idx.Ptr(), offset, kind)
	b.index = binding
}

func (b *Buffer) CmdDraw(cmd Draw) {
	b.ptr.CmdDraw(
		int(cmd.VertexCount),
		int(cmd.InstanceCount),
		uint32(cmd.VertexOffset),
		uint32(cmd.InstanceOffset))
}

func (b *Buffer) CmdDrawIndirect(buffer buffer.T, offset, count, stride int) {
	b.ptr.CmdDrawIndirect(buffer.Ptr(), offset, count, stride)
}

func (b *Buffer) CmdDrawIndexed(cmd DrawIndexed) {
	b.ptr.CmdDrawIndexed(
		int(cmd.IndexCount),
		int(cmd.InstanceCount),
		uint32(cmd.IndexOffset),
		int(cmd.VertexOffset),
		uint32(cmd.InstanceOffset))
}

func (b *Buffer) CmdDrawIndexedIndirect(buffer buffer.T, offset, count, stride int) {
	b.ptr.CmdDrawIndexedIndirect(buffer.Ptr(), offset, count, stride)
}

func (b *Buffer) CmdBeginRenderPass(pass *renderpass.Renderpass, framebuffer *framebuffer.Framebuffer) {
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

func (b *Buffer) CmdEndRenderPass() {
	b.ptr.CmdEndRenderPass()
}

func (b *Buffer) CmdNextSubpass() {
	b.ptr.CmdNextSubpass(core1_0.SubpassContentsInline)
}

func (b *Buffer) CmdSetViewport(x, y, w, h int) core1_0.Viewport {
	prev := b.viewport
	b.viewport = core1_0.Viewport{
		X:        float32(x),
		Y:        float32(y),
		Width:    float32(w),
		Height:   float32(h),
		MinDepth: 0,
		MaxDepth: 1,
	}
	b.ptr.CmdSetViewport([]core1_0.Viewport{
		b.viewport,
	})
	return prev
}

func (b *Buffer) CmdSetScissor(x, y, w, h int) core1_0.Rect2D {
	prev := b.scissor
	b.scissor = core1_0.Rect2D{
		Offset: core1_0.Offset2D{
			X: x,
			Y: y,
		},
		Extent: core1_0.Extent2D{
			Width:  w,
			Height: h,
		},
	}
	b.ptr.CmdSetScissor([]core1_0.Rect2D{
		b.scissor,
	})
	return prev
}

func (b *Buffer) CmdPushConstant(stages core1_0.ShaderStageFlags, offset int, value any) {
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

func (b *Buffer) CmdImageBarrier(srcMask, dstMask core1_0.PipelineStageFlags, image *image.Image, oldLayout, newLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags, mipLevel, levels int) {
	b.ptr.CmdPipelineBarrier(core1_0.PipelineStageFlags(srcMask), core1_0.PipelineStageFlags(dstMask), core1_0.DependencyFlags(0), nil, nil, []core1_0.ImageMemoryBarrier{
		{
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Image:     image.Ptr(),
			SubresourceRange: core1_0.ImageSubresourceRange{
				AspectMask:   core1_0.ImageAspectFlags(aspects),
				BaseMipLevel: mipLevel,
				LevelCount:   levels,
				LayerCount:   1,
			},
			SrcAccessMask: core1_0.AccessMemoryRead | core1_0.AccessMemoryWrite,
			DstAccessMask: core1_0.AccessMemoryRead | core1_0.AccessMemoryWrite,
		},
	})
}

func (b *Buffer) CmdCopyBufferToImage(source buffer.T, dst *image.Image, layout core1_0.ImageLayout) {
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

func (b *Buffer) CmdCopyImageToBuffer(src *image.Image, srcLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags, dst buffer.T) {
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

func (b *Buffer) CmdConvertImage(src *image.Image, srcLayout core1_0.ImageLayout, dst *image.Image, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags) {
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

func (b *Buffer) CmdCopyImage(src *image.Image, srcLayout core1_0.ImageLayout, dst *image.Image, dstLayout core1_0.ImageLayout, aspects core1_0.ImageAspectFlags) {
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
