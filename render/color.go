package render

import (
	"fmt"
	"image/color"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
)

var White = Color{1, 1, 1, 1}
var Black = Color{0, 0, 0, 1}
var Red = Color{1, 0, 0, 1}
var Green = Color{0, 1, 0, 1}
var Blue = Color{0, 0, 1, 1}
var Transparent = Color{0, 0, 0, 0}

/** Color type */
type Color struct {
	R, G, B, A float32
}

func Color4(r, g, b, a float32) Color {
	return Color{r, g, b, a}
}

func (c Color) RGBA() color.RGBA {
	return color.RGBA{
		uint8(255.0 * c.R),
		uint8(255.0 * c.G),
		uint8(255.0 * c.B),
		uint8(255.0 * c.A),
	}
}

func (c Color) Vec3() vec3.T {
	return vec3.New(c.R, c.G, c.B)
}

func (c Color) Vec4() vec4.T {
	return vec4.New(c.R, c.G, c.B, c.A)
}

func (c Color) String() string {
	return fmt.Sprintf("(R:%.2f G:%.2f B:%.2f A:%.2f)", c.R, c.G, c.B, c.A)
}
