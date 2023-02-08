package editor

import (
	"fmt"

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
					Vec3Editor("position", Vec3EditorProps{
						Value: target.Transform().WorldPosition(),
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

func FloatEditor(key string, props FloatEditorProps) node.T {
	const FloatFmt = "%.f"
	return node.Component(key, props, func(props FloatEditorProps) node.T {
		value, setValue := hooks.UseState(fmt.Sprintf(FloatFmt, props.Value))
		onChange := func(v string) {
			// parse, validate
			if props.OnChange != nil {
				props.OnChange(0)
			}
			setValue(v)
		}
		return rect.New(key, rect.Props{
			Style: rect.Style{
				Layout:     style.Row{},
				AlignItems: style.AlignCenter,
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
					Style:    textbox.DefaultStyle,
					Text:     value,
					OnChange: onChange,
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
				Width:          style.Pct(100),
			},
			Children: []node.T{
				FloatEditor("x", FloatEditorProps{
					Label: "X",
					Value: props.Value.X,
				}),
				FloatEditor("y", FloatEditorProps{
					Label: "Y",
					Value: props.Value.Y,
				}),
				FloatEditor("z", FloatEditorProps{
					Label: "Z",
					Value: props.Value.Z,
				}),
			},
		})
	})
}
