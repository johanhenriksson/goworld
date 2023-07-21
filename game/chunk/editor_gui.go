package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/render/color"
)

func NewGUI(e *edit, target *Mesh) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: gui.FragmentLast,
		Render: func() node.T {
			return editor.Inspector(target,
				// extend the default inspector with a color picker palette
				palette.New("palette", palette.Props{
					Palette: color.DefaultPalette,
					OnPick: func(clr color.T) {
						e.SelectColor(clr)
					},
				}),
			)
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
					{
						Key:   "new",
						Title: "New",
						OnClick: func(e mouse.Event) {
							eye := editor.Camera.Transform().WorldPosition()
							offset := editor.Camera.Transform().Forward().Scaled(3)
							object.Builder(object.Empty("New Chunk")).
								Attach(NewMesh(New(8, 0, 0))).
								Position(eye.Add(offset)).
								Parent(editor.Context.Scene).
								Create()
						},
					},
				},
			})
		},
	})
}
