package camera

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Frustum struct {
	Corners vec3.Array
	Center  vec3.T
	Min     vec3.T
	Max     vec3.T
}

var ndc_corners = vec3.Array{
	vec3.New(-1, 1, 1),  // NTL
	vec3.New(1, 1, 1),   // NTR
	vec3.New(-1, -1, 1), // NBL
	vec3.New(1, -1, 1),  // NBR
	vec3.New(-1, 1, 0),  // FTL
	vec3.New(1, 1, 0),   // FTR
	vec3.New(-1, -1, 0), // FBL
	vec3.New(1, -1, 0),  // FBR
}

// NewFrustum creates a view frustum from an inverse view projection matrix by unprojecting the corners of the NDC cube.
func NewFrustum(vpi mat4.T) Frustum {
	return Frustum{
		Corners: ndc_corners,
		Center:  vec3.Zero,
		Min:     vec3.New(-1, -1, -1),
		Max:     vec3.One,
	}.Transform(vpi)
}

// Transform returns a new frustum with all its vertices transformed by the given matrix
func (f Frustum) Transform(transform mat4.T) Frustum {
	corners := make(vec3.Array, 8)
	center := vec3.Zero
	min := vec3.New(math.InfPos, math.InfPos, math.InfPos)
	max := vec3.New(math.InfNeg, math.InfNeg, math.InfNeg)
	for i, corner := range f.Corners {
		corner = transform.TransformPoint(corner)
		center = center.Add(corner)
		min = vec3.Min(min, corner)
		max = vec3.Max(max, corner)
		corners[i] = corner
	}
	center = center.Scaled(1 / 8.0)
	return Frustum{
		Corners: corners,
		Center:  center,
		Min:     min,
		Max:     max,
	}
}
