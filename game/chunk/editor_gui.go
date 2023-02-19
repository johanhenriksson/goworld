package chunk

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

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
			return menu.Item("chunk-menu", menu.ItemProps{
				Key:      "menu-chunk",
				Title:    "Chunk",
				Style:    menu.DefaultStyle,
				OpenDown: true,
				Items: []menu.ItemProps{
					{
						Key:     "action-save",
						Title:   "Save",
						OnClick: func(e mouse.Event) { editor.saveChunkDialog() },
					},
					{
						Key:     "action-clear",
						Title:   "Clear",
						OnClick: func(e mouse.Event) { editor.clearChunk() },
					},
				},
			})
		},
	})
}
