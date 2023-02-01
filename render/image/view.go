package image

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type View interface {
	device.Resource[core1_0.ImageView]

	Image() T
	Format() core1_0.Format
}

type imgview struct {
	ptr    core1_0.ImageView
	image  T
	format core1_0.Format
	device device.T
}

func (v *imgview) Ptr() core1_0.ImageView { return v.ptr }
func (v *imgview) Image() T               { return v.image }
func (v *imgview) Format() core1_0.Format { return v.format }

func (v *imgview) Destroy() {
	if v.ptr != nil {
		v.ptr.Destroy(nil)
		v.ptr = nil
	}
	v.device = nil
	v.image = nil
}
