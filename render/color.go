package render

import (
	"image/color"

	mgl "github.com/go-gl/mathgl/mgl32"
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

func (c Color) Vec3() mgl.Vec3 {
	return mgl.Vec3{c.R, c.G, c.B}
}

func (c Color) Vec4() mgl.Vec4 {
	return mgl.Vec4{c.R, c.G, c.B, c.A}
}
