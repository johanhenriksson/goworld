package button

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

type Props struct {
	Text    string
	Style   Style
	OnClick mouse.Callback
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
				Text:  props.Text,
				Style: props.Style.Text,
			}),
		},
	})
}
