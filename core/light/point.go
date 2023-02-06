package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type PointArgs struct {
	Attenuation Attenuation
	Color       color.T
	Range       float32
	Intensity   float32
}

type pointlight struct {
	object.T

	PointArgs
}

func NewPoint(args PointArgs) T {
	return object.New(&pointlight{
		PointArgs: args,
	})
}

func (lit *pointlight) Name() string { return "PointLight" }
func (lit *pointlight) Type() Type   { return Point }

func (lit *pointlight) LightDescriptor(args render.Args) Descriptor {
	return Descriptor{
		Type:        Point,
		Position:    vec4.Extend(lit.Transform().WorldPosition(), 0),
		Color:       lit.Color,
		Intensity:   lit.Intensity,
		Range:       lit.Range,
		Attenuation: lit.Attenuation,
	}
}
