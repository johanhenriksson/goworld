package query

import (
	"github.com/johanhenriksson/goworld/core/object"
)

func Is[K any](c object.Component) bool {
	_, ok := c.(K)
	return ok
}
