package attachment

import (
	"errors"
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrArrayExhausted = errors.New("image array allocator exhausted")

type Image interface {
	Format() core1_0.Format
	Next(device *device.Device, name string, width, height int) (*image.Image, bool, error)
}

type alloc struct {
	key    string
	format core1_0.Format
	usage  core1_0.ImageUsageFlags
}

var _ Image = &alloc{}

func (im *alloc) Format() core1_0.Format {
	return im.format
}

func (im *alloc) Next(
	device *device.Device,
	name string,
	width, height int,
) (*image.Image, bool, error) {
	key := fmt.Sprintf("%s-%s", name, im.key)
	log.Println("attachment alloc", key)
	img, err := image.New2D(
		device,
		key,
		width, height,
		im.format, false, im.usage,
	)
	return img, true, err
}

func NewImage(key string, format core1_0.Format, usage core1_0.ImageUsageFlags) Image {
	return &alloc{
		key:    key,
		format: format,
		usage:  usage,
	}
}

type imageArray struct {
	images image.Array
	next   int
}

func (im *imageArray) Format() core1_0.Format {
	return im.images[0].Format()
}

func (im *imageArray) Next(
	device *device.Device,
	name string,
	width, height int,
) (*image.Image, bool, error) {
	if im.next >= len(im.images) {
		return nil, false, ErrArrayExhausted
	}
	img := im.images[im.next]
	im.next++
	return img, false, nil
}

func FromImageArray(images image.Array) Image {
	return &imageArray{
		images: images,
		next:   0,
	}
}

// FromImage returns an allocator that always returns a reference to the provided image.
func FromImage(img *image.Image) Image {
	return &imageRef{image: img}
}

type imageRef struct {
	image *image.Image
}

func (im *imageRef) Format() core1_0.Format {
	return im.image.Format()
}

func (im *imageRef) Next(
	device *device.Device,
	name string,
	width, height int,
) (*image.Image, bool, error) {
	return im.image, false, nil
}
