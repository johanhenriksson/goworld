package vkrender

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
)

type GeometryBuffer interface {
	Diffuse(int) image.View
	Normal(int) image.View
	Position(int) image.View
	Depth(int) image.View
	Output(int) image.View
	Destroy()
}

type gbuffer struct {
	frames   int
	diffuse  []image.View
	normal   []image.View
	position []image.View
	output   []image.View
	depth    []image.View
}

func NewGbuffer(backend vulkan.T, pass renderpass.T, frames int) GeometryBuffer {
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

	return &gbuffer{
		frames:   frames,
		diffuse:  diffuse,
		normal:   normal,
		position: position,
		depth:    depth,
		output:   output,
	}
}

func (b *gbuffer) Diffuse(frame int) image.View  { return b.diffuse[frame%b.frames] }
func (b *gbuffer) Normal(frame int) image.View   { return b.normal[frame%b.frames] }
func (b *gbuffer) Position(frame int) image.View { return b.position[frame%b.frames] }
func (b *gbuffer) Depth(frame int) image.View    { return b.depth[frame%b.frames] }
func (b *gbuffer) Output(frame int) image.View   { return b.output[frame%b.frames] }

func (p *gbuffer) Destroy() {
	for i := range p.diffuse {
		p.diffuse[i].Destroy()
		p.normal[i].Destroy()
		p.position[i].Destroy()
		p.depth[i].Destroy()
	}
}
