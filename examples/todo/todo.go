package main

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/app"
	"github.com/johanhenriksson/goworld/gui"
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
	"github.com/samber/lo"
)

func main() {
	app.Run(
		app.Args{
			Width:  1200,
			Height: 800,
			Title:  "goworld: todo",
		},
		func(scene object.Object) {
			object.Attach(scene, gui.New(func() node.T {
				return rect.New("gui", rect.Props{
					Style: rect.Style{
						Position: style.Absolute{},
					},
					Children: []node.T{
						NewTodo("todo", Props{}),
					},
				})
			}))
		},
	)
}

type Props struct{}

func NewTodo(key string, props Props) node.T {
	return node.Component(key, props, func(props Props) node.T {
		items, setItems := hooks.UseState([]string{"testy", "mctestface"})
		itemTitle, setItemTitle := hooks.UseState("")
		addStyle := addButtonStyle
		if len(itemTitle) == 0 {
			addStyle.BgColor = color.DarkGrey
		}

		return window.New(key, window.Props{
			Title:    "todo app",
			Position: vec2.New(500, 300),
			Floating: true,
			Style: window.Style{
				MinWidth: style.Px(200),
			},
			Children: []node.T{
				rect.New("list", rect.Props{
					Children: lo.Map(items, func(text string, idx int) node.T {
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
					Style: itemRowStyle,
					Children: []node.T{
						textbox.New("title", textbox.Props{
							Style:    itemInputStyle,
							OnChange: setItemTitle,
							Text:     itemTitle,
						}),
						button.New("add", button.Props{
							Text:  "+",
							Style: addStyle,
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

var itemRowStyle = rect.Style{
	Layout:     style.Row{},
	Width:      style.Pct(100),
	Grow:       style.Grow(1),
	AlignItems: style.AlignCenter,
	Padding:    style.Px(4),
}

var itemLabelStyle = label.Style{
	Grow: style.Grow(1),
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
	Padding:   style.Px(4),
	BgColor:   color.RGB(0.2, 0.6, 0.4),
	TextColor: color.White,
}

var removeButtonStyle = button.Style{
	Padding:   style.Px(4),
	BgColor:   window.CloseButtonStyle.BgColor,
	TextColor: color.White,
}
