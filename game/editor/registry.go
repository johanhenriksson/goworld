package editor

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Constructor func(*Context, object.T) T

var editors = map[reflect.Type]Constructor{}

// Register an Editor constructor for a given object type
func Register[K object.T, E T](obj object.T, constructor func(*Context, K) E) {
	t := reflect.TypeOf(obj).Elem()
	editors[t] = func(ctx *Context, obj object.T) T {
		k := obj.(K)
		return constructor(ctx, k)
	}
}

func ConstructEditors(ctx *Context, obj object.T, mv *mover.T) object.T {
	editNode := object.Ghost(obj)

	t := reflect.TypeOf(obj).Elem()
	if construct, exists := editors[t]; exists {
		log.Println("creating editor for", obj.Name())
		object.Attach(editNode, NewObjectEditor(
			obj,
			collider.Box{
				Center: vec3.New(16, 16, 16),
				Size:   vec3.New(32, 32, 32),
			},
			construct(ctx, obj),
		))
	}

	for _, child := range obj.Children() {
		childNode := ConstructEditors(ctx, child, mv)
		object.Attach(editNode, childNode)
	}
	return editNode
}
