package label

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/render/color"
)

var DefaultStyle = &Style{}

type Style struct {
	Extends *Style
	Hover   Hover

	Hidden     bool
	Font       FontProp
	Color      ColorProp
	LineHeight LineHeightProp
}

type Hover struct {
	Color ColorProp
}

func (style *Style) Apply(w T, state State) {
	if style.Extends == nil {
		if style != DefaultStyle {
			DefaultStyle.Apply(w, state)
		}
	} else {
		style.Extends.Apply(w, state)
	}

	if style.Font != nil {
		style.Font.ApplyFont(w)
	}
	if style.Color != nil {
		rgba := style.Color.Vec4()
		w.SetFontColor(color.RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}

	if state.Hovered {
		style.Hover.Apply(w)
	}
}

func (style *Hover) Apply(w T) {
	if style.Color != nil {
		rgba := style.Color.Vec4()
		w.SetFontColor(color.RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}
}
