package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
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
	position, setPosition := hooks.UseState(vec2.New(150, 30))
	dragOffset, setDragOffset := hooks.UseState(vec2.Zero)

	colors := util.Map(props.Palette, func(i int, c color.T) node.T {
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
		OnMouseUp: func(e mouse.Event) {},
		Style: rect.Style{
			Color:   color.Black.WithAlpha(0.9),
			Padding: Px(4),
			Layout:  Column{},
			Position: Absolute{
				Left: Px(position.X),
				Top:  Px(position.Y),
			},
		},
		Children: []node.T{
			rect.New("titlebar", rect.Props{
				OnMouseDown: func(e mouse.Event) {
					offset := e.Position().Sub(position)
					setDragOffset(offset)
				},
				OnMouseDrag: func(e mouse.Event) {
					setPosition(e.Position().Sub(dragOffset))
				},
				Children: []node.T{
					label.New("title", label.Props{
						Text:  "Palette",
						Style: TitleStyle,
					}),
				},
				Style: rect.Style{
					Color: color.RGBA(0, 0, 0, 0.8),
					Pressed: rect.Pressed{
						Color: color.RGBA(0.5, 0.5, 0.5, 0.8),
					},
				},
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
