package icon

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
)

type Icon rune

type Props struct {
	Icon        Icon
	Style       Style
	OnMouseUp   func(mouse.Event)
	OnMouseDown func(mouse.Event)
}

func New(key string, props Props) node.T {
	return label.New(key, label.Props{
		Text:        string(props.Icon),
		OnMouseUp:   props.OnMouseUp,
		OnMouseDown: props.OnMouseDown,
		Style: label.Style{
			Color:  props.Style.Color,
			Hidden: props.Icon == IconNone,
			Font: style.Font{
				Name: "fonts/MaterialIcons-Regular.ttf",
				Size: props.Style.Size,
			},
			Hover: label.Hover{
				Color: props.Style.Hover.Color,
			},
		},
	})
}
