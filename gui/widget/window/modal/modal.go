package modal

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/window"
)

type Props struct {
	Children []node.T
	Title    string
	OnClose  func()
}

func New(key string, props Props) node.T {
	return rect.New(key, rect.Props{
		Children: []node.T{
			window.New(key, window.Props{
				Title:    props.Title,
				OnClose:  props.OnClose,
				Floating: false,
				Children: props.Children,
			}),
		},
		Style: rect.Style{
			Position: style.Absolute{},
			Width:    style.Pct(100),
			Height:   style.Pct(100),

			Layout:         style.Row{},
			JustifyContent: style.JustifyCenter,
			AlignContent:   style.AlignCenter,
			AlignItems:     style.AlignCenter,
		},
	})
}
