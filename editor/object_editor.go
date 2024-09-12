package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
)

type ObjectEditor struct {
	Object
	target Object
	GUI    gui.Fragment
}

var _ T = &ObjectEditor{}

func NewObjectEditor(pool Pool, target Object) *ObjectEditor {
	props := Properties(target)
	editors := make([]node.T, 2, len(props)+2)

	return NewObject(pool, "ObjectEditor", &ObjectEditor{
		Object: Ghost(pool, target.Name(), target.Transform()),
		target: target,

		GUI: PropertyEditorFragment(pool, gui.FragmentLast, func() node.T {
			editors = editors[:2]
			editors[0] = propedit.BoolField("enabled", "Enabled", propedit.BoolProps{
				Value: target.Enabled(),
				OnChange: func(b bool) {
					Toggle(target, b)
				},
			})
			editors[1] = propedit.Transform("transform", target.Transform())

			for _, prop := range props {
				if editor := propedit.ForType(prop.Type()); editor != nil {
					editors = append(editors, editor(prop.Key, prop.Name, prop))
				}
			}

			return Inspector(
				target,
				editors...,
			)
		}),
	})
}

func (e *ObjectEditor) Target() Component { return e.target }

func (e *ObjectEditor) Select(ev mouse.Event) {
	Enable(e.GUI)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	Disable(e.GUI)
	return true
}

func (e *ObjectEditor) Actions() []Action {
	return []Action{
		{
			Name: "Move",
			Key:  keys.G,
			Icon: icon.IconMove,
			Callback: func(mgr *ToolManager) {
				mgr.MoveTool(e.target)
			},
		},
		{
			Name: "Rotate",
			Icon: icon.IconRotate,
			Key:  keys.V,
			Callback: func(mgr *ToolManager) {
				mgr.RotateTool(e.target)
			},
		},
		{
			Name: "Select Parent",
			Icon: icon.IconVerticalAlignTop,
			Key:  keys.U,
			Callback: func(mgr *ToolManager) {
				parent := e.target.Parent()
				if parent == nil {
					return
				}

				editor, hit := NewQuery[T]().Where(func(e T) bool {
					return e.Target() == parent
				}).First(Root(e))
				if !hit {
					return
				}

				mgr.Select(editor)
			},
		},
	}
}

func (e *ObjectEditor) Update(scene Component, dt float32) {
	e.Object.Update(scene, dt)
	if updatable, ok := e.target.(EditorUpdater); ok {
		updatable.EditorUpdate(scene, dt)
	}
}
