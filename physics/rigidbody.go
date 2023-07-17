package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type RigidBody struct {
	object.Object
	transform *Transform

	world  *World
	handle C.goRigidBodyHandle
	mass   float32

	Shape Shape
}

var _ object.Object = &RigidBody{}

type rigidbodyState struct {
	position vec3.T
	rotation quat.T
	mass     float32
}

func NewRigidBody(name string, mass float32) *RigidBody {
	body := object.New(name, &RigidBody{
		mass:      mass,
		transform: identity(),
	})
	runtime.SetFinalizer(body, func(b *RigidBody) {
		b.destroy()
	})
	return body
}

func (b *RigidBody) pullState() {
	if b.handle == nil {
		return
	}
	if b.mass == 0 {
		return
	}
	state := rigidbodyState{}
	C.goRigidBodyGetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
	b.Transform().SetWorldPosition(state.position)
	b.Transform().SetWorldRotation(state.rotation)
}

func (b *RigidBody) pushState() {
	if b.handle == nil {
		return
	}
	// rigidbodies can only be attached to floating objects
	// thus, we dont need to use WorldPosition/WorldRotation
	state := rigidbodyState{
		position: b.Transform().WorldPosition(),
		rotation: b.Transform().WorldRotation(),
	}
	C.goRigidBodySetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
}

func (b *RigidBody) OnEnable() {
	log.Println("Enable Rigidbody", b.Name())
	if b.Shape == nil {
		b.Shape = object.Get[Shape](b)
		if b.Shape == nil {
			log.Println("rigidbody", b.Parent().Name(), ": no shape")
			return
		}
	}

	b.handle = C.goCreateRigidBody((*C.char)(unsafe.Pointer(b)), C.goReal(b.mass), b.Shape.shape())
	b.pushState()

	b.world = object.GetInParents[*World](b)
	if b.world != nil {
		b.world.addRigidBody(b)
		log.Println("Rigidbody", b.Name(), "added to physics world", b.world.Parent().Name())
	} else {
		log.Println("Rigidbody", b.Name(), ": No physics world in parents")
	}
}

func (b *RigidBody) OnDisable() {
	b.destroy()
	b.Shape = nil
	b.world = nil
}

func (b *RigidBody) destroy() {
	if b.world != nil {
		b.world.removeRigidBody(b)
		b.world = nil
	}
	if b.handle != nil {
		log.Println("destroy rigidbody", b)
		C.goDeleteRigidBody(b.handle)
		b.handle = nil
	}
}

func (b *RigidBody) Transform() transform.T {
	if b.mass == 0 {
		return b.Object.Transform()
	}
	// floating object
	b.transform.Recalculate(nil)
	return b.transform
}
