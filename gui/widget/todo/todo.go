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
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/util"
)

var itemRowStyle = rect.Style{
	Padding: style.Px(5),
	Layout:  style.Row{},
}

var itemLabelStyle = label.Style{
	Grow: style.Grow(1),
}

var removeButtonStyle = button.Style{
	Bg: rect.Style{
		Padding: style.Px(4),
		Color:   color.Red,
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
	Bg: rect.Style{
		Padding: style.Px(4),
		Color:   color.Green,
	},
}

type Props struct{}

func New(key string, props Props) node.T {
	return node.Component(key, props, nil, func(props Props) node.T {
		items, setItems := hooks.UseState([]string{
			"apples",
			"oranges",
			"mangoes",
		})
		itemTitle, setItemTitle := hooks.UseState("")

		return window.New(key, window.Props{
			Title: "todo",
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
							Text:  "+",
							Style: addButtonStyle,
							OnClick: func(e mouse.Event) {
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
