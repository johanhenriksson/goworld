package editor

import (
	"fmt"
	"os"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

func makeGui(r renderer.T, scene object.T) {
	object.Attach(scene, gui.New(func() node.T {
		return rect.New("gui", rect.Props{
			Children: []node.T{
				makeMenu(),
				rect.New("gui-main", rect.Props{
					Style: rect.Style{
						Grow: style.Grow(1),
					},
					Children: []node.T{
						makeSidebar(scene, r),
					},
				}),
			},
		})
	}))
}

func makeMenu() node.T {
	return menu.Menu("gui-menu", menu.Props{
		Style: menu.Style{
			Color:      color.RGB(0.76, 0.76, 0.76),
			HoverColor: color.RGB(0.85, 0.85, 0.85),
			TextColor:  color.Black,
		},

		Items: []menu.ItemProps{
			{
				Key:   "menu-file",
				Title: "File",
				Items: []menu.ItemProps{
					{
						Key:   "file-exit",
						Title: "Exit",
						OnClick: func(e mouse.Event) {
							os.Exit(0)
						},
					},
				},
			},
			{
				Key:   "menu-edit",
				Title: "Edit",
				Items: []menu.ItemProps{
					{
						Key:   "edit-undo",
						Title: "Undo",
					},
					{
						Key:   "edit-redo",
						Title: "Redo",
					},
				},
			},
		},
	})
}

func makeSidebar(scene object.T, r renderer.T) node.T {
	return rect.New("sidebar", rect.Props{
		OnMouseDown: gui.ConsumeMouse,
		Style: rect.Style{
			Layout: style.Column{},
			Width:  style.Pct(15),
			Height: style.Pct(100),
			Color:  color.RGBA(0.1, 0.1, 0.11, 0.85),
		},
		Children: []node.T{
			image.New("logo", image.Props{
				Image: texture.PathRef("textures/shit_logo.png"),
				Style: image.Style{
					Width:  style.Pct(100),
					Height: style.Auto{},
				},
			}),

			// content placeholder
			rect.New("sidebar:content", rect.Props{}),

			ObjectListEntry("scene-graph", ObjectListEntryProps{
				Object: scene,
				OnSelect: func(obj object.T) {
					fmt.Println("selected", obj.Name())
				},
			}),
		},
	})
}
