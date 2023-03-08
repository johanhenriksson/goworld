package propedit

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func Container(key string, children []node.T) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Layout:     style.Column{},
			AlignItems: style.AlignStart,
			Width:      style.Pct(100),
		},
		Children: children,
	})
}

func Field(key, title string, children []node.T) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Layout:     style.Column{},
			AlignItems: style.AlignStart,
			Width:      style.Pct(100),
			Padding:    style.RectY(4),
		},
		Children: []node.T{
			label.New("label", label.Props{
				Text: title,
				Style: label.Style{
					Color: color.White,
				},
			}),
			rect.New("editor", rect.Props{
				Style: rect.Style{
					Layout:     style.Row{},
					AlignItems: style.AlignStart,
					Width:      style.Pct(100),
					Padding:    style.RectY(2),
				},
				Children: children,
			}),
		},
	})
}
