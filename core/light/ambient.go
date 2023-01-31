package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
)

type ambient struct {
	object.T
	Color     color.T
	Intensity float32
}

func (lit *ambient) LightDescriptor() Descriptor {
	return Descriptor{
		Type:      Ambient,
		Color:     lit.Color,
		Intensity: lit.Intensity,
	}
}

func (lit *ambient) Type() Type {
	return Ambient
}

func NewAmbient(clr color.T, intensity float32) T {
	return object.New(&ambient{
		Color:     clr,
		Intensity: intensity,
	})
}
