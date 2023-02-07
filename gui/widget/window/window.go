package window

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Props struct {
	Title    string
	Style    Style
	Position vec2.T
	Children []node.T
	OnClose  func()
	Floating bool
}

type Style struct {
	MinWidth MinWidthProp
	MaxWidth MaxWidthProp
}

func New(key string, props Props) node.T {
	return node.Component(key, props, render)
}

func render(props Props) node.T {
	position, setPosition := hooks.UseState(props.Position)
	dragOffset, setDragOffset := hooks.UseState(vec2.Zero)
	var cssPos style.PositionProp = style.Relative{}
	if props.Floating {
		cssPos = Absolute{
			Left: Px(position.X),
			Top:  Px(position.Y),
		}
	}

	return rect.New("window", rect.Props{
		OnMouseDown: func(e mouse.Event) {},
		Style: rect.Style{
			Position: cssPos,
			MaxWidth: props.Style.MaxWidth,
			MinWidth: props.Style.MinWidth,
		},
		Children: []node.T{
			rect.New("titlebar", rect.Props{
				Style: TitlebarStyle,
				OnMouseDown: func(e mouse.Event) {
					if !props.Floating {
						return
					}
					offset := e.Position().Sub(position)
					setDragOffset(offset)
				},
				OnMouseDrag: func(e mouse.Event) {
					if !props.Floating {
						return
					}
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

	Color: RGB(1, 1, 1),
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 16,
	},
}

var TitlebarStyle = rect.Style{
	Color:      RGBA(0, 0, 0, 0.8),
	Padding:    Px(4),
	Layout:     Row{},
	AlignItems: AlignCenter,
	Pressed: rect.Pressed{
		Color: RGBA(0.2, 0.2, 0.2, 0.8),
	},
}

var FrameStyle = rect.Style{
	Color:      RGBA(0.1, 0.1, 0.1, 0.8),
	Padding:    RectXY(10, 10),
	Layout:     Column{},
	AlignItems: AlignStart,
}

var CloseButtonStyle = button.Style{
	Bg: rect.Style{
		Color: RGB(0.597, 0.098, 0.117),
		Padding: Rect{
			Left:   5,
			Right:  5,
			Top:    2,
			Bottom: 2,
		},

		Hover: rect.Hover{
			Color: RGB(0.3, 0.3, 0.3),
		},
	},
}
