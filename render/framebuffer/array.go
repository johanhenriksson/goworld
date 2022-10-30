package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass"
)

type Array []T

func NewArray(count int, device device.T, width, height int, pass renderpass.T) (Array, error) {
	var err error
	array := make(Array, count)
	for i := range array {
		array[i], err = New(device, width, height, pass)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func (a Array) Destroy() {
	for i, fbuf := range a {
		fbuf.Destroy()
		a[i] = nil
	}
}
