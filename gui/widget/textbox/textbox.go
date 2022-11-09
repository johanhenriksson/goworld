package textbox

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

type Props struct {
	Style    Style
	Text     string
	OnChange func(string)
}

type Style struct {
	Text label.Style
	Bg   rect.Style
}

func New(key string, props Props) node.T {
	return node.Component(key, props, render)
}

func render(props Props) node.T {
	return rect.New("background", rect.Props{
		Style: props.Style.Bg,
		OnMouseUp: func(e mouse.Event) {
			// grab input focus
		},
		OnMouseDown: func(e mouse.Event) {
		},
		Children: []node.T{
			label.New("input", label.Props{
				Text:     props.Text,
				Style:    props.Style.Text,
				OnChange: props.OnChange,
			}),
		},
	})
}
