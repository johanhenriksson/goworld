package attachment

import (
	"errors"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	vk "github.com/vulkan-go/vulkan"
)

var ErrAllocArrayExhausted = errors.New("image array allocator exhausted")

type Allocator interface {
	Alloc(
		device device.T,
		width, height int,
		format vk.Format,
		usage vk.ImageUsageFlagBits,
	) (image.T, error)
}

type alloc struct{}

var _ Allocator = &alloc{}

func (im *alloc) Alloc(
	device device.T,
	width, height int,
	format vk.Format,
	usage vk.ImageUsageFlagBits,
) (image.T, error) {
	return image.New2D(
		device,
		width, height, format,
		vk.ImageUsageFlags(usage),
	)
}

type imageArray struct {
	images []image.T
	next   int
}

func (im *imageArray) Alloc(
	device device.T,
	width, height int,
	format vk.Format,
	usage vk.ImageUsageFlagBits,
) (image.T, error) {
	if im.next >= len(im.images) {
		return nil, ErrAllocArrayExhausted
	}
	img := im.images[im.next]
	im.next++
	return img, nil
}

func FromImageArray(images []image.T) Allocator {
	return &imageArray{
		images: images,
		next:   0,
	}
}

func FromSwapchain(swap swapchain.T) Allocator {
	return FromImageArray(swap.Images())
}
