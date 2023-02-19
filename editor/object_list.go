package editor

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

type SelectObjectHandler func(object.T)

type ObjectListProps struct {
	Scene       object.T
	EditorRoot  object.T
	ToolManager ToolManager
}

func ObjectList(key string, props ObjectListProps) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Padding: style.RectY(15),
		},
		Children: []node.T{
			ObjectListEntry("scene", ObjectListEntryProps{
				Object: props.Scene,
				OnSelect: func(obj object.T) {
					if !object.Is[*ObjectEditor](obj) {
						// look up an editor instead
						var hit bool
						obj, hit = object.Query[*ObjectEditor]().Where(func(e *ObjectEditor) bool {
							return e.Target() == obj
						}).First(props.EditorRoot)
						if !hit {
							return
						}
					}

					// check if we found something selectable
					if selectable, ok := obj.(Selectable); ok {
						props.ToolManager.Select(selectable)
					}
				},
			}),
		},
	})
}

type ObjectListEntryProps struct {
	Object   object.T
	OnSelect SelectObjectHandler
}

func ObjectListEntry(key string, props ObjectListEntryProps) node.T {
	return node.Component(key, props, func(props ObjectListEntryProps) node.T {
		obj := props.Object
		clr := color.White
		if !obj.Active() {
			clr = color.RGB(0.7, 0.7, 0.7)
		}

		open, setOpen := hooks.UseState(false)
		icon := "+"
		if open {
			icon = "-"
		}

		title := rect.New("title-row", rect.Props{
			Style: rect.Style{
				Layout: style.Row{},
			},
			Children: []node.T{
				label.New("toggle", label.Props{
					Text: icon,
					OnClick: func(e mouse.Event) {
						setOpen(!open)
					},
					Style: label.Style{
						Color: clr,
					},
				}),
				label.New("title", label.Props{
					Text: " " + obj.Name(),
					OnClick: func(e mouse.Event) {
						if props.OnSelect != nil {
							props.OnSelect(obj)
						}
					},
					Style: label.Style{
						Color: clr,
					},
				}),
			},
		})

		nodes := make([]node.T, 0, len(obj.Children())+1)
		nodes = append(nodes, title)

		if open {
			for _, obj := range obj.Children() {
				key := object.Key("object", obj)
				nodes = append(nodes, ObjectListEntry(key, ObjectListEntryProps{
					Object:   obj,
					OnSelect: props.OnSelect,
				}))
			}
		}

		return rect.New(key, rect.Props{
			Style: rect.Style{
				Padding: style.RectXY(5, 3),
			},
			Children: nodes,
		})
	})
}
