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
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type RigidBody struct {
	object.T

	world  *World
	handle C.goRigidBodyHandle
	mass   float32
	shape  Shape
}

type rigidbodyState struct {
	position vec3.T
	rotation quat.T
	mass     float32
}

func NewRigidBody(mass float32, shape Shape) *RigidBody {
	body := object.New(&RigidBody{
		mass:  mass,
		shape: shape,
	})
	runtime.SetFinalizer(body, func(b *RigidBody) {
		C.goDeleteRigidBody(b.handle)
	})
	body.handle = C.goCreateRigidBody(nil, C.goReal(mass), shape.shape())
	return body
}

func (b *RigidBody) fetchState() {
	state := rigidbodyState{}
	C.goRigidBodyGetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
	b.Transform().SetPosition(state.position)
	b.Transform().SetRotation(state.rotation)
}

func (b *RigidBody) Update(scene object.T, dt float32) {
	b.T.Update(scene, dt)

	if b.world == nil {
		var ok bool
		b.world, ok = object.FindInParents[*World](b)
		if ok {
			b.world.addRigidBody(b)
		} else {
			return
		}
	} else {
		// detach from world if required
		world, _ := object.FindInParents[*World](b)
		if world != b.world {
			b.world.removeRigidBody(b)
			b.world = nil
		}
	}

	state := rigidbodyState{
		position: b.Transform().Position(),
		rotation: b.Transform().Rotation(),
	}
	C.goRigidBodySetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
}
