package vkrender

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"

	vk "github.com/vulkan-go/vulkan"
)

type GeometryBuffer interface {
	Diffuse(int) vk_texture.T
	Normal(int) vk_texture.T
	Position(int) vk_texture.T
	Destroy()
}

type gbuffer struct {
	diffuse  []vk_texture.T
	normal   []vk_texture.T
	position []vk_texture.T
}

func NewGbuffer(backend vulkan.T, pass renderpass.T) GeometryBuffer {
	diffuse := make([]vk_texture.T, backend.Frames())
	normal := make([]vk_texture.T, backend.Frames())
	position := make([]vk_texture.T, backend.Frames())

	for i := 0; i < backend.Frames(); i++ {
		diffuseImage := pass.Attachment("diffuse").Image(i)
		diffuse[i] = vk_texture.FromImage(backend.Device(), diffuseImage, vk_texture.Args{
			Format: diffuseImage.Format(),
			Filter: vk.FilterLinear,
			Wrap:   vk.SamplerAddressModeRepeat,
		})

		normalImage := pass.Attachment("normal").Image(i)
		normal[i] = vk_texture.FromImage(backend.Device(), normalImage, vk_texture.Args{
			Format: normalImage.Format(),
			Filter: vk.FilterLinear,
			Wrap:   vk.SamplerAddressModeRepeat,
		})

		positionImage := pass.Attachment("position").Image(i)
		position[i] = vk_texture.FromImage(backend.Device(), positionImage, vk_texture.Args{
			Format: positionImage.Format(),
			Filter: vk.FilterLinear,
			Wrap:   vk.SamplerAddressModeRepeat,
		})
	}

	return &gbuffer{
		diffuse:  diffuse,
		normal:   normal,
		position: position,
	}
}

func (b *gbuffer) Diffuse(frame int) vk_texture.T {
	return b.diffuse[frame]
}

func (b *gbuffer) Normal(frame int) vk_texture.T {
	return b.normal[frame]
}

func (b *gbuffer) Position(frame int) vk_texture.T {
	return b.position[frame]
}

func (p *gbuffer) Destroy() {
	for i := range p.diffuse {
		p.diffuse[i].Destroy()
		p.normal[i].Destroy()
		p.position[i].Destroy()
	}
}
