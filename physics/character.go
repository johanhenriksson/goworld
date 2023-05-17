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

type Character struct {
	object.T

	handle C.goCharacterHandle
	shape  Shape
	step   float32
	world  *World
}

func NewCharacter(height, radius, stepHeight float32) *Character {
	handle := C.goCreateCharacter(nil, C.float(height), C.float(radius), C.float(stepHeight))
	character := object.New(&Character{
		handle: handle,
		step:   stepHeight,
	})
	runtime.SetFinalizer(character, func(c *Character) {
		C.goDeleteCharacter(c.handle)
	})
	return character
}

func (c *Character) Update(scene object.T, dt float32) {
	c.T.Update(scene, dt)

	if c.world == nil {
		var ok bool
		c.world, ok = object.FindInParents[*World](c)
		if ok {
			c.world.AddCharacter(c)
		} else {
			return
		}
	}

	C.goCharacterUpdate(c.handle, c.world.handle, C.float(dt))
}

func (c *Character) Walk(dir vec3.T) {
	C.goCharacterWalkDirection(c.handle, vec3ptr(&dir))
}

func (c *Character) Jump() {
	C.goCharacterJump(c.handle)
}

func (c *Character) Warp(pos vec3.T) {
	C.goCharacterWarp(c.handle, vec3ptr(&pos))
}

func (w *World) AddCharacter(character *Character) {
	C.goAddCharacter(w.handle, character.handle)
	character.world = w
}

func (w *World) RemoveCharacter(character *Character) {
	C.goRemoveCharacter(w.handle, character.handle)
	character.world = nil
}
