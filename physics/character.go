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

type Character struct {
	object.Component

	handle   C.goCharacterHandle
	shape    Shape
	step     float32
	world    *World
	grounded bool
}

type characterState struct {
	position vec3.T
	rotation quat.T
	grounded bool
}

func NewCharacter(height, radius, stepHeight float32) *Character {
	shape := NewCapsule(height, radius)
	handle := C.goCreateCharacter(shape.handle, C.float(stepHeight))
	character := object.NewComponent(&Character{
		handle: handle,
		shape:  shape,
		step:   stepHeight,
	})
	runtime.SetFinalizer(character, func(c *Character) {
		C.goDeleteCharacter(c.handle)
	})
	return character
}

func (c *Character) pullState() {
	// pull physics state
	state := characterState{}
	C.goCharacterGetState(c.handle, (*C.goCharacterState)(unsafe.Pointer(&state)))

	c.Transform().SetPosition(state.position)
	c.Transform().SetRotation(state.rotation)
	c.grounded = state.grounded
}

func (c *Character) pushState() {
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
	dir.Scale(0.016)
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

func (c *Character) OnEnable() {
	if c.world = object.GetInParents[*World](c); c.world != nil {
		c.world.AddCharacter(c)
	} else {
		log.Println("Character: No physics world in parents")
	}
}

func (c *Character) OnDisable() {
	if c.world != nil {
		c.world.RemoveCharacter(c)
		c.world = nil
	}
}
