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
	"github.com/johanhenriksson/goworld/util"
)

type ItemProps struct {
	Key      string
	Title    string
	Style    Style
	OpenDown bool

	Items   []ItemProps
	OnClick mouse.Callback
	Close   func()
}

type Style struct {
	Color      color.T
	HoverColor color.T
	TextColor  color.T
}

var DefaultStyle = Style{
	Color:      color.RGB(0.76, 0.76, 0.76),
	HoverColor: color.RGB(0.85, 0.85, 0.85),
	TextColor:  color.Black,
}

var menuPadding = style.Rect{
	Top:    4,
	Bottom: 2,
	Left:   5,
	Right:  5,
}

type Props struct {
	Style Style
	Items []ItemProps
}

func Menu(key string, props Props) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Color:   props.Style.Color,
			Width:   style.Pct(100),
			Layout:  style.Row{},
			ZOffset: 100,
		},
		Children: []node.T{
			rect.New("main-menu", rect.Props{
				Style: rect.Style{
					Layout: style.Row{},
				},
				Children: util.Map(props.Items, func(item ItemProps) node.T {
					return Item(item.Key, ItemProps{
						Key:      item.Key,
						Title:    item.Title,
						Style:    props.Style,
						Items:    item.Items,
						OnClick:  item.OnClick,
						OpenDown: true,
					})
				}),
			}),
			rect.New("menu-spacer", rect.Props{
				Style: rect.Style{
					Grow:  style.Grow(1),
					Width: style.Pct(100),
				},
			}),
		},
	})
}

func Item(key string, props ItemProps) node.T {
	return node.Component(key, props, func(prop ItemProps) node.T {
		open, setOpen := hooks.UseState(false)

		close := func() {
			setOpen(false)
			if props.Close != nil {
				props.Close()
			}
		}

		var items []node.T
		if open {
			items = util.Map(props.Items, func(item ItemProps) node.T {
				return Item(item.Key, ItemProps{
					Key:      item.Key,
					Title:    item.Title,
					Style:    props.Style,
					Items:    item.Items,
					OnClick:  item.OnClick,
					Close:    close,
					OpenDown: false,
				})
			})
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

		return rect.New(props.Key, rect.Props{
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
								if len(props.Items) > 0 {
									// open submenu
									setOpen(!open)
								} else {
									e.Consume()
									if props.OnClick != nil {
										props.OnClick(e)
									}
									close()
								}
							},
							Style: label.Style{
								Color: props.Style.TextColor,
							},
						}),
					},
				}),
				rect.New(props.Key+"-items", rect.Props{
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
