package physics

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/core/object"
)

type Shape interface {
	object.Component

	OnChange() *events.Event[Shape]

	shape() shapeHandle
}

type ShapeType int

const (
	BoxShape      = ShapeType(1)
	SphereShape   = ShapeType(2)
	CylinderShape = ShapeType(3)
	CapsuleShape  = ShapeType(4)
	MeshShape     = ShapeType(5)
	CompoundShape = ShapeType(6)
)

// validates that a shape can be restored from a physics user pointer
func checkShape(shape Shape) bool {
	ptr := reflect.ValueOf(shape).Elem().UnsafeAddr()
	restoreShape(unsafe.Pointer(ptr))
	return true
}

func restoreShape(ptr unsafe.Pointer) Shape {
	if ptr == unsafe.Pointer(uintptr(0)) {
		panic("shape pointer is nil")
	}
	kind := *(*ShapeType)(ptr)
	switch kind {
	case SphereShape:
		return (*Sphere)(ptr)
	case BoxShape:
		return (*Box)(ptr)
	case CapsuleShape:
		return (*Capsule)(ptr)
	case MeshShape:
		return (*Mesh)(ptr)
	case CompoundShape:
		return (*Compound)(ptr)
	default:
		panic(fmt.Sprintf("invalid shape kind: %d", kind))
	}
}

func Shapes(cmp object.Component) []Shape {
	shapes := []Shape{}
	group, isGroup := cmp.(object.Object)
	if !isGroup {
		if cmp.Parent() == nil {
			return nil
		} else {
			group = cmp.Parent()
		}
	}
	for _, child := range group.Children() {
		if child == cmp {
			continue
		}
		if group, isGroup := child.(object.Object); isGroup {
			if rigidbody := object.Get[*RigidBody](child); rigidbody != nil {
				// object has a rigidbody.
				// all shapes in its children belong to that rigidbody
				continue
			}
			if compound := object.Get[*Compound](child); compound != nil {
				// object has a compound shape
				// all shapes in its children belong to that compound shape
				// but the compound shape itself should be returned
				shapes = append(shapes, compound)
				continue
			}
			// the object has neither rigidbody or compound mesh
			// all its shape components belong to us
			shapes = append(shapes, object.GetAll[Shape](group)...)
		}
		if shape, isShape := child.(Shape); isShape {
			shapes = append(shapes, shape)
		}
	}
	return shapes
}

func hasParentShape(obj object.Object) bool {
	if obj == nil {
		return false
	}
	if compound := object.Get[*Compound](obj); compound != nil {
		return true
	}
	if rigidbody := object.Get[*RigidBody](obj); rigidbody != nil {
		return false
	}
	return false
}
