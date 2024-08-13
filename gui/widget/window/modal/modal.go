package modal

import (
	"fmt"

	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/window"

	"github.com/samber/lo"
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
				Children: lo.Map(props.Children, modalRow),
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

func modalRow(child node.T, index int) node.T {
	return rect.New(fmt.Sprintf("row-%d", index), rect.Props{
		Style: rect.Style{
			Layout:  style.Row{},
			Width:   style.Pct(100),
			Basis:   style.Pct(100),
			Padding: style.RectY(2),
		},
		Children: []node.T{child},
	})
}
