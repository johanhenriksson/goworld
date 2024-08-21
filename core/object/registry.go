package object

import (
	"reflect"
)

type CreateFn func(Pool) (Component, error)
type DeserializeFn func(Pool, Decoder) (Component, error)

type TypeInfo struct {
	Name        string
	Path        []string
	Create      CreateFn
	Deserialize DeserializeFn
	rtype       reflect.Type
}

type Registry map[string]TypeInfo

var types = Registry{}

func typeName(obj any) string {
	t := reflect.TypeOf(obj).Elem()
	return t.PkgPath() + "/" + t.Name()
}

func init() {
	Register[*object](TypeInfo{
		Name: "Object",
		Create: func(pool Pool) (Component, error) {
			return Empty(pool, "Object"), nil
		},
		rtype: baseObjectType,
	})
	Register[*component](TypeInfo{
		Name:  "Component",
		rtype: baseComponentType,
	})
}

func Register[T any](info TypeInfo) {
	var empty T
	kind := typeName(empty)
	info.rtype = reflect.TypeOf(empty).Elem()
	if info.Name == "" {
		info.Name = kind
	}
	// if info.Deserialize == nil {
	// 	panic("no deserializer for " + info.Name)
	// }
	types[kind] = info
}

func Types() Registry {
	return types
}
