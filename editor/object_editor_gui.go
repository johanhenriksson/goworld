package editor

import (
	"log"
	"strconv"
	"strings"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func SidebarItem(key string, children []node.T) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Layout:     style.Row{},
			AlignItems: style.AlignStart,
			Width:      style.Pct(100),
			Padding:    style.RectY(2),
		},
		Children: children,
	})
}

func objectEditorGui(target object.T) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: gui.FragmentLast,
		Render: func() node.T {
			key := object.Key("editor", target)
			return rect.New(key, rect.Props{
				Style: rect.Style{
					Layout:     style.Column{},
					AlignItems: style.AlignStart,
					Width:      style.Pct(100),
				},
				Children: []node.T{
					label.New("name", label.Props{
						Text: target.Name(),
					}),
					SidebarItem("position", []node.T{
						Vec3Editor("position", Vec3EditorProps{
							Value: target.Transform().Position(),
							OnChange: func(pos vec3.T) {
								log.Println("update position", pos)
								target.Transform().SetPosition(pos)
							},
						}),
					}),
					SidebarItem("test", []node.T{
						FloatEditor("test", FloatEditorProps{Label: "F", Value: 13.37}),
					}),
				},
			})
		},
	})
}

type FloatEditorProps struct {
	Label    string
	Value    float32
	OnChange func(float32)
}

func fmtFloat(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func FloatEditor(key string, props FloatEditorProps) node.T {
	return node.Component(key, props, func(props FloatEditorProps) node.T {
		value, setValue := hooks.UseState(props.Value)
		valid, setValid := hooks.UseState(true)
		text, setText := hooks.UseState(fmtFloat(props.Value))

		onChange := func(newText string) {
			newText = strings.ReplaceAll(newText, ",", ".")
			_, err := strconv.ParseFloat(newText, 32)
			setValid(err == nil)
			setText(newText)
		}

		updateValue := func() {
			text = strings.ReplaceAll(text, ",", ".")
			if f, err := strconv.ParseFloat(text, 32); err == nil {
				// if its a valid float, run the onchange callback
				if props.OnChange != nil && f != float64(props.Value) {
					props.OnChange(float32(f))
				}
			} else {
				// otherwise revert to the previous value
				setText(fmtFloat(value))
				setValid(true)
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
					OnKeyDown: func(e keys.Event) {
						if e.Code() == keys.Enter {
							updateValue()
							// e.Consume()
						}
						if e.Code() == keys.Escape {
							// revert
							setText(fmtFloat(value))
							// e.Consume()
						}
					},
				}),
			},
		})
	})
}

type Vec3EditorProps struct {
	Value    vec3.T
	OnChange func(vec3.T)
}

func Vec3Editor(key string, props Vec3EditorProps) node.T {
	return node.Component(key, props, func(props Vec3EditorProps) node.T {
		return rect.New(key, rect.Props{
			Style: rect.Style{
				Layout:         style.Row{},
				AlignItems:     style.AlignCenter,
				JustifyContent: style.JustifySpaceBetween,
				Grow:           style.Grow(1),
				Shrink:         style.Shrink(1),
				Basis:          style.Pct(100),
			},
			Children: []node.T{
				FloatEditor("x", FloatEditorProps{
					Label: "X",
					Value: props.Value.X,
					OnChange: func(x float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(x, props.Value.Y, props.Value.Z))
						}
					},
				}),
				FloatEditor("y", FloatEditorProps{
					Label: "Y",
					Value: props.Value.Y,
					OnChange: func(y float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(props.Value.X, y, props.Value.Z))
						}
					},
				}),
				FloatEditor("z", FloatEditorProps{
					Label: "Z",
					Value: props.Value.Z,
					OnChange: func(z float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(props.Value.X, props.Value.Y, z))
						}
					},
				}),
			},
		})
	})
}
