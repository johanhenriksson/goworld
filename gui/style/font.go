package style

import (
	"github.com/johanhenriksson/goworld/assets"
)

type Font struct {
	Name string
	Size int
}

func (f Font) ApplyFont(fw FontWidget) {
	font := assets.GetFont(f.Name, f.Size)
	fw.SetFont(font)
	fw.SetFontSize(f.Size)
}
