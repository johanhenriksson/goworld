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
		b.destroy()
	})
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
		b.create()
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
func (b *RigidBody) OnAttach(obj object.T) {
	// if its a shape, we need to recreate
}

func (b *RigidBody) OnDetach(obj object.T) {
	// if its a shape, we need to recreate
}

func (b *RigidBody) create() {
	var ok bool
	if b.shape == nil {
		b.shape, ok = object.FindInChildren[Shape](b)
		if !ok {
			return
		}
		log.Println("rigidbody", b, ": found shape", b.shape)
	}

	if b.handle == nil {
		b.handle = C.goCreateRigidBody((*C.char)(unsafe.Pointer(b)), C.goReal(b.mass), b.shape.shape())
		log.Println("rigidbody", b, ": created")
	}

	if b.world == nil {
		b.world, ok = object.FindInParents[*World](b)
		if ok {
			b.world.addRigidBody(b)
			log.Println("rigidbody", b, ": added to world", b.world)
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
