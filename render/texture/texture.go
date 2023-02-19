package texture

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vkerror"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type T interface {
	device.Resource[core1_0.Sampler]
	Image() image.T
	View() image.View
	Size() vec3.T
}

type Args struct {
	Key    string
	Width  int
	Height int
	Format core1_0.Format
	Filter core1_0.Filter
	Wrap   core1_0.SamplerAddressMode
	Aspect core1_0.ImageAspectFlags
	Usage  core1_0.ImageUsageFlags
}

type vktexture struct {
	Args
	ptr    core1_0.Sampler
	device device.T
	image  image.T
	view   image.View
}

func New(device device.T, args Args) (T, error) {
	if args.Key == "" {
		panic("texture must have a key")
	}
	if args.Usage == 0 {
		args.Usage = core1_0.ImageUsageFlags(core1_0.ImageUsageSampled | core1_0.ImageUsageTransferDst)
	}

	img, err := image.New2D(device, args.Key, args.Width, args.Height, args.Format, args.Usage)
	if err != nil {
		return nil, err
	}

	if args.Key != "" {
		device.SetDebugObjectName(driver.VulkanHandle(img.Ptr().Handle()),
			core1_0.ObjectTypeImage, args.Key)
	}

	tex, err := FromImage(device, img, args)
	if err != nil {
		img.Destroy()
		return nil, err
	}

	return tex, nil
}

func FromImage(device device.T, img image.T, args Args) (T, error) {
	if args.Key == "" {
		args.Key = img.Key()
	}
	if args.Aspect == 0 {
		args.Aspect = core1_0.ImageAspectFlags(core1_0.ImageAspectColor)
	}
	if args.Format == 0 {
		args.Format = img.Format()
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
	if args.Key == "" {
		panic("texture must have a key")
	}
	info := core1_0.SamplerCreateInfo{
		MinFilter:    args.Filter,
		MagFilter:    args.Filter,
		AddressModeU: args.Wrap,
		AddressModeV: args.Wrap,
		AddressModeW: args.Wrap,

		MipmapMode: core1_0.SamplerMipmapModeLinear,
	}

	ptr, result, err := device.Ptr().CreateSampler(nil, info)
	if err != nil {
		return nil, err
	}
	if result != core1_0.VKSuccess {
		return nil, vkerror.FromResult(result)
	}

	if args.Key != "" {
		device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()),
			core1_0.ObjectTypeSampler, args.Key)
	}

	return &vktexture{
		Args:   args,
		ptr:    ptr,
		device: device,
		image:  view.Image(),
		view:   view,
	}, nil
}

func (t *vktexture) Ptr() core1_0.Sampler {
	return t.ptr
}

func (t *vktexture) Image() image.T   { return t.image }
func (t *vktexture) View() image.View { return t.view }
func (t *vktexture) Size() vec3.T     { return t.image.Size() }

func (t *vktexture) Destroy() {
	t.ptr.Destroy(nil)
	t.ptr = nil

	t.view.Destroy()
	t.view = nil

	t.image.Destroy()
	t.image = nil

	t.device = nil
}
