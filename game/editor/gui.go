package editor

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func ToolButton(tool Tool, editor *editor) node.T {
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
					Left:   4,
					Right:  4,
				},
				Hover: rect.Hover{},
			},
		},
		OnClick: func(ev mouse.Event) {
			editor.SelectTool(tool)
		},
	})
}

func NewGUI(editor *editor) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: gui.FragmentLast,
		Render: func() node.T {
			return rect.New("editor", rect.Props{
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
