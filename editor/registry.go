package editor

import (
	"log"
	"reflect"

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
func Register[K object.Component, E T](obj K, constructor func(*Context, K) E) {
	t := typeOf[K](obj)
	log.Println("register editor for", t.Name())
	editors[t] = func(ctx *Context, obj object.Component) T {
		k := obj.(K)
		return constructor(ctx, k)
	}
}

func ConstructEditors(ctx *Context, current object.Component, target object.Component) object.Component {
	var editor *EditorGhost

	if current != nil {
		editor = current.(*EditorGhost)
		if editor.target != target {
			panic("unexpected editor target")
		}
	} else {
		t := typeOf(target)
		var customEditor T
		if construct, exists := editors[t]; exists {
			log.Println("creating custom editor for", target.Name(), "of type", t.Name())
			customEditor = construct(ctx, target)
		} else {
			log.Println("creating object editor for", target.Name(), "of type", t.Name())
		}
		editor = NewEditorGhost(
			target,
			customEditor,
		)
	}

	existingEditors := map[object.Component]*EditorGhost{}
	for _, child := range editor.Children() {
		childEdit, isEdit := child.(*EditorGhost)
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
