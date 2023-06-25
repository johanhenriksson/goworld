package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type raycastResult struct {
	shape  unsafe.Pointer
	point  vec3.T
	normal vec3.T
}

type RaycastHit struct {
	Shape  Shape
	Point  vec3.T
	Normal vec3.T
}

func (w *World) Raycast(from, to vec3.T) (hit RaycastHit, exists bool) {
	result := raycastResult{}
	hits := C.goRayCast(w.handle, vec3ptr(&from), vec3ptr(&to), (*C.goRayCastResult)(unsafe.Pointer(&result)))
	if hits > 0 {
		exists = true
		hit = RaycastHit{
			Shape:  restoreShape(result.shape),
			Point:  result.point,
			Normal: result.normal,
		}
	}
	return
}
