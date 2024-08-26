package object

import (
	"reflect"
)

type CreateFn func(Pool) (Component, error)

type Type struct {
	Name   string
	Path   []string
	Create CreateFn

	rtype reflect.Type
}

type Registry map[string]*Type

var types = Registry{}

func typeName(obj any) string {
	t := reflect.TypeOf(obj).Elem()
	return t.PkgPath() + "/" + t.Name()
}

func init() {
	Register[*object](Type{
		Name: "Object",
		Create: func(pool Pool) (Component, error) {
			return Empty(pool, "Object"), nil
		},
		rtype: baseObjectType,
	})
	Register[*component](Type{
		Name:  "Component",
		rtype: baseComponentType,
	})
}

func Register[T any](info Type) {
	var empty T
	kind := typeName(empty)
	info.rtype = reflect.TypeOf(empty).Elem()
	if info.Name == "" {
		info.Name = kind
	}
	types[kind] = &info
}

func Types() Registry {
	return types
}
