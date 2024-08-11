package texture

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vkerror"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Array []*Texture

type Args struct {
	Filter  Filter
	Wrap    Wrap
	Aspect  core1_0.ImageAspectFlags
	Usage   core1_0.ImageUsageFlags
	Border  core1_0.BorderColor
	Mipmaps bool
}

type Texture struct {
	Args
	ptr    core1_0.Sampler
	key    string
	device *device.Device
	image  *image.Image
	view   *image.View
}

func New(device *device.Device, key string, width, height int, format core1_0.Format, args Args) (*Texture, error) {
	if key == "" {
		panic("texture must have a key")
	}
	args.Usage |= core1_0.ImageUsageFlags(core1_0.ImageUsageSampled | core1_0.ImageUsageTransferSrc | core1_0.ImageUsageTransferDst)

	img, err := image.New2D(device, key, width, height, format, args.Mipmaps, args.Usage)
	if err != nil {
		return nil, err
	}

	device.SetDebugObjectName(driver.VulkanHandle(img.Ptr().Handle()),
		core1_0.ObjectTypeImage, key)

	tex, err := FromImage(device, key, img, args)
	if err != nil {
		img.Destroy()
		return nil, err
	}

	return tex, nil
}

func FromImage(device *device.Device, key string, img *image.Image, args Args) (*Texture, error) {
	if key == "" {
		key = img.Key()
	}
	if args.Aspect == 0 {
		args.Aspect = core1_0.ImageAspectFlags(core1_0.ImageAspectColor)
	}

	view, err := img.View(img.Format(), args.Aspect)
	if err != nil {
		return nil, err
	}

	tex, err := FromView(device, key, view, args)
	if err != nil {
		// clean up
		view.Destroy()
		return nil, err
	}

	return tex, nil
}

func FromView(device *device.Device, key string, view *image.View, args Args) (*Texture, error) {
	if key == "" {
		panic("texture must have a key")
	}
	info := core1_0.SamplerCreateInfo{
		MinFilter:    core1_0.Filter(args.Filter),
		MagFilter:    core1_0.Filter(args.Filter),
		AddressModeU: core1_0.SamplerAddressMode(args.Wrap),
		AddressModeV: core1_0.SamplerAddressMode(args.Wrap),
		AddressModeW: core1_0.SamplerAddressMode(args.Wrap),
		BorderColor:  args.Border,

		MipmapMode: core1_0.SamplerMipmapModeNearest,
		MinLod:     0,
		MaxLod:     float32(view.Image().MipLevels()),
		MipLodBias: 0,
	}

	ptr, result, err := device.Ptr().CreateSampler(nil, info)
	if err != nil {
		return nil, err
	}
	if result != core1_0.VKSuccess {
		return nil, vkerror.FromResult(result)
	}

	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()),
		core1_0.ObjectTypeSampler, key)

	return &Texture{
		Args:   args,
		key:    key,
		ptr:    ptr,
		device: device,
		image:  view.Image(),
		view:   view,
	}, nil
}

func (t *Texture) Ptr() core1_0.Sampler {
	return t.ptr
}

func (t *Texture) Key() string         { return t.key }
func (t *Texture) Image() *image.Image { return t.image }
func (t *Texture) View() *image.View   { return t.view }
func (t *Texture) Size() vec3.T        { return t.image.Size() }

func (t *Texture) Destroy() {
	t.ptr.Destroy(nil)
	t.ptr = nil

	t.view.Destroy()
	t.view = nil

	t.image.Destroy()
	t.image = nil

	t.device = nil
}
