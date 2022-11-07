package menu

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

type MainItemProps struct {
	Title    string
	Style    Style
	OpenDown bool
}

type Style struct {
	Color      color.T
	HoverColor color.T
	TextColor  color.T
}

var menuPadding = style.Rect{
	Top:    4,
	Bottom: 2,
	Left:   5,
	Right:  5,
}

type Props struct {
	Style Style
}

func Menu(key string, props Props) node.T {
	return node.Component(key, props, nil, func(props Props) node.T {
		menuStyle := Style{
			Color:      color.RGB(0.76, 0.76, 0.76),
			HoverColor: color.RGB(0.85, 0.85, 0.85),
			TextColor:  color.Black,
		}
		return rect.New("gui-menu", rect.Props{
			Style: rect.Style{
				Color:   color.RGB(0.76, 0.76, 0.76),
				Width:   style.Pct(100),
				Layout:  style.Row{},
				ZOffset: 100,
			},
			Children: []node.T{
				MainItem("gui-menu-file", MainItemProps{
					Title:    "File",
					Style:    menuStyle,
					OpenDown: true,
				}),
				MainItem("gui-menu-edit", MainItemProps{
					Title:    "Edit",
					Style:    menuStyle,
					OpenDown: true,
				}),
				rect.New("gui-menu-spacer", rect.Props{
					Style: rect.Style{
						Width: style.Pct(100),
						Grow:  style.Grow(1),
					},
				}),
			},
		})
	})
}

func MainItem(key string, props MainItemProps) node.T {
	return node.Component(key, props, nil, func(prop MainItemProps) node.T {
		open, setOpen := hooks.UseState(false)

		items := make([]node.T, 0, 16)
		if open {
			items = append(items, MainItem("exit", MainItemProps{
				Title:    "exit",
				Style:    props.Style,
				OpenDown: false,
			}))
		}

		itemPos := style.Absolute{
			Top:  style.Pct(100),
			Left: style.Pct(0),
		}
		if !props.OpenDown {
			itemPos = style.Absolute{
				Left: style.Pct(100),
				Top:  style.Pct(0),
			}
		}

		return rect.New("gui-menu-"+props.Title, rect.Props{
			Style: rect.Style{
				Shrink: style.Shrink(1),
				Color:  props.Style.Color,
				Hover: rect.Hover{
					Color: props.Style.HoverColor,
				},
			},
			Children: []node.T{
				rect.New("title-box", rect.Props{
					Style: rect.Style{
						Padding: menuPadding,
					},
					Children: []node.T{
						label.New("title", label.Props{
							Text: props.Title,
							OnClick: func(e mouse.Event) {
								setOpen(!open)
								e.Consume()
							},
							Style: label.Style{
								Color: props.Style.TextColor,
							},
						}),
					},
				}),
				rect.New("gui-menu-"+props.Title+"-items", rect.Props{
					OnMouseUp: gui.ConsumeMouse,
					Style: rect.Style{
						Position: itemPos,
						MinWidth: style.Px(200),
						Color:    prop.Style.Color,
					},
					Children: items,
				}),
			},
		})
	})
}
