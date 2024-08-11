package image

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ViewArray []*View

type View struct {
	ptr    core1_0.ImageView
	image  *Image
	format core1_0.Format
	device *device.Device
}

func (v *View) Ptr() core1_0.ImageView { return v.ptr }
func (v *View) Image() *Image          { return v.image }
func (v *View) Format() core1_0.Format { return v.format }

func (v *View) Destroy() {
	if v.ptr != nil {
		v.ptr.Destroy(nil)
		v.ptr = nil
	}
	v.device = nil
	v.image = nil
}
