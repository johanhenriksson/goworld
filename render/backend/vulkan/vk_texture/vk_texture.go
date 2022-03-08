package vk_texture

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Sampler]
	Image() image.T
	View() image.View
}

type Args struct {
	Width  int
	Height int
	Format vk.Format
	Filter vk.Filter
	Wrap   vk.SamplerAddressMode
}

type vktexture struct {
	Args
	ptr    vk.Sampler
	device device.T
	image  image.T
	view   image.View
}

func New(device device.T, args Args) T {
	img := image.New2D(device,
		args.Width, args.Height, args.Format,
		vk.ImageUsageFlags(vk.ImageUsageSampledBit|vk.ImageUsageTransferDstBit))

	return FromImage(device, img, args)

}

func FromImage(device device.T, img image.T, args Args) T {
	view := img.View(args.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))

	info := vk.SamplerCreateInfo{
		SType: vk.StructureTypeSamplerCreateInfo,

		MinFilter:    args.Filter,
		MagFilter:    args.Filter,
		AddressModeU: args.Wrap,
		AddressModeV: args.Wrap,
		AddressModeW: args.Wrap,

		MipmapMode: vk.SamplerMipmapModeLinear,
	}

	var ptr vk.Sampler
	vk.CreateSampler(device.Ptr(), &info, nil, &ptr)

	return &vktexture{
		Args:   args,
		ptr:    ptr,
		device: device,
		image:  img,
		view:   view,
	}
}

func (t *vktexture) Ptr() vk.Sampler {
	return t.ptr
}

func (t *vktexture) Image() image.T   { return t.image }
func (t *vktexture) View() image.View { return t.view }

func (t *vktexture) Destroy() {
	vk.DestroySampler(t.device.Ptr(), t.ptr, nil)
	t.ptr = nil

	t.view.Destroy()
	t.image.Destroy()
}
