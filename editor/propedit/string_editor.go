package propedit

import (
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
)

type StringProps struct {
	Label    string
	Value    string
	OnChange func(string)
	Validate func(string) bool
}

func String(key string, props StringProps) node.T {
	return node.Component(key, props, func(props StringProps) node.T {
		value, setValue := hooks.UseState(props.Value)
		valid, setValid := hooks.UseState(true)
		text, setText := hooks.UseState(props.Value)

		validate := func(string) bool { return true }
		if props.Validate != nil {
			validate = props.Validate
		}

		onChange := func(newText string) {
			setValid(validate(newText))
			setText(newText)
		}

		revert := func() {
			setText(value)
			setValid(true)
		}

		updateValue := func() {
			if validate(text) {
				// if its a valid float, run the onchange callback
				if props.OnChange != nil && text != props.Value {
					props.OnChange(text)
				}
			} else {
				revert()
			}
		}

		// how to properly observe changing values with hooks?
		// this hook would make sure the label is updated any time the prop changes,
		// even due to an external change (or OnChange effects)
		hooks.UseEffect(func() {
			setValue(props.Value)
			setValid(true)
			setText(props.Value)
		}, props.Value)

		// set style depending on validity of input
		textboxStyle := textbox.DefaultStyle
		if !valid {
			textboxStyle = textbox.InputErrorStyle
		}

		return rect.New(key, rect.Props{
			Style: rect.Style{
				Layout:     style.Row{},
				AlignItems: style.AlignCenter,
				Basis:      style.Pct(100),
				Grow:       style.Grow(0),
				Shrink:     style.Shrink(1),
			},
			Children: []node.T{
				label.New("label", label.Props{
					Text: props.Label,
					Style: label.Style{
						Grow:   style.Grow(0),
						Shrink: style.Shrink(0),
					},
				}),
				textbox.New("value", textbox.Props{
					Text:     text,
					Style:    textboxStyle,
					OnChange: onChange,
					OnBlur:   updateValue,
					OnAccept: updateValue,
					OnReject: revert,
				}),
			},
		})
	})
}
