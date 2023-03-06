package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func EditorRow(key, title string, children []node.T) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Layout:     style.Column{},
			AlignItems: style.AlignStart,
			Width:      style.Pct(100),
			Padding:    style.RectY(4),
		},
		Children: []node.T{
			label.New("label", label.Props{
				Text: title,
				Style: label.Style{
					Color: color.White,
				},
			}),
			rect.New("editor", rect.Props{
				Style: rect.Style{
					Layout:     style.Row{},
					AlignItems: style.AlignStart,
					Width:      style.Pct(100),
					Padding:    style.RectY(2),
				},
				Children: children,
			}),
		},
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
					label.New("title", label.Props{
						Text: target.Name(),
						Style: label.Style{
							Font: style.Font{
								Size: 16,
							},
						},
					}),
					rect.New("underline", rect.Props{
						Style: rect.Style{
							Width: style.Pct(100),
							Border: style.Border{
								Width: style.Px(0.5),
								Color: color.White,
							},
						},
					}),
					EditorRow("position", "Position", []node.T{
						propedit.Vec3("position", propedit.Vec3Props{
							Value: target.Transform().Position(),
							OnChange: func(pos vec3.T) {
								target.Transform().SetPosition(pos)
							},
						}),
					}),
					EditorRow("rotation", "Rotation", []node.T{
						propedit.Vec3("rotation", propedit.Vec3Props{
							Value: target.Transform().Rotation(),
							OnChange: func(rot vec3.T) {
								target.Transform().SetRotation(rot)
							},
						}),
					}),
				},
			})
		},
	})
}
