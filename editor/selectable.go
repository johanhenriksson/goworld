package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
)

type Selectable interface {
	object.T
	Select(mouse.Event, collider.T)
	Deselect(mouse.Event) bool
}

type selectable struct {
	object.T
	Bounds     collider.T
	OnSelect   func()
	OnDeselect func() bool
}

func NewSelectable(bounds collider.T, onSelect func(), onDeselect func() bool) Selectable {
	return object.New(&selectable{
		Bounds:     bounds,
		OnSelect:   onSelect,
		OnDeselect: onDeselect,
	})
}

func (g *selectable) Select(e mouse.Event, collider collider.T) {
	if g.OnSelect != nil {
		g.OnSelect()
	}
}

func (g *selectable) Deselect(e mouse.Event) bool {
	if g.OnDeselect != nil {
		return g.OnDeselect()
	}
	return true
}
