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

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Shape interface {
	object.T

	Type() ShapeType

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
	kind   ShapeType
	handle C.goShapeHandle
}

func (s *shapeBase) Type() ShapeType {
	return s.kind
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
		fmt.Println("invalid shape kind: %d", base.kind)
		return nil
	}
}

//
// Box shape
//

type Box struct {
	shapeBase
	object.T
	size vec3.T
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	box := object.New(&Box{
		shapeBase: shapeBase{
			kind: BoxShape,
		},
		size: size,
	})
	box.handle = C.goNewBoxShape((*C.char)(unsafe.Pointer(box)), vec3ptr(&size))

	runtime.SetFinalizer(box, func(b *Box) {
		C.goDeleteShape(b.shape())
	})
	return box
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.size)
}

func (b *Box) Size() vec3.T {
	return b.size
}

//
// Capsule shape
//

type Capsule struct {
	shapeBase
	object.T
	height float32
	radius float32
}

var _ = &Capsule{}

func NewCapsule(height, radius float32) *Capsule {
	capsule := object.New(&Capsule{
		shapeBase: shapeBase{
			kind: CapsuleShape,
		},
		radius: radius,
		height: height,
	})
	capsule.handle = C.goNewCapsuleShape((*C.char)(unsafe.Pointer(capsule)), C.float(radius), C.float(height))
	runtime.SetFinalizer(capsule, func(c *Capsule) {
		C.goDeleteShape(c.shape())
	})
	return capsule
}

func (c *Capsule) String() string {
	return fmt.Sprintf("Capsule[Height=%.2f,Radius=%.2f]", c.height, c.radius)
}
