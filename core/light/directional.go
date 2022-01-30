package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type DirectionalArgs struct {
	Direction vec3.T
	Color     color.T
	Intensity float32
	Shadows   bool
}

type dirlight struct {
	object.Component

	DirectionalArgs
}

func NewDirectional(args DirectionalArgs) T {
	return &dirlight{
		Component:       object.NewComponent(),
		DirectionalArgs: args,
	}
}

func (lit *dirlight) LightDescriptor() Descriptor {
	desc := Descriptor{
		Type:       Directional,
		Position:   lit.Direction,
		Color:      lit.Color,
		Intensity:  lit.Intensity,
		Shadows:    lit.Shadows,
		Projection: mat4.Orthographic(-31, 60, -20, 40, -10, 50),
	}
	return desc
}
