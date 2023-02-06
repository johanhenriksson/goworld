package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type DirectionalArgs struct {
	Direction vec3.T
	Color     color.T
	Intensity float32
	Shadows   bool
}

type dirlight struct {
	object.T

	DirectionalArgs
}

func NewDirectional(args DirectionalArgs) T {
	return object.New(&dirlight{
		DirectionalArgs: args,
	})
}

func (lit *dirlight) Name() string { return "DirectionalLight" }
func (lit *dirlight) Type() Type   { return Directional }

func (lit *dirlight) LightDescriptor(args render.Args) Descriptor {
	position := lit.Direction.Scaled(-1).Normalized() // turn direction into a position

	width := float32(args.Viewport.Width)
	height := float32(args.Viewport.Height)
	ar := width / height
	tanHalfHfov := math.Tan(math.DegToRad(args.Fov * ar / 2))
	tanHalfVfov := math.Tan(math.DegToRad(args.Fov / 2))

	xn := args.Near * tanHalfHfov
	xf := args.Far * tanHalfHfov
	yn := args.Near * tanHalfVfov
	yf := args.Far * tanHalfVfov

	frustumCorners := []vec4.T{
		// near face
		vec4.New(xn, yn, args.Near, 1.0),
		vec4.New(-xn, yn, args.Near, 1.0),
		vec4.New(xn, -yn, args.Near, 1.0),
		vec4.New(-xn, -yn, args.Near, 1.0),

		// far face
		vec4.New(xf, yf, args.Far, 1.0),
		vec4.New(-xf, yf, args.Far, 1.0),
		vec4.New(xf, -yf, args.Far, 1.0),
		vec4.New(-xf, -yf, args.Far, 1.0),
	}

	// create light view matrix
	lv := mat4.LookAt(position, vec3.Zero)

	min, max := vec3.InfPos, vec3.InfNeg
	for _, corner := range frustumCorners {
		cornerWorld := args.ViewInv.VMul(corner)
		cornerLight := lv.VMul(cornerWorld).XYZ()
		min = vec3.Min(min, cornerLight)
		max = vec3.Max(max, cornerLight)
	}

	// these calculations will need to know about the camera frustum later
	lp := mat4.OrthographicVK(min.X, max.X, min.Y, max.Y, min.Z, max.Z)
	lvp := lp.Mul(&lv)

	desc := Descriptor{
		Type:       Directional,
		Position:   vec4.Extend(position, 0),
		Color:      lit.Color,
		Intensity:  lit.Intensity,
		Projection: lp,
		View:       lv,
		ViewProj:   lvp,
	}
	if lit.Shadows {
		desc.Shadows = 1
	}
	return desc
}
