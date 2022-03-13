package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec4"
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

func (lit *pointlight) Name() string { return "PointLight" }

func (lit *pointlight) LightDescriptor() Descriptor {
	return Descriptor{
		Type:        Point,
		Position:    vec4.Extend(lit.Transform().WorldPosition(), 0),
		Color:       lit.Color,
		Intensity:   lit.Intensity,
		Range:       lit.Range,
		Attenuation: lit.Attenuation,
	}
}
