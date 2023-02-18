package attachment

import (
	"errors"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrAllocArrayExhausted = errors.New("image array allocator exhausted")

type Allocator interface {
	Alloc(
		device device.T,
		width, height int,
		format core1_0.Format,
		usage core1_0.ImageUsageFlags,
	) (image.T, error)
}

type alloc struct{}

var _ Allocator = &alloc{}

func (im *alloc) Alloc(
	device device.T,
	width, height int,
	format core1_0.Format,
	usage core1_0.ImageUsageFlags,
) (image.T, error) {
	return image.New2D(
		device,
		"",
		width, height, format,
		core1_0.ImageUsageFlags(usage),
	)
}

type imageArray struct {
	images []image.T
	next   int
}

func (im *imageArray) Alloc(
	device device.T,
	width, height int,
	format core1_0.Format,
	usage core1_0.ImageUsageFlags,
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

// FromImage returns an allocator that always returns a reference to the provided image.
func FromImage(img image.T) Allocator {
	return &imageRef{image: img}
}

type imageRef struct {
	image image.T
}

func (im *imageRef) Alloc(
	device device.T,
	width, height int,
	format core1_0.Format,
	usage core1_0.ImageUsageFlags,
) (image.T, error) {
	return im.image, nil
}
