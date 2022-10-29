package style

// awkward attempt at making styles compatible with the render packages color values

import (
	icolor "image/color"

	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

type Colorizable interface {
	SetColor(color.T)
}

type ColorProp interface {
	RGBA() icolor.RGBA
	Vec4() vec4.T
}

func RGB(r, g, b float32) color.T {
	// alias
	return color.RGB(r, g, b)
}

func RGBA(r, g, b, a float32) color.T {
	// alias
	return color.RGBA(r, g, b, a)
}
