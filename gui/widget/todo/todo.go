package todo

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
	"github.com/johanhenriksson/goworld/gui/widget/window"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/util"
)

var itemRowStyle = rect.Style{
	Layout:     style.Row{},
	AlignItems: style.AlignCenter,
}

var itemLabelStyle = label.Style{
	Grow: style.Grow(1),
}

var removeButtonStyle = button.Style{
	BgColor: color.Red,
	Padding: style.Px(4),
	Margin: style.Rect{
		Left: 40,
	},
}

var itemInputStyle = textbox.Style{
	Text: label.Style{
		Color: color.Black,
	},
	Bg: rect.Style{
		Padding: style.Px(4),
		Grow:    style.Grow(1),
		Color:   color.White,
	},
}

var addButtonStyle = button.Style{
	Padding: style.Px(4),
	BgColor: color.Green,
	Margin: style.Rect{
		Left: 40,
	},
}

type Props struct{}

func New(key string, props Props) node.T {
	return node.Component(key, props, func(props Props) node.T {
		items, setItems := hooks.UseState([]string{})
		itemTitle, setItemTitle := hooks.UseState("")
		addColor := color.Green
		if len(itemTitle) == 0 {
			addColor = color.DarkGrey
		}

		return window.New(key, window.Props{
			Title:    "todo app",
			Position: vec2.New(250, 400),
			Style: window.Style{
				MinWidth: style.Px(200),
			},
			Children: []node.T{
				rect.New("list", rect.Props{
					Children: util.MapIdx(items, func(text string, idx int) node.T {
						return rect.New(fmt.Sprintf("item:%d", idx), rect.Props{
							Style: itemRowStyle,
							Children: []node.T{
								label.New("title", label.Props{
									Text:  text,
									Style: itemLabelStyle,
								}),
								button.New("remove", button.Props{
									Text:  "X",
									Style: removeButtonStyle,
									OnClick: func(e mouse.Event) {
										newItems := append(items[:idx], items[idx+1:]...)
										setItems(newItems)
									},
								}),
							},
						})
					}),
				}),
				rect.New("entry", rect.Props{
					Style: rect.Style{
						Layout: style.Row{},
					},
					Children: []node.T{
						textbox.New("title", textbox.Props{
							Style:    itemInputStyle,
							OnChange: setItemTitle,
							Text:     itemTitle,
						}),
						button.New("add", button.Props{
							Text: "+",
							Style: button.Style{
								Padding: style.Px(4),
								BgColor: addColor,
							},
							OnClick: func(e mouse.Event) {
								if len(itemTitle) == 0 {
									return
								}
								newItems := append(items, itemTitle)
								setItems(newItems)
								setItemTitle("")
							},
						}),
					},
				}),
			},
		})
	})
}
