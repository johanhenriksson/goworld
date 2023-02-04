package chunk

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func ToolButton(tool Tool, editor *edit) node.T {
	bg := color.DarkGrey
	if tool == editor.Tool {
		bg = color.RGB(0.7, 0.7, 0.7)
	}
	return button.New(tool.Name(), button.Props{
		Text: tool.Name(),
		Style: button.Style{
			Bg: rect.Style{
				Color:   bg,
				Padding: style.Px(4),
				Margin: style.Rect{
					Bottom: 4,
				},
				Hover: rect.Hover{},
			},
		},
		OnClick: func(ev mouse.Event) {
			if editor.Tool != tool {
				editor.SelectTool(tool)
			} else {
				editor.SelectTool(nil)
			}
		},
	})
}

func NewGUI(editor *edit) gui.Fragment {
	key := fmt.Sprintf("chunk:%s", editor.mesh.meshdata.Key())
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: gui.FragmentLast,
		Render: func() node.T {
			return rect.New(key, rect.Props{
				Children: []node.T{
					palette.New("palette", palette.Props{
						Palette: color.DefaultPalette,
						OnPick: func(clr color.T) {
							editor.SelectColor(clr)
						},
					}),

					ToolButton(editor.PlaceTool, editor),
					ToolButton(editor.EraseTool, editor),
					ToolButton(editor.ReplaceTool, editor),
					ToolButton(editor.SampleTool, editor),
				},
			})
		},
	})
}

func NewMenu(editor *edit) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "main-menu",
		Position: gui.FragmentLast,
		Render: func() node.T {
			return menu.Item("editor-menu", menu.ItemProps{
				Key:      "menu-editor",
				Title:    "Editor",
				Style:    menu.DefaultStyle,
				OpenDown: true,
				Items: []menu.ItemProps{
					{
						Key:   "tool",
						Title: "Tool",
						Items: []menu.ItemProps{
							{
								Key:     "tool-place",
								Title:   "Place",
								OnClick: func(e mouse.Event) { editor.SelectTool(editor.PlaceTool) },
							},
							{
								Key:     "tool-erase",
								Title:   "Erase",
								OnClick: func(e mouse.Event) { editor.SelectTool(editor.EraseTool) },
							},
							{
								Key:     "tool-replace",
								Title:   "Replace",
								OnClick: func(e mouse.Event) { editor.SelectTool(editor.ReplaceTool) },
							},
							{
								Key:     "tool-sample",
								Title:   "Sample",
								OnClick: func(e mouse.Event) { editor.SelectTool(editor.SampleTool) },
							},
						},
					},
				},
			})
		},
	})
}
