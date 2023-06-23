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
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type RigidBody struct {
	object.T

	world  *World
	handle C.goRigidBodyHandle
	mass   float32
	shape  Shape
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

func (b *RigidBody) PreDraw(args render.Args, scene object.T) error {
	b.Transform().SetPosition(b.Position())
	b.Transform().SetRotation(b.Rotation())
	return nil
}

func (b *RigidBody) Update(scene object.T, dt float32) {
	b.T.Update(scene, dt)

	if b.world == nil {
		var ok bool
		b.world, ok = object.FindInParents[*World](b)
		if ok {
			b.world.AddRigidBody(b)
		} else {
			return
		}
	} else {
		// detach from world if required
		world, _ := object.FindInParents[*World](b)
		if world != b.world {
			b.world.RemoveRigidBody(b)
			b.world = nil
		}
	}

	b.SetPosition(b.Transform().Position())
	// b.SetRotation(b.Transform().Rotation())
}

func (b *RigidBody) Position() vec3.T {
	position := vec3.New(0, 0, 0)
	C.goGetPosition(b.handle, vec3ptr(&position))
	return position
}

func (b *RigidBody) SetPosition(position vec3.T) {
	C.goSetPosition(b.handle, vec3ptr(&position))
}

func (b *RigidBody) Rotation() quat.T {
	q := quat.Ident()
	C.goGetRotation(b.handle, quatPtr(&q))
	return q
}

func (b *RigidBody) SetRotation(q quat.T) {
	C.goSetRotation(b.handle, quatPtr(&q))
}
