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
	"github.com/johanhenriksson/goworld/render"
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
	character := object.New(&Character{
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

func (c *Character) PreDraw(args render.Args, scene object.T) error {
	// pull physics state
	state := characterState{}
	C.goCharacterGetState(c.handle, (*C.goCharacterState)(unsafe.Pointer(&state)))

	c.Transform().SetPosition(state.position)
	c.Transform().SetRotation(state.rotation)
	c.grounded = state.grounded

	return nil
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

func (p *Character) KeyEvent(e keys.Event) {
	p.keys.KeyEvent(e)

	if !p.grounded {
		return
	}

	if e.Code() == keys.Space {
		p.Jump()
	}

	forward, right := float32(0), float32(0)
	if p.keys.Down(keys.RightArrow) {
		right += 1
	}
	if p.keys.Down(keys.LeftArrow) {
		right -= 1
	}
	if p.keys.Down(keys.UpArrow) {
		forward += 1
	}
	if p.keys.Down(keys.DownArrow) {
		forward -= 1
	}

	dirForward := p.Transform().Forward().Scaled(forward)
	dirRight := p.Transform().Right().Scaled(right)
	dir := dirForward.Add(dirRight)
	p.Move(dir.Scaled(p.speed))
}

func (c *Character) Move(dir vec3.T) {
	C.goCharacterMove(c.handle, vec3ptr(&dir))
}

func (c *Character) Jump() {
	C.goCharacterJump(c.handle)
}

func (w *World) AddCharacter(character *Character) {
	C.goAddCharacter(w.handle, character.handle)
	character.world = w
}

func (w *World) RemoveCharacter(character *Character) {
	C.goRemoveCharacter(w.handle, character.handle)
	character.world = nil
}
