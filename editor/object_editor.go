package editor

import (
	"log"

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

func NewObjectEditor(target object.Object) *ObjectEditor {
	return object.New("ObjectEditor", &ObjectEditor{
		Object: object.Ghost("Ghost:"+target.Name(), target.Transform()),
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
	log.Println("enable gui", e.target.Name())
	object.Enable(e.GUI)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	log.Println("disable gui", e.target.Name())
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
	}
}
