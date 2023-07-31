package pass

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type GeometryBuffer interface {
	Diffuse() []image.T
	Normal() []image.T
	Position() []image.T
	Destroy()
}

type gbuffer struct {
	diffuse  []image.T
	normal   []image.T
	position []image.T
}

func NewGbuffer(device device.T, size vulkan.TargetSize) (GeometryBuffer, error) {
	frames, width, height := size.Frames, size.Width, size.Height
	diffuseFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	normalFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	positionFmt := core1_0.FormatR16G16B16A16SignedFloat
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment

	var err error
	diffuses := make([]image.T, frames)
	normals := make([]image.T, frames)
	positions := make([]image.T, frames)

	for i := 0; i < frames; i++ {
		diffuses[i], err = image.New2D(device, "diffuse", width, height, diffuseFmt, usage)
		if err != nil {
			return nil, err
		}

		normals[i], err = image.New2D(device, "normal", width, height, normalFmt, usage|core1_0.ImageUsageTransferSrc)
		if err != nil {
			return nil, err
		}

		positions[i], err = image.New2D(device, "position", width, height, positionFmt, usage|core1_0.ImageUsageTransferSrc)
		if err != nil {
			return nil, err
		}
	}

	return &gbuffer{
		diffuse:  diffuses,
		normal:   normals,
		position: positions,
	}, nil
}

func (b *gbuffer) Diffuse() []image.T  { return b.diffuse }
func (b *gbuffer) Normal() []image.T   { return b.normal }
func (b *gbuffer) Position() []image.T { return b.position }

func (b *gbuffer) pixelOffset(pos vec2.T, img image.T, size int) int {
	denormPos := pos.Mul(img.Size().XY())
	return size * (int(denormPos.Y)*img.Width() + int(denormPos.X))
}

func (p *gbuffer) Destroy() {
	for _, img := range p.diffuse {
		img.Destroy()
	}
	p.diffuse = nil

	for _, img := range p.normal {
		img.Destroy()
	}
	p.normal = nil

	for _, img := range p.position {
		img.Destroy()
	}
	p.position = nil
}
