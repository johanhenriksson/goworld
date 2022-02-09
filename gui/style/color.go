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
