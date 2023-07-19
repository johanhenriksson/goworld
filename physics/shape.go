package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Shape interface {
	object.Component

	Type() ShapeType

	OnChange() *events.Event[Shape]

	shape() C.goShapeHandle
}

type ShapeType int

const (
	BoxShape      = ShapeType(1)
	SphereShape   = ShapeType(2)
	CylinderShape = ShapeType(3)
	CapsuleShape  = ShapeType(4)
	MeshShape     = ShapeType(5)
)

type shapeBase struct {
	kind    ShapeType
	handle  C.goShapeHandle
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

func (s *shapeBase) shape() C.goShapeHandle {
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
	default:
		panic(fmt.Sprintf("invalid shape kind: %d", base.kind))
	}
}

//
// Box shape
//

type Box struct {
	shapeBase
	object.Component

	Extents *object.Property[vec3.T]
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	box := object.NewComponent(&Box{
		shapeBase: newShapeBase(BoxShape),
		Extents:   object.NewProperty(size),
	})

	// resize shape when extents are modified
	box.Extents.OnChange().Subscribe(box, box.resize)

	// trigger initial resize
	box.resize(size)

	runtime.SetFinalizer(box, func(b *Box) {
		b.destroy()
	})
	return box
}

func (b *Box) resize(size vec3.T) {
	b.destroy()
	b.handle = C.goNewBoxShape((*C.char)(unsafe.Pointer(b)), vec3ptr(&size))
	b.OnChange().Emit(b)
}

func (b *Box) destroy() {
	if b.handle != nil {
		C.goDeleteShape(b.handle)
		b.handle = nil
	}
}

func (b *Box) Name() string {
	return "BoxCollider"
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.Extents.Get())
}

//
// Capsule shape
//

type Capsule struct {
	shapeBase
	object.Component

	Radius *object.Property[float32]
	Height *object.Property[float32]
}

var _ = &Capsule{}

func NewCapsule(height, radius float32) *Capsule {
	capsule := object.NewComponent(&Capsule{
		shapeBase: newShapeBase(CapsuleShape),
		Radius:    object.NewProperty(radius),
		Height:    object.NewProperty(height),
	})

	capsule.Radius.OnChange().Subscribe(capsule, func(radius float32) {
		capsule.resize(radius, capsule.Height.Get())
	})
	capsule.Height.OnChange().Subscribe(capsule, func(height float32) {
		capsule.resize(capsule.Radius.Get(), height)
	})

	// trigger initial resize
	capsule.resize(radius, height)

	runtime.SetFinalizer(capsule, func(c *Capsule) {
		c.destroy()
	})
	return capsule
}

func (c *Capsule) resize(radius, height float32) {
	c.destroy()
	c.handle = C.goNewCapsuleShape((*C.char)(unsafe.Pointer(c)), C.float(radius), C.float(height))
	c.OnChange().Emit(c)
}

func (c *Capsule) destroy() {
	if c.handle != nil {
		C.goDeleteShape(c.handle)
		c.handle = nil
	}
}

func (c *Capsule) String() string {
	return fmt.Sprintf("Capsule[Height=%.2f,Radius=%.2f]", c.Height.Get(), c.Radius.Get())
}
