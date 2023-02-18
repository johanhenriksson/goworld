package pass

import (
	"image/color"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/x448/float16"
)

type GeometryBuffer interface {
	RenderTarget

	Diffuse() image.T
	Normal() image.T
	Position() image.T
	Destroy()

	SamplePosition(cursor vec2.T) (vec3.T, bool)
	SampleNormal(cursor vec2.T) (vec3.T, bool)
	RecordBufferCopy(command.Recorder)
}

type gbuffer struct {
	RenderTarget

	diffuse  image.T
	normal   image.T
	position image.T

	normalBuf   image.T
	positionBuf image.T
}

func NewGbuffer(
	app vulkan.App,
	rt RenderTarget,
) (GeometryBuffer, error) {
	diffuseFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	normalFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	positionFmt := core1_0.FormatR16G16B16A16SignedFloat
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment

	diffuse, err := image.New2D(app.Device(), "diffuse", rt.Width(), rt.Height(), diffuseFmt, usage)
	if err != nil {
		return nil, err
	}

	normal, err := image.New2D(app.Device(), "normal", rt.Width(), rt.Height(), normalFmt, usage|core1_0.ImageUsageTransferSrc)
	if err != nil {
		return nil, err
	}

	position, err := image.New2D(app.Device(), "position", rt.Width(), rt.Height(), positionFmt, usage|core1_0.ImageUsageTransferSrc)
	if err != nil {
		return nil, err
	}

	positionBuf, err := image.New(app.Device(), image.Args{
		Type:   core1_0.ImageType2D,
		Key:    "positionBuffer",
		Width:  rt.Width(),
		Height: rt.Height(),
		Format: positionFmt,
		Tiling: core1_0.ImageTilingLinear,
		Usage:  core1_0.ImageUsageTransferDst,
		Memory: core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
	})
	if err != nil {
		return nil, err
	}

	normalBuf, err := image.New(app.Device(), image.Args{
		Type:   core1_0.ImageType2D,
		Key:    "normalBuffer",
		Width:  rt.Width(),
		Height: rt.Height(),
		Format: normalFmt,
		Tiling: core1_0.ImageTilingLinear,
		Usage:  core1_0.ImageUsageTransferDst,
		Memory: core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
	})
	if err != nil {
		return nil, err
	}

	// move images to ImageLayoutGeneral to avoid errors on first copy
	worker := app.Transferer()
	worker.Queue(func(b command.Buffer) {
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, positionBuf, core1_0.ImageLayoutUndefined, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, normalBuf, core1_0.ImageLayoutUndefined, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)
	})
	worker.Submit(command.SubmitInfo{
		Marker: "GBufferInit",
	})

	return &gbuffer{
		RenderTarget: rt,

		diffuse:  diffuse,
		normal:   normal,
		position: position,

		positionBuf: positionBuf,
		normalBuf:   normalBuf,
	}, nil
}

func (b *gbuffer) Diffuse() image.T  { return b.diffuse }
func (b *gbuffer) Normal() image.T   { return b.normal }
func (b *gbuffer) Position() image.T { return b.position }

func (b *gbuffer) pixelOffset(pos vec2.T, img image.T, size int) int {
	x := int(pos.X * float32(img.Width()) / float32(b.Width()))
	y := int(pos.Y * float32(img.Height()) / float32(b.Height()))

	return size * (y*img.Width() + x)
}

func (b *gbuffer) SamplePosition(cursor vec2.T) (vec3.T, bool) {
	if cursor.X < 0 || cursor.Y < 0 || cursor.X > float32(b.Width()) || cursor.Y > float32(b.Height()) {
		return vec3.Zero, false
	}

	offset := b.pixelOffset(cursor, b.positionBuf, 8)
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
	if cursor.X < 0 || cursor.Y < 0 || cursor.X > float32(b.Width()) || cursor.Y > float32(b.Height()) {
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
	p.RenderTarget.Destroy()

	p.diffuse.Destroy()
	p.diffuse = nil

	p.normal.Destroy()
	p.normal = nil

	p.position.Destroy()
	p.position = nil

	p.positionBuf.Destroy()
	p.positionBuf = nil

	p.normalBuf.Destroy()
	p.normalBuf = nil
}

func (p *gbuffer) RecordBufferCopy(cmds command.Recorder) {
	cmds.Record(func(b command.Buffer) {
		b.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			p.position,
			core1_0.ImageLayoutShaderReadOnlyOptimal,
			core1_0.ImageLayoutGeneral,
			core1_0.ImageAspectColor)

		b.CmdCopyImage(p.position, core1_0.ImageLayoutGeneral, p.positionBuf, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)

		b.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			p.normal,
			core1_0.ImageLayoutShaderReadOnlyOptimal,
			core1_0.ImageLayoutGeneral,
			core1_0.ImageAspectColor)

		b.CmdCopyImage(p.normal, core1_0.ImageLayoutGeneral, p.normalBuf, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)
	})
}
