package propedit

import (
	"reflect"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
)

type PropEditor func(key, name string, prop object.GenericProp) node.T

var registry map[reflect.Type]PropEditor = make(map[reflect.Type]PropEditor, 100)

func For[T object.PropValue](prop object.Property[T]) PropEditor {
	var empty T
	t := reflect.TypeOf(empty)
	return ForType(t)
}

func ForType(t reflect.Type) PropEditor {
	return registry[t]
}

func Register[T object.PropValue](editor PropEditor) {
	var empty T
	t := reflect.TypeOf(empty)
	registry[t] = editor
}
