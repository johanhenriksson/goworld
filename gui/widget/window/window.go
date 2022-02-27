package window

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Props struct {
	Title    string
	Children []node.T
	OnClose  func()
}

func New(key string, props Props) node.T {
	return node.Component(key, props, nil, render)
}

func render(props Props) node.T {
	position, setPosition := hooks.UseState(vec2.New(150, 30))
	dragOffset, setDragOffset := hooks.UseState(vec2.Zero)

	return rect.New("window", rect.Props{
		OnMouseDown: func(e mouse.Event) {},
		Style: rect.Style{
			Position: Absolute{
				Left: Px(position.X),
				Top:  Px(position.Y),
			},
		},
		Children: []node.T{
			rect.New("titlebar", rect.Props{
				Style: TitlebarStyle,
				OnMouseDown: func(e mouse.Event) {
					offset := e.Position().Sub(position)
					setDragOffset(offset)
				},
				OnMouseDrag: func(e mouse.Event) {
					setPosition(e.Position().Sub(dragOffset))
				},
				Children: []node.T{
					label.New("title", label.Props{
						Text:  props.Title,
						Style: TitleStyle,
					}),
					button.New("close", button.Props{
						Style: CloseButtonStyle,
						Text:  "X",
						OnClick: func(e mouse.Event) {
							if props.OnClose != nil {
								props.OnClose()
							}
						},
					}),
				},
			}),
			rect.New("frame", rect.Props{
				Style:    FrameStyle,
				Children: props.Children,
			}),
		},
	})
}

var TitleStyle = label.Style{
	Grow: Grow(1),

	Color: color.White,
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 16,
	},
}

var TitlebarStyle = rect.Style{
	Color:      color.RGBA(0, 0, 0, 0.8),
	Padding:    Px(4),
	Layout:     Row{},
	AlignItems: AlignCenter,
	Pressed: rect.Pressed{
		Color: color.RGBA(0.2, 0.2, 0.2, 0.8),
	},
}

var FrameStyle = rect.Style{
	Color:   color.RGBA(0.1, 0.1, 0.1, 0.8),
	Padding: Px(4),
}

var CloseButtonStyle = button.Style{
	Background: color.RGB(0.597, 0.098, 0.117),
}
