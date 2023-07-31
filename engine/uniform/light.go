package uniform

import (
	"github.com/johanhenriksson/goworld/math/mat4"
)

type Light struct {
	ViewProj  [4]mat4.T
	Shadowmap [4]uint32
	Distance  [4]float32
}
