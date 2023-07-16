package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
)

type T interface {
	object.Component

	// EditorGUI(object.T) node.T
	Actions() []Action
}
