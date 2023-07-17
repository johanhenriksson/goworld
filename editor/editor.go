package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/physics"
)

type T interface {
	object.Component

	// EditorGUI(object.T) node.T
	Actions() []Action

	Bounds() physics.Shape
}
