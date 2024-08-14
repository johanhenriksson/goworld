package pass

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type GeometryBuffer interface {
	Width() int
	Height() int
	Frames() int
	Diffuse() image.Array
	Normal() image.Array
	Position() image.Array
	Destroy()
}

type gbuffer struct {
	diffuse  image.Array
	normal   image.Array
	position image.Array
	width    int
	height   int
}

func NewGbuffer(device *device.Device, size engine.TargetSize) (GeometryBuffer, error) {
	frames, width, height := size.Frames, size.Width, size.Height
	diffuseFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	normalFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	positionFmt := core1_0.FormatR32G32B32A32SignedFloat
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment

	var err error
	diffuses := make(image.Array, frames)
	normals := make(image.Array, frames)
	positions := make(image.Array, frames)

	for i := 0; i < frames; i++ {
		diffuses[i], err = image.New2D(device, "diffuse", width, height, diffuseFmt, false, usage)
		if err != nil {
			return nil, err
		}

		normals[i], err = image.New2D(device, "normal", width, height, normalFmt, false, usage|core1_0.ImageUsageTransferSrc)
		if err != nil {
			return nil, err
		}

		positions[i], err = image.New2D(device, "position", width, height, positionFmt, false, usage|core1_0.ImageUsageTransferSrc)
		if err != nil {
			return nil, err
		}
	}

	return &gbuffer{
		diffuse:  diffuses,
		normal:   normals,
		position: positions,
		width:    width,
		height:   height,
	}, nil
}

func (b *gbuffer) Width() int            { return b.width }
func (b *gbuffer) Height() int           { return b.height }
func (b *gbuffer) Frames() int           { return len(b.diffuse) }
func (b *gbuffer) Diffuse() image.Array  { return b.diffuse }
func (b *gbuffer) Normal() image.Array   { return b.normal }
func (b *gbuffer) Position() image.Array { return b.position }

func (b *gbuffer) pixelOffset(pos vec2.T, img *image.Image, size int) int {
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
