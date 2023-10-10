package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
	"github.com/johanhenriksson/goworld/physics"
)

type Tool interface {
	object.Component
	CanDeselect() bool
	ToolMouseEvent(e mouse.Event, hover physics.RaycastHit)
}

const ToolLayer = physics.Mask(2)

type Action struct {
	Name     string
	Icon     icon.Icon
	Key      keys.Code
	Modifier keys.Modifier
	Callback func(*App)
}
