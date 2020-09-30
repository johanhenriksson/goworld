package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Box struct {
	X      float32
	Y      float32
	Z      float32
	Center vec3.T
}

func (box Box) Min() vec3.T {
	return vec3.T{
		box.Center.X - box.X,
		box.Center.Y - box.Y,
		box.Center.Z - box.Z,
	}
}

func (box Box) Max() vec3.T {
	return vec3.T{
		box.Center.X + box.X,
		box.Center.Y + box.Y,
		box.Center.Z + box.Z,
	}
}
