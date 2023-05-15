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

type Box struct {
	handle C.goShapeHandle
	size   vec3.T
}

var _ Shape = &Box{}

func NewBox(size vec3.T) *Box {
	handle := C.goNewBoxShape(vec3ptr(&size))
	shape := &Box{
		handle: handle,
		size:   size,
	}
	runtime.SetFinalizer(shape, func(b *Box) {
		C.goDeleteShape(b.shape())
	})
	return shape
}

func (b *Box) shape() C.goShapeHandle {
	return b.handle
}
