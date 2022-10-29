package vkrender

import (
	"image/color"
	"log"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"

	vk "github.com/vulkan-go/vulkan"
	"github.com/x448/float16"
)

type GeometryBuffer interface {
	engine.BufferOutput

	Diffuse(int) image.View
	Normal(int) image.View
	Position(int) image.View
	Depth(int) image.View
	Output(int) image.View
	Destroy()

	CopyBuffers(command.Worker, int, command.Wait)
}

type gbuffer struct {
	frames int
	width  int
	height int

	diffuse  []image.View
	normal   []image.View
	position []image.View
	output   []image.View
	depth    []image.View

	normalBuf   image.T
	positionBuf image.T
}

func NewGbuffer(backend vulkan.T, pass renderpass.T) GeometryBuffer {
	frames := pass.Frames()
	log.Println("creating gbuffer with", frames, "frames")

	diffuse := make([]image.View, frames)
	normal := make([]image.View, frames)
	position := make([]image.View, frames)
	output := make([]image.View, frames)
	depth := make([]image.View, frames)

	for i := 0; i < frames; i++ {
		diffuse[i] = pass.Attachment("diffuse").View(i)
		normal[i] = pass.Attachment("normal").View(i)
		position[i] = pass.Attachment("position").View(i)
		output[i] = pass.Attachment("output").View(i)
		depth[i] = pass.Depth().View(i)
	}

	positionBuf := image.New(backend.Device(), image.Args{
		Type:   vk.ImageType2d,
		Width:  position[0].Image().Width(),
		Height: position[0].Image().Height(),
		Format: position[0].Format(),
		Tiling: vk.ImageTilingLinear,
		Usage:  vk.ImageUsageTransferDstBit,
		Memory: vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
	})

	normalBuf := image.New(backend.Device(), image.Args{
		Type:   vk.ImageType2d,
		Width:  normal[0].Image().Width(),
		Height: normal[0].Image().Height(),
		Format: normal[0].Format(),
		Tiling: vk.ImageTilingLinear,
		Usage:  vk.ImageUsageTransferDstBit,
		Memory: vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
	})

	// move images to ImageLayoutGeneral to avoid errors on first copy
	worker := backend.Transferer()
	worker.Queue(func(b command.Buffer) {
		b.CmdImageBarrier(vk.PipelineStageTopOfPipeBit, vk.PipelineStageTransferBit, positionBuf, vk.ImageLayoutUndefined, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
		b.CmdImageBarrier(vk.PipelineStageTopOfPipeBit, vk.PipelineStageTransferBit, normalBuf, vk.ImageLayoutUndefined, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
	})
	worker.Submit(command.SubmitInfo{})
	worker.Wait()

	return &gbuffer{
		frames: frames,
		width:  backend.Width(),
		height: backend.Height(),

		diffuse:  diffuse,
		normal:   normal,
		position: position,
		depth:    depth,
		output:   output,

		positionBuf: positionBuf,
		normalBuf:   normalBuf,
	}
}

func (b *gbuffer) Diffuse(frame int) image.View  { return b.diffuse[frame%b.frames] }
func (b *gbuffer) Normal(frame int) image.View   { return b.normal[frame%b.frames] }
func (b *gbuffer) Position(frame int) image.View { return b.position[frame%b.frames] }
func (b *gbuffer) Depth(frame int) image.View    { return b.depth[frame%b.frames] }
func (b *gbuffer) Output(frame int) image.View   { return b.output[frame%b.frames] }

func (b *gbuffer) pixelOffset(pos vec2.T, img image.T, size int) int {
	x := int(pos.X * float32(img.Width()) / float32(b.width))
	y := int(pos.Y * float32(img.Height()) / float32(b.height))

	return size * (y*img.Width() + x)
}

func (b *gbuffer) SamplePosition(cursor vec2.T) (vec3.T, bool) {
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
	for i := range p.diffuse {
		p.diffuse[i].Destroy()
		p.normal[i].Destroy()
		p.position[i].Destroy()
		p.depth[i].Destroy()
	}

	p.positionBuf.Destroy()
	p.normalBuf.Destroy()
}

func (p *gbuffer) CopyBuffers(worker command.Worker, frame int, wait command.Wait) {
	worker.Queue(func(b command.Buffer) {
		b.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			p.position[frame%p.frames].Image(),
			vk.ImageLayoutShaderReadOnlyOptimal,
			vk.ImageLayoutGeneral,
			vk.ImageAspectColorBit)

		b.CmdCopyImage(p.position[frame%p.frames].Image(), vk.ImageLayoutGeneral, p.positionBuf, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)

		b.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			p.normal[frame%p.frames].Image(),
			vk.ImageLayoutShaderReadOnlyOptimal,
			vk.ImageLayoutGeneral,
			vk.ImageAspectColorBit)

		b.CmdCopyImage(p.normal[frame%p.frames].Image(), vk.ImageLayoutGeneral, p.normalBuf, vk.ImageLayoutGeneral, vk.ImageAspectColorBit)
	})
	worker.Submit(command.SubmitInfo{
		Wait: []command.Wait{wait},
	})
}
