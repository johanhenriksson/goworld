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
	Color     color.T
	Intensity float32
	Shadows   bool
}

type dirlight struct {
	object.T
	args DirectionalArgs
}

func NewDirectional(args DirectionalArgs) T {
	return object.New(&dirlight{
		args: args,
	})
}

func (lit *dirlight) Name() string  { return "DirectionalLight" }
func (lit *dirlight) Type() Type    { return Directional }
func (lit *dirlight) Shadows() bool { return lit.args.Shadows }

func (lit *dirlight) LightDescriptor(args render.Args) Descriptor {
	frustumCorners := []vec3.T{
		vec3.New(-1, 1, -1),  // NTL
		vec3.New(1, 1, -1),   // NTR
		vec3.New(-1, -1, -1), // NBL
		vec3.New(1, -1, -1),  // NBR
		vec3.New(-1, 1, 1),   // FTL
		vec3.New(1, 1, 1),    // FTR
		vec3.New(-1, -1, 1),  // FBL
		vec3.New(1, -1, 1),   // FBR
	}

	center := vec3.Zero
	for i, corner := range frustumCorners {
		cornerWorld := args.VPInv.TransformPoint(corner)
		frustumCorners[i] = cornerWorld
		center = center.Add(cornerWorld)
	}
	center = center.Scaled(float32(1) / 8)

	radius := float32(0)
	for _, corner := range frustumCorners {
		distance := vec3.Distance(corner, center)
		radius = math.Max(radius, distance)
	}
	radius = math.Snap(radius, 16)

	// create light view matrix looking at the center of the
	// camera frustum
	ldir := lit.Transform().Forward()
	position := center.Sub(ldir.Scaled(radius))
	lview := mat4.LookAt(position, center, vec3.UnitY)

	lproj := mat4.Orthographic(
		-radius, radius,
		-radius, radius,
		0, 2*radius)

	lvp := lproj.Mul(&lview)

	return Descriptor{
		Type:       Directional,
		Position:   vec4.Extend(ldir, 0),
		Color:      lit.args.Color,
		Intensity:  lit.args.Intensity,
		Projection: lproj,
		View:       lview,
		ViewProj:   lvp,
	}
}
