package button

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

type Props struct {
	Text       string
	Background style.Sheet
	Label      style.Sheet
	// how to do hover styles?
	OnClick mouse.Callback
}

func New(key string, props Props) node.T {
	return node.Component(key, props, nil, render)
}

func render(props Props) node.T {
	onclick := func(e mouse.Event) {
		if props.OnClick != nil {
			props.OnClick(e)
		}
		e.Consume()
	}
	return rect.New("background", rect.Props{
		Style:   props.Background,
		OnClick: onclick,
		Children: []node.T{
			label.New("label", label.Props{
				Text:    props.Text,
				Style:   props.Label,
				OnClick: onclick,
			}),
		},
	})
}
