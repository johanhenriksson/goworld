package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
)

type PointArgs struct {
	Attenuation Attenuation
	Color       color.T
	Range       float32
	Intensity   float32
}

type pointlight struct {
	object.Component

	PointArgs
}

func NewPoint(args PointArgs) T {
	return &pointlight{
		Component: object.NewComponent(),
		PointArgs: args,
	}
}

func (lit *pointlight) LightDescriptor() Descriptor {
	return Descriptor{
		Type:        Point,
		Position:    lit.Transform().WorldPosition(),
		Color:       lit.Color,
		Intensity:   lit.Intensity,
		Range:       lit.Range,
		Attenuation: lit.Attenuation,
	}
}
