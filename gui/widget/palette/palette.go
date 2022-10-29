package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/window"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/util"
)

type Props struct {
	Palette color.Palette
	OnPick  func(color.T)
}

func New(key string, props Props) node.T {
	return node.Component(key, props, nil, render)
}

func render(props Props) node.T {
	perRow := 5

	selected, setSelected := hooks.UseState(props.Palette[3])

	colors := util.MapIdx(props.Palette, func(c color.T, i int) node.T {
		return rect.New(fmt.Sprintf("color%d", i), rect.Props{
			Style: SwatchStyle.Extend(rect.Style{
				Color: c,
			}),
			OnMouseUp: func(e mouse.Event) {
				setSelected(c)
				if props.OnPick != nil {
					props.OnPick(c)
				}
			},
		})
	})

	rows := util.MapIdx(util.Chunks(colors, perRow), func(colors []node.T, i int) node.T {
		return rect.New(fmt.Sprintf("row%d", i), rect.Props{
			Style: rect.Style{
				Width:  Pct(100),
				Layout: Row{},
			},
			Children: colors,
		})
	})

	return window.New("palette", window.Props{
		Title: "Palette",
		Children: []node.T{
			rect.New("selected", rect.Props{
				Style: rect.Style{
					Layout:     Row{},
					MaxWidth:   Pct(100),
					AlignItems: AlignCenter,
				},
				Children: []node.T{
					label.New("selected", label.Props{
						Text: "Selected",
						Style: label.Style{
							Color: color.White,
							Grow:  Grow(1),
						},
					}),
					rect.New("preview", rect.Props{
						Style: rect.Style{
							Color:  selected,
							Height: Px(20),
							Basis:  Pct(16),
						},
					}),
				},
			}),
			rect.New("grid", rect.Props{
				Children: rows,
			}),
		},
	})
}
