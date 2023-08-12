package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
)

type ObjectEditor struct {
	object.Object
	target object.Object
	GUI    gui.Fragment
}

var _ T = &ObjectEditor{}

func NewObjectEditor(target object.Object) *ObjectEditor {
	return object.New("ObjectEditor", &ObjectEditor{
		Object: object.Ghost(target.Name(), target.Transform()),
		target: target,

		GUI: SidebarFragment(gui.FragmentLast, func() node.T {
			return Inspector(
				target,
				propedit.Transform("transform", target.Transform()),
			)
		}),
	})
}

func (e *ObjectEditor) Target() object.Component { return e.target }

func (e *ObjectEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	return true
}

func (e *ObjectEditor) Actions() []Action {
	return []Action{
		{
			Name: "Move",
			Key:  keys.G,
			Callback: func(mgr ToolManager) {
				mgr.MoveTool(e.target)
			},
		},
		{
			Name: "Rotate",
			Key:  keys.V,
			Callback: func(mgr ToolManager) {
				mgr.RotateTool(e.target)
			},
		},
		{
			Name: "Select Parent",
			Key:  keys.U,
			Callback: func(mgr ToolManager) {
				parent := e.target.Parent()
				if parent == nil {
					return
				}

				editor, hit := object.NewQuery[T]().Where(func(e T) bool {
					return e.Target() == parent
				}).First(object.Root(e))
				if !hit {
					return
				}

				mgr.Select(editor)
			},
		},
	}
}
