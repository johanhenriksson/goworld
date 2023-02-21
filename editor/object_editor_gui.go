package editor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/gui/widget/textbox"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
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

func FloatEditor(key string, props FloatEditorProps) node.T {
	const FloatFmt = "%.f"
	return node.Component(key, props, func(props FloatEditorProps) node.T {
		value, setValue := hooks.UseState(fmt.Sprintf(FloatFmt, props.Value))
		onChange := func(v string) {
			if props.OnChange != nil {
				// if its a valid float, run the onchange callback
				// todo: OnChange should probably only be called if the user
				//       leaves the input, or presses enter.
				if f, err := strconv.ParseFloat(v, 32); err == nil {
					props.OnChange(float32(f))
				} else {
					log.Println("invalid float", err, "from", v)
				}
			}
			setValue(v)
		}
		// how to properly observe changing values with hooks?
		// this hook would make sure the label is updated any time the prop changes,
		// even due to an external change (or OnChange effects)
		hooks.UseEffect(func() { setValue(fmt.Sprintf(FloatFmt, props.Value)) }, props.Value)

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
					Text:     value,
					OnChange: onChange,
					Style: textbox.Style{
						Text: label.Style{
							Color: color.Black,
						},
						Bg: rect.Style{
							Color:   color.White,
							Padding: style.RectXY(4, 2),
							Basis:   style.Pct(100),
							Shrink:  style.Shrink(1),
							Grow:    style.Grow(1),
						},
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
