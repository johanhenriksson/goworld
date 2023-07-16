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
	Cascades  int
}

type Cascade struct {
	View      mat4.T
	Proj      mat4.T
	ViewProj  mat4.T
	NearSplit float32
	FarSplit  float32
}

type dirlight struct {
	object.Component
	args     DirectionalArgs
	cascades []Cascade
}

func NewDirectional(args DirectionalArgs) T {
	return object.NewComponent(&dirlight{
		args:     args,
		cascades: make([]Cascade, args.Cascades),
	})
}

func (lit *dirlight) Name() string        { return "DirectionalLight" }
func (lit *dirlight) Type() Type          { return Directional }
func (lit *dirlight) Shadows() bool       { return lit.args.Shadows }
func (lit *dirlight) Cascades() []Cascade { return lit.cascades }

func farSplitDist(cascade, cascades int, near, far float32) float32 {
	cascadeSplitLambda := float32(0.90)
	clipRange := far - near
	minZ := near
	maxZ := near + clipRange

	rnge := maxZ - minZ
	ratio := maxZ / minZ

	// Calculate split depths based on view camera frustum
	// Based on method presented in https://developer.nvidia.com/gpugems/GPUGems3/gpugems3_ch10.html
	p := (float32(cascade) + 1) / float32(cascades)
	log := minZ * math.Pow(ratio, p)
	uniform := minZ + rnge*p
	d := cascadeSplitLambda*(log-uniform) + uniform
	return (d - near) / clipRange
}

func nearSplitDist(cascade, cascades int, near, far float32) float32 {
	if cascade == 0 {
		return 0
	}
	return farSplitDist(cascade-1, cascades, near, far)
}

func (lit *dirlight) PreDraw(args render.Args, scene object.G) error {
	lit.updateCascades(args)
	return nil
}

func (lit *dirlight) LightDescriptor(args render.Args, cascade int) Descriptor {
	ldir := lit.Transform().Forward()
	return Descriptor{
		Type:       Directional,
		Position:   vec4.Extend(ldir, 0),
		Color:      lit.args.Color,
		Intensity:  lit.args.Intensity,
		View:       lit.cascades[cascade].View,
		Projection: lit.cascades[cascade].Proj,
		ViewProj:   lit.cascades[cascade].ViewProj,
	}
}

func (lit *dirlight) updateCascades(args render.Args) {
	for i := 0; i < lit.args.Cascades; i++ {
		lit.cascades[i] = lit.calculateCascade(args, i, lit.args.Cascades)
	}
}

func (lit *dirlight) calculateCascade(args render.Args, cascade, cascades int) Cascade {
	texSize := float32(2048)

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

	// transform frustum into world space
	for i, corner := range frustumCorners {
		frustumCorners[i] = args.VPInv.TransformPoint(corner)
	}

	// squash
	nearSplit := nearSplitDist(cascade, cascades, args.Near, args.Far)
	farSplit := farSplitDist(cascade, cascades, args.Near, args.Far)
	for i := 0; i < 4; i++ {
		dist := frustumCorners[i+4].Sub(frustumCorners[i])
		frustumCorners[i] = frustumCorners[i].Add(dist.Scaled(nearSplit))
		frustumCorners[i+4] = frustumCorners[i].Add(dist.Scaled(farSplit))
	}

	// calculate frustum center
	center := vec3.Zero
	for _, corner := range frustumCorners {
		center = center.Add(corner)
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
		-radius-0.01, radius+0.01,
		-radius-0.01, radius+0.01,
		0, 2*radius)

	lvp := lproj.Mul(&lview)

	// round the center of the lights projection to the nearest texel
	origin := lvp.TransformPoint(vec3.New(0, 0, 0)).Scaled(texSize / 2.0)
	offset := origin.Round().Sub(origin)
	offset.Scale(2.0 / texSize)
	lproj[12] = offset.X
	lproj[13] = offset.Y

	// re-create view-projection after rounding
	lvp = lproj.Mul(&lview)

	return Cascade{
		Proj:      lproj,
		View:      lview,
		ViewProj:  lvp,
		NearSplit: nearSplit * args.Far,
		FarSplit:  farSplit * args.Far,
	}
}
