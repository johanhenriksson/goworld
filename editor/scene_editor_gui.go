package editor

import (
	"log"
	"os"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/image"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

func MakeGUI(editor *Editor) gui.Manager {
	return gui.New(func() node.T {
		return rect.New("gui", rect.Props{
			Children: []node.T{
				makeMenu(editor),
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

func makeMenu(editor *Editor) node.T {
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
							// it'd be cool if there was a decent way of doing things like exiting
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
			{
				Key:   "menu-object",
				Title: "Object",
				Items: []menu.ItemProps{
					{
						Key:   "menu-new-object",
						Title: "New",
						OnClick: func(e mouse.Event) {
							position := editor.Player.Transform().Position().Add(editor.Player.Transform().Forward())
							obj := object.Builder(object.Empty("New Object")).
								Position(position).
								Create()
							object.Attach(editor.workspace, obj)
							editor.Refresh()
							editor.Tools.Select(editor.Lookup(obj))
						},
					},
					{
						Key:   "add-point-light",
						Title: "Add Point Light",
						OnClick: func(e mouse.Event) {
							if len(editor.Tools.Selected()) < 1 {
								log.Println("no selection?")
								return
							}
							obj := editor.Tools.Selected()[0].Target().(object.Object)
							object.Attach(obj, light.NewPoint(light.PointArgs{
								Color:     color.Purple,
								Range:     10,
								Intensity: 3,
							}))
							editor.Refresh()
							editor.Tools.Select(editor.Lookup(obj))
						},
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

func SidebarFragment(position gui.FragmentPosition, render node.RenderFunc) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar:content",
		Position: position,
		Render:   render,
	})
}
