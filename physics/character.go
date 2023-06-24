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

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Character struct {
	object.T

	handle   C.goCharacterHandle
	shape    Shape
	step     float32
	world    *World
	keys     keys.State
	speed    float32
	grounded bool
}

type characterState struct {
	position vec3.T
	rotation quat.T
	grounded bool
}

func NewCharacter(height, radius, stepHeight float32) *Character {
	handle := C.goCreateCharacter(nil, C.float(height), C.float(radius), C.float(stepHeight))
	character := object.Component(&Character{
		handle: handle,
		step:   stepHeight,
		keys:   keys.NewState(),
		speed:  0.5,
	})
	runtime.SetFinalizer(character, func(c *Character) {
		C.goDeleteCharacter(c.handle)
	})
	return character
}

func (c *Character) fetchState() {
	// pull physics state
	state := characterState{}
	C.goCharacterGetState(c.handle, (*C.goCharacterState)(unsafe.Pointer(&state)))

	c.Transform().SetPosition(state.position)
	c.Transform().SetRotation(state.rotation)
	c.grounded = state.grounded
}

func (c *Character) Update(scene object.T, dt float32) {
	if c.world == nil {
		var ok bool
		c.world, ok = object.FindInParents[*World](c)
		if ok {
			c.world.AddCharacter(c)
		} else {
			return
		}
	} else {
		// detach from world if required
		world, _ := object.FindInParents[*World](c)
		if world != c.world {
			c.world.RemoveCharacter(c)
			c.world = nil
		}
	}

	// push engine state
	// todo: not required unless we changed something
	// todo: include movement dir?
	state := characterState{
		position: c.Transform().Position(),
		rotation: c.Transform().Rotation(),
	}
	C.goCharacterSetState(c.handle, (*C.goCharacterState)(unsafe.Pointer(&state)))
}

// Move the character controller. Called every frame to apply movement.
func (c *Character) Move(dir vec3.T) {
	C.goCharacterMove(c.handle, vec3ptr(&dir))
}

// Jump applies a jumping force to the character
// todo: configurable?
func (c *Character) Jump() {
	C.goCharacterJump(c.handle)
}

// Grounded returns true if the character is in contact with ground.
func (c *Character) Grounded() bool {
	return c.grounded
}
