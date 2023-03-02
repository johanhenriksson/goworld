package style

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

type Border struct {
	Width Px
	Color ColorProp
}

func (b Border) ApplyBorder(w BorderWidget) {
	w.Flex().StyleSetBorder(flex.EdgeAll, float32(b.Width))
	c := b.Color.Vec4()
	w.SetBorderColor(color.RGBA(c.X, c.Y, c.Z, c.W))
}

type BorderProp interface {
	ApplyBorder(BorderWidget)
}

type BorderWidget interface {
	FlexWidget
	SetBorderColor(color.T)
}
