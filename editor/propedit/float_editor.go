package propedit

import (
	"strconv"
	"strings"

	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
)

type FloatProps struct {
	Label    string
	Value    float32
	OnChange func(float32)
}

func fmtFloat(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func Float(key string, props FloatProps) node.T {
	return node.Component(key, props, func(props FloatProps) node.T {
		value, setValue := hooks.UseState(props.Value)
		valid, setValid := hooks.UseState(true)
		text, setText := hooks.UseState(fmtFloat(props.Value))

		onChange := func(newText string) {
			newText = strings.ReplaceAll(newText, ",", ".")
			_, err := strconv.ParseFloat(newText, 32)
			setValid(err == nil)
			setText(newText)
		}

		revert := func() {
			setText(fmtFloat(value))
			setValid(true)
		}

		updateValue := func() {
			text = strings.ReplaceAll(text, ",", ".")
			if f, err := strconv.ParseFloat(text, 32); err == nil {
				// if its a valid float, run the onchange callback
				if props.OnChange != nil && f != float64(props.Value) {
					props.OnChange(float32(f))
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
			setText(fmtFloat(props.Value))
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
