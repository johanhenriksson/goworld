package gui

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
)

type Fragment interface {
	object.Component
	Render() node.T
	Slot() string
}

type fragment struct {
	object.Component
	renderer node.RenderFunc
	slot     string
}

func NewFragment(slot string, render node.RenderFunc) Fragment {
	return &fragment{
		Component: object.NewComponent(),
		renderer:  render,
		slot:      slot,
	}
}

func (f *fragment) Name() string { return "UIFragment" }
func (f *fragment) Slot() string { return f.slot }

func (f *fragment) Render() node.T {
	return f.renderer()
}
