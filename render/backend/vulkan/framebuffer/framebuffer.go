package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
)

type T interface {
}

type framebuf struct {
	device device.T
	depth  image.View
	color  image.View
}

func New(device device.T) {

}
