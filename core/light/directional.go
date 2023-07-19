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

type Directional struct {
	object.Component
	cascades []Cascade

	Color     *object.Property[color.T]
	Intensity *object.Property[float32]
	Shadows   *object.Property[bool]
}

var _ T = &Directional{}

func NewDirectional(args DirectionalArgs) *Directional {
	return object.NewComponent(&Directional{
		cascades: make([]Cascade, args.Cascades),

		Color:     object.NewProperty(args.Color),
		Intensity: object.NewProperty(args.Intensity),
		Shadows:   object.NewProperty(args.Shadows),
	})
}

func (lit *Directional) Name() string        { return "DirectionalLight" }
func (lit *Directional) Type() Type          { return TypeDirectional }
func (lit *Directional) CastShadows() bool   { return lit.Shadows.Get() }
func (lit *Directional) Cascades() []Cascade { return lit.cascades }

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

func (lit *Directional) PreDraw(args render.Args, scene object.Object) error {
	lit.updateCascades(args)
	return nil
}

func (lit *Directional) LightDescriptor(args render.Args, cascade int) Descriptor {
	ldir := lit.Transform().Forward()
	return Descriptor{
		Type:       TypeDirectional,
		Position:   vec4.Extend(ldir, 0),
		Color:      lit.Color.Get(),
		Intensity:  lit.Intensity.Get(),
		View:       lit.cascades[cascade].View,
		Projection: lit.cascades[cascade].Proj,
		ViewProj:   lit.cascades[cascade].ViewProj,
	}
}

func (lit *Directional) updateCascades(args render.Args) {
	for i, _ := range lit.cascades {
		lit.cascades[i] = lit.calculateCascade(args, i, len(lit.cascades))
	}
}

func (lit *Directional) calculateCascade(args render.Args, cascade, cascades int) Cascade {
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
