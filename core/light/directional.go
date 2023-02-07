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

	// create light view matrix looking at the center of the
	// camera frustum
	ldir := lit.Transform().Forward()
	position := center.Add(ldir.Scaled(-1))
	lview := mat4.LookAtLH(position, center)

	// project the camera frustum into light view space, and
	// find the minimum bounding box
	min, max := vec3.InfPos, vec3.InfNeg
	for _, corner := range frustumCorners {
		cornerLight := lview.TransformPoint(corner)
		min = vec3.Min(min, cornerLight)
		max = vec3.Max(max, cornerLight)
	}

	extents := math.Max(max.X-min.X, max.Y-min.Y) * 0.5
	snap := 2 * extents / float32(4096)
	extents = math.Snap(extents, snap)

	// create an orthographic projection from the minimum bounding box
	lproj := mat4.OrthographicLH_ZO(
		-extents, extents,
		-extents, extents, min.Z, max.Z)
	lvp := lproj.Mul(&lview)

	desc := Descriptor{
		Type:       Directional,
		Position:   vec4.Extend(ldir, 0),
		Color:      lit.Color,
		Intensity:  lit.Intensity,
		Projection: lproj,
		View:       lview,
		ViewProj:   lvp,
	}
	if lit.Shadows {
		desc.Shadows = 1
	}
	return desc
}
