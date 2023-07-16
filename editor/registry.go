package editor

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Constructor func(*Context, object.Component) T

var editors = map[reflect.Type]Constructor{}

// Register an Editor constructor for a given object type
func Register[K object.Component, E T](obj object.Component, constructor func(*Context, K) E) {
	t := reflect.TypeOf(obj).Elem()
	editors[t] = func(ctx *Context, obj object.Component) T {
		k := obj.(K)
		return constructor(ctx, k)
	}
}

func ConstructEditors(ctx *Context, current object.Component, target object.Component) object.Component {
	var editor *ObjectEditor

	if current != nil {
		editor = current.(*ObjectEditor)
		if editor.target != target {
			panic("unexpected editor target")
		}
	} else {
		t := reflect.TypeOf(target).Elem()
		var customEditor T
		if construct, exists := editors[t]; exists {
			log.Println("creating custom editor for", target.Name())
			customEditor = construct(ctx, target)
		} else {
			log.Println("creating object editor for", target.Name())
		}
		editor = NewObjectEditor(
			target,
			// todo: how to figure out the bounds?
			// todo: how to deal with other shapes?
			collider.Box{
				Center: vec3.New(0.5, 0.5, 0.5),
				Size:   vec3.New(1, 1, 1),
			},
			customEditor,
		)
	}

	existingEditors := map[object.Component]*ObjectEditor{}
	for _, child := range editor.Children() {
		childEdit, isEdit := child.(*ObjectEditor)
		if !isEdit {
			continue
		}
		existingEditors[childEdit.target] = childEdit
	}

	for _, child := range object.Children(target) {
		current, exists := existingEditors[child]
		if exists {
			ConstructEditors(ctx, current, child)
			delete(existingEditors, child)
		} else {
			childEdit := ConstructEditors(ctx, nil, child)
			object.Attach(editor, childEdit)
		}
	}

	// any remaining editor is no longer used
	for _, childEdit := range existingEditors {
		log.Println("deleting unused editor for", childEdit.target.Name())
		object.Detach(childEdit)
	}

	return editor
}
