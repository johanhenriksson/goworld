package pass

import (
	"image/color"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/x448/float16"
)

type GeometryBuffer interface {
	Diffuse() image.T
	Normal() image.T
	Position() image.T
	Destroy()

	NormalBuf() image.T

	SamplePosition(point vec2.T) (vec3.T, bool)
	SampleNormal(point vec2.T) (vec3.T, bool)
	RecordBufferCopy(command.Recorder)
}

type gbuffer struct {
	diffuse  image.T
	normal   image.T
	position image.T

	normalBuf   image.T
	positionBuf image.T
}

func NewGbuffer(
	device device.T,
	width, height int,
) (GeometryBuffer, error) {
	diffuseFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	normalFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	positionFmt := core1_0.FormatR16G16B16A16SignedFloat
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment

	diffuse, err := image.New2D(device, "diffuse", width, height, diffuseFmt, usage)
	if err != nil {
		return nil, err
	}

	normal, err := image.New2D(device, "normal", width, height, normalFmt, usage|core1_0.ImageUsageTransferSrc)
	if err != nil {
		return nil, err
	}

	position, err := image.New2D(device, "position", width, height, positionFmt, usage|core1_0.ImageUsageTransferSrc)
	if err != nil {
		return nil, err
	}

	positionBuf, err := image.New(device, image.Args{
		Type:   core1_0.ImageType2D,
		Key:    "positionBuffer",
		Width:  width,
		Height: height,
		Format: positionFmt,
		Tiling: core1_0.ImageTilingLinear,
		Usage:  core1_0.ImageUsageTransferDst,
		Memory: core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
	})
	if err != nil {
		return nil, err
	}

	normalBuf, err := image.New(device, image.Args{
		Type:   core1_0.ImageType2D,
		Key:    "normalBuffer",
		Width:  width,
		Height: height,
		Format: normalFmt,
		Tiling: core1_0.ImageTilingLinear,
		Usage:  core1_0.ImageUsageTransferDst,
		Memory: core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
	})
	if err != nil {
		return nil, err
	}

	return &gbuffer{
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

func (b *gbuffer) NormalBuf() image.T { return b.normalBuf }

func (b *gbuffer) pixelOffset(pos vec2.T, img image.T, size int) int {
	denormPos := pos.Mul(img.Size().XY())
	return size * (int(denormPos.Y)*img.Width() + int(denormPos.X))
}

func (b *gbuffer) SamplePosition(point vec2.T) (vec3.T, bool) {
	if point.X < 0 || point.Y < 0 || point.X > 1 || point.Y > 1 {
		return vec3.Zero, false
	}

	offset := b.pixelOffset(point, b.positionBuf, 8)
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

func (b *gbuffer) SampleNormal(point vec2.T) (vec3.T, bool) {
	if point.X < 0 || point.Y < 0 || point.X > 1 || point.Y > 1 {
		return vec3.Zero, false
	}

	offset := b.pixelOffset(point, b.normalBuf, 4)
	var output color.RGBA
	b.normalBuf.Memory().Read(offset, &output)

	if output.R == 0 && output.G == 0 && output.B == 0 {
		return vec3.Zero, false
	}

	// unpack normal
	return vec3.New(
		2*float32(output.R)/255-1,
		2*float32(output.G)/255-1,
		2*float32(output.B)/255-1,
	).Normalized(), true
}

func (p *gbuffer) Destroy() {
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
		// transfer position buffer layout to TransferSrcOptimal
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, p.position, core1_0.ImageLayoutShaderReadOnlyOptimal, core1_0.ImageLayoutTransferSrcOptimal, core1_0.ImageAspectColor)

		// transfer position copy buffer layout to General
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, p.positionBuf, core1_0.ImageLayoutUndefined, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)

		// copy position buffer
		b.CmdCopyImage(p.position, core1_0.ImageLayoutTransferSrcOptimal, p.positionBuf, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)

		// transfer normal buffer layout to TransferSrcOptimal
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, p.normal, core1_0.ImageLayoutShaderReadOnlyOptimal, core1_0.ImageLayoutTransferSrcOptimal, core1_0.ImageAspectColor)

		// transfer normal copy buffer layout to General
		b.CmdImageBarrier(core1_0.PipelineStageTopOfPipe, core1_0.PipelineStageTransfer, p.normalBuf, core1_0.ImageLayoutUndefined, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)

		// copy normal buffer
		b.CmdCopyImage(p.normal, core1_0.ImageLayoutTransferSrcOptimal, p.normalBuf, core1_0.ImageLayoutGeneral, core1_0.ImageAspectColor)
	})
}
