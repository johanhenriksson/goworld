package pass

import (
	"image/color"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
	"github.com/x448/float16"
)

type BufferOutput interface {
	SamplePosition(cursor vec2.T) (vec3.T, bool)
	SampleNormal(cursor vec2.T) (vec3.T, bool)
}

type GeometryBuffer interface {
	BufferOutput

	Diffuse() image.View
	Normal() image.View
	Position() image.View
	Depth() image.View
	Output() image.View
	Destroy()

	CopyBuffers() command.CommandFn
}

type gbuffer struct {
	width  int
	height int

	diffuse  image.View
	normal   image.View
	position image.View
	output   image.View
	depth    image.View

	normalBuf   image.T
	positionBuf image.T
}

func NewGbuffer(
	target vulkan.Target,
	diffuse, normal, position, output, depth image.View,
) GeometryBuffer {
	positionBuf, err := image.New(target.Device(), image.Args{
		Type:   vk.ImageType2d,
		Width:  position.Image().Width(),
		Height: position.Image().Height(),
		Format: position.Format(),
		Tiling: vk.ImageTilingLinear,
		Usage:  vk.ImageUsageTransferDstBit,
		Memory: vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
	})
	if err != nil {
		panic(err)
	}

	normalBuf, err := image.New(target.Device(), image.Args{
		Type:   vk.ImageType2d,
		Width:  normal.Image().Width(),
		Height: normal.Image().Height(),
		Format: normal.Format(),
		Tiling: vk.ImageTilingLinear,
		Usage:  vk.ImageUsageTransferDstBit,
		Memory: vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
	})
	if err != nil {
		panic(err)
	}

	// move images to ImageLayoutGeneral to avoid errors on first copy
	worker := target.Transferer()
	worker.Queue(func(b command.Buffer) {
		b.CmdImageBarrier(vk.PipelineStageTopOfPipeBit, vk.PipelineStageTransferBit, positionBuf, vk.ImageLayoutUndefined, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
		b.CmdImageBarrier(vk.PipelineStageTopOfPipeBit, vk.PipelineStageTransferBit, normalBuf, vk.ImageLayoutUndefined, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
	})
	worker.Submit(command.SubmitInfo{
		Marker: "GBufferInit",
	})

	return &gbuffer{
		width:  target.Width(),
		height: target.Height(),

		diffuse:  diffuse,
		normal:   normal,
		position: position,
		depth:    depth,
		output:   output,

		positionBuf: positionBuf,
		normalBuf:   normalBuf,
	}
}

func (b *gbuffer) Diffuse() image.View  { return b.diffuse }
func (b *gbuffer) Normal() image.View   { return b.normal }
func (b *gbuffer) Position() image.View { return b.position }
func (b *gbuffer) Depth() image.View    { return b.depth }
func (b *gbuffer) Output() image.View   { return b.output }

func (b *gbuffer) pixelOffset(pos vec2.T, img image.T, size int) int {
	x := int(pos.X * float32(img.Width()) / float32(b.width))
	y := int(pos.Y * float32(img.Height()) / float32(b.height))

	return size * (y*img.Width() + x)
}

func (b *gbuffer) SamplePosition(cursor vec2.T) (vec3.T, bool) {
	if cursor.X < 0 || cursor.Y < 0 || cursor.X > float32(b.width) || cursor.Y > float32(b.height) {
		return vec3.Zero, false
	}

	offset := b.pixelOffset(cursor, b.normalBuf, 8)
	output := make([]uint16, 4)
	b.positionBuf.Memory().Read(offset, output)

	if output[0] == 0 && output[1] == 0 && output[2] == 0 {
		return vec3.Zero, false
	}

	return vec3.New(
		float16.Frombits(output[0]).Float32(),
		float16.Frombits(output[1]).Float32(),
		float16.Frombits(output[2]).Float32(),
	), true
}

func (b *gbuffer) SampleNormal(cursor vec2.T) (vec3.T, bool) {
	if cursor.X < 0 || cursor.Y < 0 || cursor.X > float32(b.width) || cursor.Y > float32(b.height) {
		return vec3.Zero, false
	}

	offset := b.pixelOffset(cursor, b.normalBuf, 4)
	var output color.RGBA
	b.normalBuf.Memory().Read(offset, &output)

	if output.R == 0 && output.G == 0 && output.B == 0 {
		return vec3.Zero, false
	}

	return vec3.New(
		2*float32(output.R)/255-1,
		2*float32(output.G)/255-1,
		2*float32(output.B)/255-1,
	).Normalized(), true
}

func (p *gbuffer) Destroy() {
	p.positionBuf.Destroy()
	p.normalBuf.Destroy()
}

func (p *gbuffer) CopyBuffers() command.CommandFn {
	return func(b command.Buffer) {
		b.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			p.position.Image(),
			vk.ImageLayoutShaderReadOnlyOptimal,
			vk.ImageLayoutGeneral,
			vk.ImageAspectColorBit)

		b.CmdCopyImage(p.position.Image(), vk.ImageLayoutGeneral, p.positionBuf, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)

		b.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			p.normal.Image(),
			vk.ImageLayoutShaderReadOnlyOptimal,
			vk.ImageLayoutGeneral,
			vk.ImageAspectColorBit)

		b.CmdCopyImage(p.normal.Image(), vk.ImageLayoutGeneral, p.normalBuf, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
	}
}
