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

func vec3ptr(v *vec3.T) *C.goReal {
	return (*C.goReal)(unsafe.Pointer(v))
}
