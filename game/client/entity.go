package client

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Entity struct {
	object.Object
	Sprite *CharacterSprite

	id server.Identity

	moveFrom   vec3.T
	moveTo     vec3.T
	moveVel    vec3.T
	rotFrom    float32
	rotTo      float32
	rotVel     float32
	lastUpdate time.Time
	duration   time.Duration
	animating  bool
	rotating   bool
	stopAfter  bool
}

func NewEntity(id server.Identity, pos vec3.T, rot float32) *Entity {
	spriteIndex := int(id % 365)
	sprite := NewCharacterSprite(spriteIndex)

	return object.New(fmt.Sprintf("Entity %x", id), &Entity{
		Sprite: sprite,

		id:         id,
		animating:  false,
		moveFrom:   pos,
		rotFrom:    rot,
		lastUpdate: time.Now(),
	})
}

func (e *Entity) EntityID() server.Identity {
	return e.id
}

func (e *Entity) Move(ev EntityMoveEvent) {
	e.rotFrom = e.Transform().Rotation().Euler().Y
	e.moveFrom = e.Transform().Position()
	e.rotTo = ev.Rotation
	e.moveTo = ev.Position
	e.stopAfter = ev.Stopped
	e.duration = time.Duration(ev.Delta * float32(time.Second))
	e.lastUpdate = time.Now()

	e.animating = true // vec3.Distance(e.moveFrom, e.moveTo) > 0.01 || math.Abs(e.rotFrom-ev.Rotation) > 0.001
	e.moveVel = e.moveTo.Sub(e.moveFrom).Scaled(1 / float32(e.duration.Seconds()))
	e.rotVel = (e.rotTo - e.rotFrom) / float32(e.duration.Seconds())
}

func (e *Entity) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)

	e.Sprite.animating = e.moveVel.LengthSqr() > 0.3

	if e.animating {
		elapsed := time.Now().Sub(e.lastUpdate)
		f := float32(elapsed.Seconds() / e.duration.Seconds())

		e.Transform().SetPosition(e.Transform().Position().Add(e.moveVel.Scaled(dt)))
		e.Transform().SetRotation(quat.Euler(0, e.Transform().Rotation().Euler().Y+e.rotVel*dt, 0))

		// if we reached the end of the move, stop
		done := f >= 1
		if e.stopAfter && done {
			e.animating = false
		}
	}
}
