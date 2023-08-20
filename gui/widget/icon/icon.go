package icon

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/render/color"
)

type Icon rune

type IconProps struct {
	Icon        Icon
	Size        int
	Color       color.T
	OnMouseUp   func(mouse.Event)
	OnMouseDown func(mouse.Event)
}

func New(key string, props IconProps) node.T {
	return label.New(key, label.Props{
		Text:        string(props.Icon),
		OnMouseUp:   props.OnMouseUp,
		OnMouseDown: props.OnMouseDown,
		Style: label.Style{
			Color: props.Color,
			Hidden: props.Icon == IconNone,
			Font: style.Font{
				Name: "fonts/MaterialIcons-Regular.ttf",
				Size: props.Size,
			},
		},
	})
}
