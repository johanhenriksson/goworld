package modal

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
	"github.com/johanhenriksson/goworld/render/color"
)

type InputProps struct {
	Title    string
	Message  string
	OnClose  func()
	OnAccept func(string)
}

func NewInput(key string, props InputProps) node.T {
	return node.Component(key, props, renderInput)
}

func renderInput(props InputProps) node.T {
	text, setText := hooks.UseState("")
	return New("inputbox", Props{
		Title:   props.Title,
		OnClose: props.OnClose,
		Children: []node.T{
			label.New("message", label.Props{
				Text: props.Message,
			}),
			textbox.New("input", textbox.Props{
				Text:     text,
				OnChange: setText,
				Style:    textbox.DefaultStyle,
			}),
			button.New("ok", button.Props{
				Text: "OK",
				OnClick: func(e mouse.Event) {
					props.OnAccept(text)
					props.OnClose()
				},
				Style: button.Style{
					Text: label.Style{
						Color: color.Black,
					},
					Bg: rect.Style{
						Color:      style.RGBA(0.5, 0.5, 0.5, 1),
						Padding:    style.RectXY(20, 4),
						AlignItems: style.AlignCenter,
					},
				},
			}),
		},
	})
}
