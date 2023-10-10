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
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/util"
)

func makeToolbar(editor *App) node.T {
	actions := editor.Actions()
	if len(actions) == 0 {
		return nil
	}

	return rect.New("toolbar", rect.Props{
		Style: rect.Style{
			Color:   color.RGB(0.66, 0.66, 0.66),
			Layout:  style.Row{},
			Padding: style.RectAll(2),
		},
		Children: util.Map(actions, func(action Action) node.T {
			return button.New(action.Name, button.Props{
				// Text: action.Name,
				Icon: action.Icon,
				Style: button.Style{
					TextColor: color.Black,
					BgColor:   color.RGB(0.76, 0.76, 0.76),
					Padding:   style.RectXY(8, 4),
					Margin:    style.Px(2),
					Radius:    style.Px(4),
				},
				OnClick: func(e mouse.Event) {
					action.Callback(editor)
				},
			})
		}),
	})
}

func MakeGUI(editor *App) gui.Manager {
	return gui.New(func() node.T {
		return rect.New("gui", rect.Props{
			Children: []node.T{
				// menu & toolbar
				makeMenu(editor),
				makeToolbar(editor),

				// main content area
				rect.New("gui-main", rect.Props{
					Style: rect.Style{
						Grow:   style.Grow(1),
						Layout: style.Row{},
					},
					Children: []node.T{
						// left sidebar: scene graph
						makeSidebarLeft(editor),

						// middle: scene view + asset browser
						rect.New("gui-content", rect.Props{
							Style: rect.Style{
								Grow:   style.Grow(1),
								Layout: style.Column{},
							},
							Children: []node.T{
								// middle top: scene view
								rect.New("viewport", rect.Props{
									Style: rect.Style{
										Grow: style.Grow(1),
									},
								}),
								// middle bottom: asset browser
								// rect.New("asset-browser", rect.Props{
								// 	Style: rect.Style{
								// 		Height: style.Pct(20),
								// 		Color:  color.RGBA(0.1, 0.1, 0.11, 0.85),
								// 	},
								// }),
							},
						}),

						// right sidebar: object property editors
						makeSidebarRight(editor),
					},
				}),
			},
		})
	})
}

func makeMenu(editor *App) node.T {
	attach := func(cmp object.Component) {
		if len(editor.Selected()) < 1 {
			log.Println("no selection?")
			return
		}
		obj := editor.Selected()[0].Target().(object.Object)
		object.Attach(obj, cmp)
		editor.Refresh()
		editor.Select(editor.Lookup(obj))
	}
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
						Key:   "file-new",
						Title: "New",
					},
					{
						Key:   "file-open",
						Title: "Open...",
					},
					{
						Key:   "file-save",
						Title: "Save...",
						OnClick: func(e mouse.Event) {
						},
					},
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
							editor.Select(editor.Lookup(obj))
						},
					},
					{
						Key:   "add-point-light",
						Title: "Add Point Light",
						OnClick: func(e mouse.Event) {
							attach(light.NewPoint(light.PointArgs{
								Color:     color.Purple,
								Range:     10,
								Intensity: 3,
							}))
						},
					},
					{
						Key:   "add-dir-light",
						Title: "Add Directional Light",
						OnClick: func(e mouse.Event) {
							attach(light.NewDirectional(light.DirectionalArgs{
								Color:     color.Purple,
								Intensity: 3,
								Shadows:   true,
								Cascades:  3,
							}))
						},
					},
					{
						Key:   "add-rigidbody-light",
						Title: "Add Rigidbody",
						OnClick: func(e mouse.Event) {
							attach(physics.NewRigidBody(0))
						},
					},
				},
			},
		},
	})
}

func makeSidebarLeft(editor *App) node.T {
	return rect.New("sidebar-left", rect.Props{
		OnMouseDown: gui.ConsumeMouse,
		Style: rect.Style{
			Layout:  style.Column{},
			Grow:    style.Grow(0),
			Width:   style.Px(200),
			Height:  style.Pct(100),
			Color:   color.RGBA(0.1, 0.1, 0.11, 0.85),
			Padding: style.RectAll(10),
		},
		Children: []node.T{
			ObjectList("scene-graph", ObjectListProps{
				Scene:  editor.workspace,
				Editor: editor,
			}),
		},
	})
}

func makeSidebarRight(editor *App) node.T {
	return rect.New("sidebar-right", rect.Props{
		Style: rect.Style{
			Layout:  style.Column{},
			Color:   color.RGBA(0.1, 0.1, 0.11, 0.85),
			Grow:    style.Grow(0),
			Width:   style.Px(200),
			Height:  style.Pct(100),
			Padding: style.RectAll(10),
		},
		Children: []node.T{
			// property editor placeholder
			rect.New("sidebar-right:property-editor", rect.Props{}),
		},
	})
}

func PropertyEditorFragment(position gui.FragmentPosition, render node.RenderFunc) gui.Fragment {
	return gui.NewFragment(gui.FragmentArgs{
		Slot:     "sidebar-right:property-editor",
		Position: position,
		Render:   render,
	})
}
