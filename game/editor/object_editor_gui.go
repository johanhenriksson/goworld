package editor

import (
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

func objectEditorGui() gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: gui.FragmentLast,
		Render: func() node.T {
			return rect.New("object-editor", rect.Props{
				Children: []node.T{
					label.New("hello", label.Props{
						Text: "object editor",
					}),
				},
			})
		},
	})
}
