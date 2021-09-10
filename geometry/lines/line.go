package lines

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type Line struct {
	Start vec3.T
	End   vec3.T
	Color color.T
}

// L creates a new line segment
func L(start, end vec3.T, color color.T) Line {
	return Line{start, end, color}
}
