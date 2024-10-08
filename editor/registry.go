package editor

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
)

type Constructor func(*Context, object.Component) T

var editors = map[reflect.Type]Constructor{}

func typeOf[K object.Component](obj K) reflect.Type {
	t := reflect.TypeOf(obj)
	if t == nil {
		t = reflect.TypeOf((*K)(nil))
	}
	return t.Elem()
}

// Register an Editor constructor for a given object type
func RegisterEditor[K object.Component, E T](obj K, constructor func(*Context, K) E) {
	t := typeOf[K](obj)
	log.Println("register editor for", t.Name())
	editors[t] = func(ctx *Context, obj object.Component) T {
		k := obj.(K)
		return constructor(ctx, k)
	}
}

func ConstructEditors(ctx *Context, current object.Component, target object.Component) object.Component {
	var editor T

	if current != nil {
		editor = current.(T)
		if editor.Target() != target {
			panic("unexpected editor target")
		}
	} else {
		t := typeOf(target)
		if construct, exists := editors[t]; exists {
			log.Println("creating custom editor for", target.Name(), "of type", t.Name())
			editor = construct(ctx, target)
		} else {
			log.Println("creating default editor for", target.Name(), "of type", t.Name())
			if obj, isObject := target.(object.Object); isObject {
				// use default object editor
				editor = NewObjectEditor(ctx.Objects, obj)
			} else {
				// use default component editor
				editor = NewComponentEditor(ctx.Objects, target)
			}
		}

		// start deselected
		editor.Deselect(mouse.NopEvent())
	}

	existingEditors := map[object.Component]T{}
	for child := range editor.Children() {
		childEdit, isEdit := child.(T)
		if !isEdit {
			continue
		}
		existingEditors[childEdit.Target()] = childEdit
	}

	for child := range object.Children(target) {
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
		log.Println("deleting unused editor for", childEdit.Target().Name())
		object.Detach(childEdit)
	}

	return editor
}
