package light

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
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

	Color     object.Property[color.T]
	Intensity object.Property[float32]
	Shadows   object.Property[bool]

	CascadeLambda object.Property[float32]
	CascadeBlend  object.Property[float32]
}

var _ T = &Directional{}

func init() {
	object.Register[*Directional](object.TypeInfo{
		Name:        "Directional Light",
		Deserialize: DeserializeDirectional,
		Create: func() (object.Component, error) {
			return NewDirectional(DirectionalArgs{
				Color:     color.White,
				Intensity: 1,
				Shadows:   true,
				Cascades:  4,
			}), nil
		},
	})
}

func NewDirectional(args DirectionalArgs) *Directional {
	lit := object.NewComponent(&Directional{
		cascades: make([]Cascade, args.Cascades),

		Color:     object.NewProperty(args.Color),
		Intensity: object.NewProperty(args.Intensity),
		Shadows:   object.NewProperty(args.Shadows),

		CascadeLambda: object.NewProperty[float32](0.9),
		CascadeBlend:  object.NewProperty[float32](3.0),
	})
	return lit
}

func (lit *Directional) Name() string      { return "DirectionalLight" }
func (lit *Directional) Type() Type        { return TypeDirectional }
func (lit *Directional) CastShadows() bool { return lit.Shadows.Get() }

func farSplitDist(cascade, cascades int, near, far, splitLambda float32) float32 {
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
	d := splitLambda*(log-uniform) + uniform
	return (d - near) / clipRange
}

func nearSplitDist(cascade, cascades int, near, far, splitLambda float32) float32 {
	if cascade == 0 {
		return 0
	}
	return farSplitDist(cascade-1, cascades, near, far, splitLambda)
}

func (lit *Directional) PreDraw(args draw.Args, scene object.Object) error {
	// update cascades
	for i, _ := range lit.cascades {
		lit.cascades[i] = lit.calculateCascade(args, i, len(lit.cascades))
	}
	return nil
}

func (lit *Directional) calculateCascade(args draw.Args, cascade, cascades int) Cascade {
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
		frustumCorners[i] = args.Camera.ViewProjInv.TransformPoint(corner)
	}

	// squash
	nearSplit := nearSplitDist(cascade, cascades, args.Camera.Near, args.Camera.Far, lit.CascadeLambda.Get())
	farSplit := farSplitDist(cascade, cascades, args.Camera.Near, args.Camera.Far, lit.CascadeLambda.Get())
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
		NearSplit: nearSplit * args.Camera.Far,
		FarSplit:  farSplit * args.Camera.Far,
	}
}

func (lit *Directional) LightData(shadowmaps ShadowmapStore) uniform.Light {
	ldir := lit.Transform().Forward()
	entry := uniform.Light{
		Type:      uint32(TypeDirectional),
		Position:  vec4.Extend(ldir, 0),
		Color:     lit.Color.Get(),
		Intensity: lit.Intensity.Get(),
		Range:     lit.CascadeBlend.Get(),
	}

	for cascadeIndex, cascade := range lit.cascades {
		entry.ViewProj[cascadeIndex] = cascade.ViewProj
		entry.Distance[cascadeIndex] = cascade.FarSplit
		if handle, exists := shadowmaps.Lookup(lit, cascadeIndex); exists {
			entry.Shadowmap[cascadeIndex] = uint32(handle)
		}
	}

	return entry
}

func (lit *Directional) Shadowmaps() int {
	return len(lit.cascades)
}

func (lit *Directional) ShadowProjection(mapIndex int) uniform.Camera {
	cascade := lit.cascades[mapIndex]
	return uniform.Camera{
		Proj:        cascade.Proj,
		View:        cascade.View,
		ViewProj:    cascade.ViewProj,
		ProjInv:     cascade.Proj.Invert(),
		ViewInv:     cascade.View.Invert(),
		ViewProjInv: cascade.ViewProj.Invert(),
		Eye:         vec4.Extend(lit.Transform().Position(), 0),
		Forward:     vec4.Extend(lit.Transform().Forward(), 0),
	}
}

type DirectionalState struct {
	object.ComponentState
	DirectionalArgs
}

func (lit *Directional) Serialize(enc object.Encoder) error {
	return enc.Encode(DirectionalState{
		// send help
		ComponentState: object.NewComponentState(lit.Component),
		DirectionalArgs: DirectionalArgs{
			Color:     lit.Color.Get(),
			Intensity: lit.Intensity.Get(),
			Shadows:   lit.Shadows.Get(),
		},
	})
}

func DeserializeDirectional(dec object.Decoder) (object.Component, error) {
	var state DirectionalState
	if err := dec.Decode(&state); err != nil {
		return nil, err
	}

	obj := NewDirectional(state.DirectionalArgs)
	obj.Component = state.ComponentState.New()
	return obj, nil
}
