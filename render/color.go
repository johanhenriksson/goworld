package render

import (
	"image/color"
)

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
