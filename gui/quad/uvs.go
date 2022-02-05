package quad

import (
	"github.com/johanhenriksson/goworld/math/vec2"
)

type UV struct {
	// Top left UV
	A vec2.T
	// Bottom right UV
	B vec2.T
}

var DefaultUVs = UV{
	A: vec2.New(0, 0),
	B: vec2.New(1, 1),
}

func (uv UV) Inverted() UV {
	return UV{
		A: vec2.New(uv.A.X, 1-uv.A.Y),
		B: vec2.New(uv.B.X, 1-uv.B.Y),
	}
}
