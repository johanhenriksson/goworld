package editor

import (
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

func MakeGUI(render renderer.T) gui.Manager {
	return gui.New(func() node.T {
		return rect.New("gui", rect.Props{
			Children: []node.T{
				makeMenu(),
				rect.New("gui-main", rect.Props{
					Style: rect.Style{
						Grow: style.Grow(1),
					},
					Children: []node.T{
						makeSidebar(),
					},
				}),
			},
		})
	})
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

func makeSidebar() node.T {
	return rect.New("sidebar", rect.Props{
		OnMouseDown: gui.ConsumeMouse,
		Style: rect.Style{
			Layout:  style.Column{},
			Grow:    style.Grow(0),
			Width:   style.Px(200),
			Height:  style.Pct(100),
			Color:   color.RGBA(0.1, 0.1, 0.11, 0.85),
			Padding: style.RectAll(15),
		},
		Children: []node.T{
			rect.New("logo-container", rect.Props{
				Style: rect.Style{
					Padding: style.Rect{Bottom: 15},
				},
				Children: []node.T{
					image.New("logo", image.Props{
						Image: texture.PathRef("textures/shit_logo.png"),
						Style: image.Style{
							Width:  style.Pct(100),
							Height: style.Auto{},
						},
					}),
				},
			}),

			// content placeholder
			rect.New("sidebar:content", rect.Props{}),
		},
	})
}

func SidebarFragment(position gui.FragmentPosition, render node.RenderFunc) object.T {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: position,
		Render:   render,
	})
}
