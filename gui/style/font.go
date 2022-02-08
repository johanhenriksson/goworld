package style

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
)

type Font struct {
	Name string
	Size int
}

type FontWidget interface {
	SetFont(font.T)
	SetFontSize(int)
	SetFontColor(color.T)
	SetLineHeight(float32)
}

func (f Font) ApplyFont(w widget.T) {
	font := assets.GetFont(f.Name, f.Size)
	if fw, ok := w.(FontWidget); ok {
		fw.SetFont(font)
		fw.SetFontSize(f.Size)
	}
}
