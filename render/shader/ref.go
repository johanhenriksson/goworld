package shader

import (
	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/device"
)

type ref struct {
	name string
}

func Ref(name string) *ref {
	return &ref{name: name}
}

func (r *ref) Key() string {
	return r.name
}

func (r *ref) Version() int {
	return 1
}

func (r *ref) LoadShader(assets fs.Filesystem, dev *device.Device) *Shader {
	return New(dev, assets, r.name)
}
