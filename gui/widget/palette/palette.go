package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
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

	selected, setSelected := hooks.UseState(props.Palette[0])

	colors := util.Map(props.Palette, func(i int, c color.T) node.T {
		return rect.New(fmt.Sprintf("color%d", i), rect.Props{
			Style: SwatchStyle.Extend(rect.Style{
				Color: c,
			}),
			OnClick: func(e mouse.Event) {
				setSelected(c)
				if props.OnPick != nil {
					props.OnPick(c)
				}
			},
		})
	})

	rows := util.Map(util.Chunks(colors, perRow), func(i int, colors []node.T) node.T {
		return rect.New(fmt.Sprintf("row%d", i), rect.Props{
			Style: rect.Style{
				Width:  Pct(100),
				Layout: Row{},
			},
			Children: colors,
		})
	})

	return rect.New("window", rect.Props{
		OnClick: func(e mouse.Event) {},
		Style: rect.Style{
			Color:   color.Black.WithAlpha(0.9),
			Padding: Px(4),
			Layout:  Column{},
			Position: Absolute{
				Top:  Px(0),
				Left: Pct(100),
			},
		},
		Children: []node.T{
			label.New("title", label.Props{
				Text:  "Palette",
				Style: TitleStyle,
			}),
			rect.New("selected", rect.Props{
				Style: rect.Style{
					Layout:   Row{},
					MaxWidth: Pct(100),
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
						Style: SwatchStyle.Extend(rect.Style{
							Color: selected,
						}),
					}),
				},
			}),
			rect.New("grid", rect.Props{
				Children: rows,
			}),
		},
	})
}
