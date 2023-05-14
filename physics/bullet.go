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

func Hello() {
	world := NewWorld()
	world.SetGravity(vec3.New(0, -10, 0))

}

type World struct {
	handle C.goDynamicsWorldHandle
}

func NewWorld() *World {
	handle := C.goCreateDynamicsWorld()
	return &World{
		handle: handle,
	}
}

func (w *World) SetGravity(gravity vec3.T) {
	C.goSetGravity(w.handle, vec3ptr(&gravity))
}

func (w *World) Step(timestep float32) {
	C.goStepSimulation(w.handle, (C.goReal)(timestep))
}

func (w *World) AddRigidBody(body *RigidBody) {
	C.goAddRigidBody(w.handle, body.handle)
}

func vec3ptr(v *vec3.T) *C.goReal {
	return (*C.goReal)(unsafe.Pointer(v))
}

type Shape interface {
	shape() C.goCollisionShapeHandle
}

type BoxShape struct {
	handle C.goCollisionShapeHandle
	size   vec3.T
}

func NewBoxShape(size vec3.T) *BoxShape {
	handle := C.goNewBoxShape(vec3ptr(&size))
	return &BoxShape{
		handle: handle,
		size:   size,
	}
}

func (b *BoxShape) shape() C.goCollisionShapeHandle {
	return b.handle
}

type RigidBody struct {
	handle C.goRigidBodyHandle
	mass   float32
	shape  Shape
}

func NewRigidBody(mass float32, shape *BoxShape) *RigidBody {
	body := &RigidBody{
		mass:  mass,
		shape: shape,
	}
	body.handle = C.goCreateRigidBody(nil, C.goReal(mass), shape.handle)
	return body
}

func (b *RigidBody) Position() vec3.T {
	position := vec3.New(0, 0, 0)
	C.goGetPosition(b.handle, vec3ptr(&position))
	return position
}

func (b *RigidBody) SetPosition(position vec3.T) {
	C.goSetPosition(b.handle, vec3ptr(&position))
}

func (b *RigidBody) Rotation() vec3.T {
	rotation := vec3.New(0, 0, 0)
	C.goGetRotation(b.handle, vec3ptr(&rotation))
	return rotation
}
