package camera

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Frustum struct {
	NTL vec3.T
	NTR vec3.T
	NBL vec3.T
	NBR vec3.T
	FTL vec3.T
	FTR vec3.T
	FBL vec3.T
	FBR vec3.T
}

func (f *Frustum) Bounds() (vec3.T, vec3.T) {
	min := vec3.New(math.InfPos, math.InfPos, math.InfPos)
	max := vec3.New(math.InfNeg, math.InfNeg, math.InfNeg)

	// this does not seem very efficient
	min = vec3.Min(min, f.NTL)
	max = vec3.Max(max, f.NTL)
	min = vec3.Min(min, f.NTR)
	max = vec3.Max(max, f.NTR)
	min = vec3.Min(min, f.NBL)
	max = vec3.Max(max, f.NBL)
	min = vec3.Min(min, f.NBR)
	max = vec3.Max(max, f.NBR)
	min = vec3.Min(min, f.FTL)
	max = vec3.Max(max, f.FTL)
	min = vec3.Min(min, f.FTR)
	max = vec3.Max(max, f.FTR)
	min = vec3.Min(min, f.FBL)
	max = vec3.Max(max, f.FBL)
	min = vec3.Min(min, f.FBR)
	max = vec3.Max(max, f.FBR)

	return min, max
}
