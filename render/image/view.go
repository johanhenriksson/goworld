package image

import (
	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

var NilView View = &imgview{
	device: device.Nil,
	ptr:    vk.NullImageView,
	image:  Nil,
}

type View interface {
	device.Resource[vk.ImageView]

	Image() T
	Format() vk.Format
}

type imgview struct {
	ptr    vk.ImageView
	image  T
	format vk.Format
	device device.T
}

func (v *imgview) Ptr() vk.ImageView { return v.ptr }
func (v *imgview) Image() T          { return v.image }
func (v *imgview) Format() vk.Format { return v.format }

func (v *imgview) Destroy() {
	if v.ptr != vk.NullImageView {
		vk.DestroyImageView(v.device.Ptr(), v.ptr, nil)
		v.ptr = vk.NullImageView
	}
	v.device = device.Nil
	v.image = Nil
}
