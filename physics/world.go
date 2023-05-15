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

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type World struct {
	object.T
	handle C.goDynamicsWorldHandle
}

func NewWorld() *World {
	handle := C.goCreateDynamicsWorld()
	world := object.New(&World{
		handle: handle,
	})
	runtime.SetFinalizer(world, func(w *World) {
		C.goDeleteDynamicsWorld(w.handle)
	})
	return world
}

func (w *World) Update(scene object.T, dt float32) {
	w.T.Update(scene, dt)

	w.Step(dt)
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

func (w *World) RemoveRigidBody(body *RigidBody) {
	C.goRemoveRigidBody(w.handle, body.handle)
}
