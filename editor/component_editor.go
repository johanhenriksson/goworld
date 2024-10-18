package editor

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
)

type ComponentEditor struct {
	Object
	target Component
	GUI    gui.Fragment
}

var _ T = &ComponentEditor{}

func NewComponentEditor(pool Pool, target Component) *ComponentEditor {
	props := Properties(target)
	editors := make([]node.T, 0, len(props))

	return NewObject(pool, "ComponentEditor", &ComponentEditor{
		Object: Ghost(pool, target.Name(), target.Transform()),
		target: target,

		GUI: PropertyEditorFragment(pool, gui.FragmentLast, func() node.T {
			clear(editors)
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
	Enable(e.GUI)
}

func (e *ComponentEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	Disable(e.GUI)
	return true
}

func (e *ComponentEditor) Target() Component { return e.target }
func (e *ComponentEditor) Actions() []Action { return nil }

func (e *ComponentEditor) Update(scene Component, dt float32) {
	e.Object.Update(scene, dt)
	if updatable, ok := e.target.(EditorUpdater); ok {
		updatable.EditorUpdate(scene, dt)
	}
}
