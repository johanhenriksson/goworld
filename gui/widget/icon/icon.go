package icon

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/render/color"
)

type Icon string

const (
	IconExpandMore   = Icon(rune(0xe5cf))
	IconChevronRight = Icon(rune(0xe5cc))
)

type IconProps struct {
	Icon    Icon
	Size    int
	Color   color.T
	OnClick func(mouse.Event)
}

func New(key string, props IconProps) node.T {
	return label.New(key, label.Props{
		Text:    string(props.Icon),
		OnClick: props.OnClick,
		Style: label.Style{
			Color: props.Color,
			Font: style.Font{
				Name: "fonts/MaterialIcons-Regular.ttf",
				Size: props.Size,
			},
		},
	})
}
