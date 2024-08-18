package object

import (
	"reflect"
)

type CreateFn func(Context) (Component, error)
type DeserializeFn func(Context, Decoder) (Component, error)

type TypeInfo struct {
	Name        string
	Path        []string
	Create      CreateFn
	Deserialize DeserializeFn
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
		Deserialize: DeserializeObject,
	})
}

func Register[T Serializable](info TypeInfo) {
	var empty T
	kind := typeName(empty)
	if info.Name == "" {
		t := reflect.TypeOf(empty).Elem()
		info.Name = t.Name()
	}
	if info.Deserialize == nil {
		panic("no deserializer for " + info.Name)
	}
	types[kind] = info
}

func Types() Registry {
	return types
}
