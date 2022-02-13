package image

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

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

func (v *imgview) Ptr() vk.ImageView {
	return v.ptr
}

func (v *imgview) Image() T {
	return v.image
}

func (v *imgview) Format() vk.Format {
	return v.format
}

func (v *imgview) Destroy() {
	vk.DestroyImageView(v.device.Ptr(), v.ptr, nil)
	v.ptr = nil
}
