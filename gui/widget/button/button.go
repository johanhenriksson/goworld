package button

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

type Props struct {
	Text    string
	Style   Style
	OnClick mouse.Callback
}

func New(key string, props Props) node.T {
	return node.Component(key, props, render)
}

func render(props Props) node.T {
	return rect.New("background", rect.Props{
		Style: rect.Style{
			Layout:     style.Row{},
			AlignItems: style.AlignCenter,

			Color:   props.Style.BgColor,
			Padding: props.Style.Padding,
			Margin:  props.Style.Margin,
			Radius:  props.Style.Radius,
			Border:  props.Style.Border,

			Hover: rect.Hover{
				Color: props.Style.Hover.BgColor,
			},
		},
		OnMouseUp: func(e mouse.Event) {
			if props.OnClick != nil {
				props.OnClick(e)
			}
			e.Consume()
		},
		OnMouseDown: func(e mouse.Event) {
			e.Consume()
		},
		Children: []node.T{
			label.New("label", label.Props{
				Text: props.Text,
				Style: label.Style{
					Color: props.Style.TextColor,
					Hover: label.Hover{
						Color: props.Style.Hover.TextColor,
					},
				},
			}),
		},
	})
}
