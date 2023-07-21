package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
)

type ComponentEditor struct {
	object.Object
	target object.Component
	GUI    gui.Fragment
}

func NewComponentEditor(target object.Component) *ComponentEditor {
	props := object.Properties(target)
	editors := make([]node.T, 0, len(props))
	return object.New("ComponentEditor", &ComponentEditor{
		target: target,
		GUI: SidebarFragment(gui.FragmentLast, func() node.T {
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

func (d *ComponentEditor) Actions() []Action     { return nil }
func (d *ComponentEditor) Bounds() physics.Shape { return nil }
