package editor

import (
	"os"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/util"
)

func makeToolbar(editor *App) node.T {
	actions := editor.Tools.Actions()
	if len(actions) == 0 {
		return nil
	}

	return rect.New("toolbar", rect.Props{
		Style: rect.Style{
			Color:   color.RGB(0.66, 0.66, 0.66),
			Layout:  style.Row{},
			Padding: style.RectAll(2),
		},
		Children: util.Map(editor.Tools.Actions(), func(action Action) node.T {
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
					action.Callback(editor.Tools)
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
	createItems := make([]menu.ItemProps, 0, len(object.Types()))
	for _, t := range object.Types() {
		info := t // todo: fix in go 1.22
		if info.Create == nil {
			continue
		}
		createItems = append(createItems, menu.ItemProps{
			Key:   "new:" + t.Name,
			Title: "New " + t.Name,
			OnClick: func(e mouse.Event) {
				parent := editor.workspace
				if len(editor.Tools.Selected()) > 0 {
					parent = editor.Tools.Selected()[0].Target().(object.Object)
				}

				thing, err := info.Create()
				if err != nil {
					// todo: handle errors properly
					panic("failed to create " + t.Name + ": " + err.Error())
				}
				object.Attach(parent, thing)
				editor.Refresh()

				if obj, ok := thing.(object.Object); ok {
					position := editor.Player.Transform().Position().Add(editor.Player.Camera.Transform().Forward())
					obj.Transform().SetWorldPosition(position)
					editor.Tools.Select(editor.Lookup(obj))
				} else {
					editor.Tools.Select(editor.Lookup(parent))
				}
			},
		})
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
				Key:   "menu-create",
				Title: "Create",
				Items: createItems,
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
			Padding: style.RectAll(0),
		},
		Children: []node.T{
			ObjectList("scene-graph", ObjectListProps{
				Scene:       editor.workspace,
				EditorRoot:  editor,
				ToolManager: editor.Tools,
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
