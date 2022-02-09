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

func (lit *dirlight) Name() string { return "DirectionalLight" }

func (lit *dirlight) LightDescriptor() Descriptor {
	position := lit.Direction.Scaled(-1).Normalized() // turn direction into a position

	// these calculations will need to know about the camera frustum later
	lp := mat4.Orthographic(-16, 20, -10, 30, -20, 30)
	lv := mat4.LookAt(position, vec3.Zero)
	lvp := lp.Mul(&lv)

	desc := Descriptor{
		Type:       Directional,
		Position:   position,
		Color:      lit.Color,
		Intensity:  lit.Intensity,
		Shadows:    lit.Shadows,
		Projection: lp,
		View:       lv,
		ViewProj:   lvp,
	}
	return desc
}
