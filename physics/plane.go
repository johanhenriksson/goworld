package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Plane struct {
	Normal vec3.T
	D      float32
}
