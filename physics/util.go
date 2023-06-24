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

	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func init() {
	// sanity check

	vsize := unsafe.Sizeof(vec3.T{})
	if vsize != 12 {
		panic("expected vec3 to be 12 bytes")
	}

	qsize := unsafe.Sizeof(quat.T{})
	if qsize != 16 {
		panic("expected quaternion to be 16 bytes")
	}
}

func vec3ptr(v *vec3.T) *C.goVector3 {
	return (*C.goVector3)(unsafe.Pointer(v))
}

func quatPtr(q *quat.T) *C.goQuaternion {
	return (*C.goQuaternion)(unsafe.Pointer(q))
}
