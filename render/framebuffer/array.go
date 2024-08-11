package framebuffer

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass"
)

type Array []*Framebuffer

func NewArray(count int, device *device.Device, name string, width, height int, pass *renderpass.Renderpass) (Array, error) {
	var err error
	array := make(Array, count)
	for i := range array {
		array[i], err = New(device, fmt.Sprintf("%s[%d]", name, i), width, height, pass)
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
