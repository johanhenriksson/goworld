package texture

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vkerror"

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
	Aspect vk.ImageAspectFlags
	Usage  vk.ImageUsageFlags
}

type vktexture struct {
	Args
	ptr    vk.Sampler
	device device.T
	image  image.T
	view   image.View
}

func New(device device.T, args Args) (T, error) {
	if args.Usage == 0 {
		args.Usage = vk.ImageUsageFlags(vk.ImageUsageSampledBit | vk.ImageUsageTransferDstBit)
	}

	img, err := image.New2D(device, args.Width, args.Height, args.Format, args.Usage)
	if err != nil {
		return nil, err
	}

	tex, err := FromImage(device, img, args)
	if err != nil {
		img.Destroy()
		return nil, err
	}

	return tex, nil
}

func FromImage(device device.T, img image.T, args Args) (T, error) {
	if args.Aspect == 0 {
		args.Aspect = vk.ImageAspectFlags(vk.ImageAspectColorBit)
	}

	view, err := img.View(args.Format, args.Aspect)
	if err != nil {
		return nil, err
	}

	tex, err := FromView(device, view, args)
	if err != nil {
		// clean up
		view.Destroy()
		return nil, err
	}

	return tex, nil
}

func FromView(device device.T, view image.View, args Args) (T, error) {
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
	result := vk.CreateSampler(device.Ptr(), &info, nil, &ptr)
	if result != vk.Success {
		return nil, vkerror.FromResult(result)
	}

	return &vktexture{
		Args:   args,
		ptr:    ptr,
		device: device,
		image:  view.Image(),
		view:   view,
	}, nil
}

func (t *vktexture) Ptr() vk.Sampler {
	return t.ptr
}

func (t *vktexture) Image() image.T   { return t.image }
func (t *vktexture) View() image.View { return t.view }

func (t *vktexture) Destroy() {
	vk.DestroySampler(t.device.Ptr(), t.ptr, nil)
	t.ptr = vk.NullSampler

	t.view.Destroy()
	t.view = nil

	t.image.Destroy()
	t.image = nil

	t.device = nil
}
