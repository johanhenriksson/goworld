package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func Inspector(target object.Component, propEditors ...node.T) node.T {
	title := "Component: "
	if _, isObject := target.(object.Object); isObject {
		title = "Object: "
	}

	key := object.Key("inspector", target)
	children := make([]node.T, 0, 4+len(propEditors))
	children = append(children, []node.T{
		label.New("title", label.Props{
			Text: title + target.Name(),
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
	}...)
	children = append(children, propEditors...)
	return propedit.Container(key, children)
}
