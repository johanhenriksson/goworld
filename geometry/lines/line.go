package lines

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Line struct {
	Start vec3.T
	End   vec3.T
	Color render.Color
}

// L creates a new line segment
func L(start, end vec3.T, color render.Color) Line {
	return Line{start, end, color}
}
