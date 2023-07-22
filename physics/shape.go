package physics

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/core/object"
)

type Shape interface {
	object.Component

	Type() ShapeType

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

type shapeBase struct {
	kind    ShapeType
	handle  shapeHandle
	changed *events.Event[Shape]
}

func newShapeBase(kind ShapeType) shapeBase {
	return shapeBase{
		kind:    kind,
		changed: events.New[Shape](),
	}
}

func (s *shapeBase) Type() ShapeType {
	return s.kind
}

func (s *shapeBase) OnChange() *events.Event[Shape] {
	return s.changed
}

func (s *shapeBase) shape() shapeHandle {
	return s.handle
}

func restoreShape(ptr unsafe.Pointer) Shape {
	if ptr == unsafe.Pointer(uintptr(0)) {
		panic("shape pointer is nil")
	}
	base := (*shapeBase)(ptr)
	switch base.kind {
	case BoxShape:
		return (*Box)(ptr)
	case CapsuleShape:
		return (*Capsule)(ptr)
	case MeshShape:
		return (*Mesh)(ptr)
	case CompoundShape:
		return (*Compound)(ptr)
	default:
		panic(fmt.Sprintf("invalid shape kind: %d", base.kind))
	}
}
