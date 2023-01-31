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
		objEditor := construct(ctx, obj)
		objEditor.SetActive(false)
		object.Attach(editNode, objEditor)
		object.Attach(editNode, NewSelectable(
			collider.NewBox(collider.Box{
				Center: vec3.New(4, 4, 4),
				Size:   vec3.New(8, 8, 8),
			}),
			func() {
				objEditor.SetActive(true)
				mv.SetTarget(obj.Transform())
				mv.SetActive(true)
			},
			func() bool {
				if objEditor.CanDeselect() {
					objEditor.SetActive(false)
					mv.SetTarget(nil)
					mv.SetActive(false)
					return true
				}
				return false
			}))
	}

	for _, child := range obj.Children() {
		childNode := ConstructEditors(ctx, child, mv)
		object.Attach(editNode, childNode)
	}
	return editNode
}
