package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type Ambient struct {
	object.Component
	Color     color.T
	Intensity float32
}

var _ T = &Ambient{}

func NewAmbient(clr color.T, intensity float32) T {
	return object.NewComponent(&Ambient{
		Color:     clr,
		Intensity: intensity,
	})
}

func (lit *Ambient) LightDescriptor(args render.Args, _ int) Descriptor {
	return Descriptor{
		Type:      TypeAmbient,
		Color:     lit.Color,
		Intensity: lit.Intensity,
	}
}

func (lit *Ambient) CastShadows() bool { return false }
func (lit *Ambient) Type() Type {
	return TypeAmbient
}
func (lit *Ambient) Cascades() []Cascade { return nil }
