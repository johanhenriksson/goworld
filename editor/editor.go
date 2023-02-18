package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
)

type T interface {
	object.T

	CanDeselect() bool
}
