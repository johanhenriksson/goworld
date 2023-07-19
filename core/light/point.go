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

type Point struct {
	object.Component

	Attenuation Attenuation

	Color     *object.Property[color.T]
	Range     *object.Property[float32]
	Intensity *object.Property[float32]
}

var _ T = &Point{}

func NewPoint(args PointArgs) *Point {
	return object.NewComponent(&Point{
		Attenuation: args.Attenuation,

		Color:     object.NewProperty(args.Color),
		Range:     object.NewProperty(args.Range),
		Intensity: object.NewProperty(args.Intensity),
	})
}

func (lit *Point) Name() string        { return "PointLight" }
func (lit *Point) Type() Type          { return TypePoint }
func (lit *Point) CastShadows() bool   { return false }
func (lit *Point) Cascades() []Cascade { return nil }

func (lit *Point) LightDescriptor(args render.Args, _ int) Descriptor {
	return Descriptor{
		Type:        TypePoint,
		Position:    vec4.Extend(lit.Transform().WorldPosition(), 0),
		Color:       lit.Color.Get(),
		Intensity:   lit.Intensity.Get(),
		Range:       lit.Range.Get(),
		Attenuation: lit.Attenuation,
	}
}
