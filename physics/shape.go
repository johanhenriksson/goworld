package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"runtime"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Shape interface {
	shape() C.goShapeHandle
}

//
// Box shape
//

type Box struct {
	handle C.goShapeHandle
	size   vec3.T
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	handle := C.goNewBoxShape(vec3ptr(&size))
	box := &Box{
		handle: handle,
		size:   size,
	}
	runtime.SetFinalizer(box, func(b *Box) {
		C.goDeleteShape(b.shape())
	})
	return box
}

func (b *Box) shape() C.goShapeHandle {
	return b.handle
}

//
// Capsule shape
//

type Capsule struct {
	handle C.goShapeHandle
	height float32
	radius float32
}

var _ = &Capsule{}

func NewCapsule(radius, height float32) *Capsule {
	handle := C.goNewCapsuleShape(C.float(radius), C.float(height))
	capsule := &Capsule{
		handle: handle,
	}
	runtime.SetFinalizer(capsule, func(c *Capsule) {
		C.goDeleteShape(c.shape())
	})
	return capsule
}

func (c *Capsule) shape() C.goShapeHandle {
	return c.handle
}
