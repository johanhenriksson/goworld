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

type Object interface {
	fetchState()
}

type World struct {
	object.G
	handle C.goDynamicsWorldHandle
	debug  bool

	objects []Object
}

func NewWorld() *World {
	handle := C.goCreateDynamicsWorld()
	world := object.Group("World", &World{
		handle: handle,
	})
	runtime.SetFinalizer(world, func(w *World) {
		C.goDeleteDynamicsWorld(w.handle)
	})
	return world
}

func (w *World) Update(scene object.T, dt float32) {
	w.G.Update(scene, dt)

	w.step(dt)

	if w.debug {
		w.debugDraw()
	}
}

func (w *World) SetGravity(gravity vec3.T) {
	C.goSetGravity(w.handle, vec3ptr(&gravity))
}

func (w *World) step(timestep float32) {
	C.goStepSimulation(w.handle, (C.goReal)(timestep))
	// fetch physics object state
	for _, obj := range w.objects {
		obj.fetchState()
	}
}

func (w *World) addObject(obj Object) bool {
	w.objects = append(w.objects, obj)
	return true
}

func (w *World) removeObject(obj Object) bool {
	for i, o := range w.objects {
		if o == obj {
			w.objects = append(w.objects[:i], w.objects[i+1:]...)
			return true
		}
	}
	return false
}

func (w *World) addRigidBody(body *RigidBody) {
	if w.addObject(body) {
		C.goAddRigidBody(w.handle, body.handle)
	}
}

func (w *World) removeRigidBody(body *RigidBody) {
	if w.removeObject(body) {
		C.goRemoveRigidBody(w.handle, body.handle)
	}
}

func (w *World) AddCharacter(character *Character) {
	if w.addObject(character) {
		C.goAddCharacter(w.handle, character.handle)
	}
}

func (w *World) RemoveCharacter(character *Character) {
	if w.removeObject(character) {
		C.goRemoveCharacter(w.handle, character.handle)
	}
}

func (w *World) Debug(enabled bool) {
	if enabled {
		enableDebug(w)
	} else {
		disableDebug(w)
	}
	w.debug = enabled
}

func (w *World) debugDraw() {
	C.goDebugDraw(w.handle)
}
