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
	object.G
	transform *Transform

	world  *World
	handle C.goRigidBodyHandle
	mass   float32

	Shape Shape
}

type rigidbodyState struct {
	position vec3.T
	rotation quat.T
	mass     float32
}

func NewRigidBody(name string, mass float32) *RigidBody {
	body := object.Group(name, &RigidBody{
		mass:      mass,
		transform: identity(),
	})
	runtime.SetFinalizer(body, func(b *RigidBody) {
		b.destroy()
	})
	return body
}

func (b *RigidBody) fetchState() {
	if b.handle == nil {
		return
	}
	state := rigidbodyState{}
	C.goRigidBodyGetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
	b.Transform().SetWorldPosition(state.position)
	b.Transform().SetWorldRotation(state.rotation)
}

func (b *RigidBody) Update(scene object.Component, dt float32) {
	b.G.Update(scene, dt)

	if b.world == nil {
		b.create()
	} else {
		// detach from world if required
		world := object.GetInParents[*World](b)
		if world != b.world {
			log.Println("remove rigidbody", b.Parent().Name(), "from world", world.Parent().Name())
			b.world.removeRigidBody(b)
			b.world = nil
		}
	}

	if b.handle != nil {
		// rigidbodies can only be attached to floating objects
		// thus, we dont need to use WorldPosition/WorldRotation
		state := rigidbodyState{
			position: b.Transform().WorldPosition(),
			rotation: b.Transform().WorldRotation(),
		}
		C.goRigidBodySetState(b.handle, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
	}
}

func (b *RigidBody) OnActivate() {
	// find shapes

	// if no shapes, exit.

	// create a rigidbody

	// find world in parent

	// if no world, exit.

	// add to world
}

func (b *RigidBody) OnDeactivate() {
	// remove from world

	// destroy ?
}

// called when a child object is attached
// or any decendant? :think:
func (b *RigidBody) OnAttach(obj object.Component) {
	// if its a shape, we need to recreate
}

func (b *RigidBody) OnDetach(obj object.Component) {
	// if its a shape, we need to recreate
}

func (b *RigidBody) create() {
	if b.Shape == nil {
		b.Shape = object.Get[Shape](b)
		if b.Shape == nil {
			log.Println("rigidbody", b.Parent().Name(), ": no shape")
			return
		}
		log.Println("rigidbody", b.Parent().Name(), ": found shape", b.Shape.Name())
	}

	if b.Shape.shape() == nil {
		// the shape is not available yet
		return
	}

	if b.handle == nil {
		b.handle = C.goCreateRigidBody((*C.char)(unsafe.Pointer(b)), C.goReal(b.mass), b.Shape.shape())
		log.Println("rigidbody", b, ": created")
	}

	if b.world == nil {
		b.world = object.GetInParents[*World](b)
		if b.world != nil {
			b.world.addRigidBody(b)
			log.Println("rigidbody", b.Parent().Name(), ": added to world", b.world.Parent().Name(), "at", b.Transform().WorldPosition())
		}
	}
}

func (b *RigidBody) destroy() {
	if b.handle != nil {
		log.Println("destroy rigidbody", b)
		C.goDeleteRigidBody(b.handle)
		b.handle = nil
	}
}

func (b *RigidBody) Transform() transform.T {
	b.transform.Recalculate(nil)
	return b.transform
}
