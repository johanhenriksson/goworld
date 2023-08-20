package editor

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
)

type ComponentEditor struct {
	object.Object
	target object.Component
	GUI    gui.Fragment
}

var _ T = &ComponentEditor{}

func NewComponentEditor(target object.Component) *ComponentEditor {
	props := object.Properties(target)
	editors := make([]node.T, 0, len(props))
	return object.New("ComponentEditor", &ComponentEditor{
		Object: object.Ghost(target.Name(), target.Transform()),
		target: target,

		GUI: PropertyEditorFragment(gui.FragmentLast, func() node.T {
			editors = editors[:0]
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

func (e *ComponentEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
}

func (e *ComponentEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	return true
}

func (e *ComponentEditor) Target() object.Component { return e.target }
func (e *ComponentEditor) Actions() []Action        { return nil }
