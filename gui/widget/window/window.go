package window

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
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

func New(key string, props Props) node.T {
	return node.Component(key, props, render)
}

func render(props Props) node.T {
	position, setPosition := hooks.UseState(props.Position)
	dragOffset, setDragOffset := hooks.UseState(vec2.Zero)
	var cssPos style.PositionProp = style.Relative{}
	if props.Floating {
		cssPos = style.Absolute{
			Left: style.Px(position.X),
			Top:  style.Px(position.Y),
		}
	}

	return rect.New("window", rect.Props{
		OnMouseDown: func(e mouse.Event) {
			e.Consume()
		},
		Style: rect.Style{
			Position: cssPos,
			MaxWidth: props.Style.MaxWidth,
			MinWidth: props.Style.MinWidth,
		},
		Children: []node.T{
			rect.New("titlebar", rect.Props{
				Style:     TitlebarStyle,
				OnMouseUp: func(e mouse.Event) { e.Consume() },
				OnMouseDown: func(e mouse.Event) {
					if !props.Floating {
						return
					}
					offset := e.Position().Sub(position)
					setDragOffset(offset)
					e.Consume()
				},
				OnMouseDrag: func(e mouse.Event) {
					if !props.Floating {
						return
					}
					setPosition(e.Position().Sub(dragOffset))
					e.Consume()
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
