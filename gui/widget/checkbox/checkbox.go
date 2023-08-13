package checkbox

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

type Props struct {
	Style    Style
	Checked  bool
	OnChange func(bool)
}

type Style struct {
}

func New(key string, props Props) node.T {
	return node.Component(key, props, render)
}

func render(props Props) node.T {
	checker := icon.IconCheckboxBlank
	if props.Checked {
		checker = icon.IconCheckboxChecked
	}

	click := func(e mouse.Event) {
		props.OnChange(!props.Checked)
	}

	return rect.New("background", rect.Props{
		OnMouseUp: func(e mouse.Event) {
			// grab input focus
			// ... but how?
			e.Consume()
		},
		OnMouseDown: click,

		Style: rect.Style{
			Margin: style.Rect{
				Right: 4,
			},
		},

		Children: []node.T{
			icon.New("checked", icon.IconProps{
				Size:  18,
				Icon:  checker,
				Color: color.White,
			}),
		},
	})
}
