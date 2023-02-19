package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
)

type ObjectEditor struct {
	object.T

	Bounds collider.T
	GUI    gui.Fragment
	Custom T

	target object.T
}

func NewObjectEditor(target object.T, bounds collider.Box, editor T) *ObjectEditor {
	var boundsCol collider.T
	if editor != nil {
		editor.SetActive(false)
		boundsCol = collider.NewBox(bounds)
	}
	edit := object.New(&ObjectEditor{
		T:      object.Ghost(target),
		target: target,
		Custom: editor,
		Bounds: boundsCol,
		GUI:    objectEditorGui(target),
	})
	edit.GUI.SetActive(false)
	return edit
}

var _ Selectable = &ObjectEditor{}

func (e *ObjectEditor) Select(ev mouse.Event, collider collider.T) {
	if e.Custom != nil {
		e.Custom.SetActive(true)
	}
	e.GUI.SetActive(true)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	if e.Custom != nil {
		e.Custom.SetActive(false)
	}
	e.GUI.SetActive(false)
	return true
}

func (e *ObjectEditor) Target() object.T {
	return e.target
}

func (e *ObjectEditor) Actions() []Action {
	actions := []Action{
		{
			Name: "Move",
			Key:  keys.G,
			Callback: func(mgr ToolManager) {
				mgr.MoveTool(e.target)
			},
		},
	}
	if e.Custom != nil {
		actions = append(actions, e.Custom.Actions()...)
	}
	return actions
}
