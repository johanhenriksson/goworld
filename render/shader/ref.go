package shader

import "github.com/johanhenriksson/goworld/render/device"

type Ref interface {
	Key() string
	Version() int

	Load(*device.Device) *Shader
}

type ref struct {
	name string
}

func NewRef(name string) Ref {
	return &ref{name: name}
}

func (r *ref) Key() string {
	return r.name
}

func (r *ref) Version() int {
	return 1
}

func (r *ref) Load(dev *device.Device) *Shader {
	return New(dev, r.name)
}
