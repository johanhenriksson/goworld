package image

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type View interface {
	device.Resource
	Image() T
	Ptr() vk.ImageView
}

type imgview struct {
	ptr    vk.ImageView
	image  T
	device device.T
}

func (v *imgview) Ptr() vk.ImageView {
	return v.ptr
}

func (v *imgview) Image() T {
	return v.image
}

func (v *imgview) Destroy() {
	vk.DestroyImageView(v.device.Ptr(), v.ptr, nil)
	v.ptr = nil
}
