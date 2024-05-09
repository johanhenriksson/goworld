package srv

import (
	"log"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Target interface {
	Name() string
	Attacked(by *Unit, damage uint)
}

type Unit struct {
	Emitter

	name     string
	position vec3.T
}

var _ Actor = (*Unit)(nil)

func NewUnit(name string) *Unit {
	return &Unit{
		name:     name,
		position: vec3.Zero,
	}
}

func (u *Unit) Name() string {
	return u.name
}

func (u *Unit) Position() vec3.T {
	return u.position
}

// Updates the units world position.
// Emits a position update
func (u *Unit) SetPosition(p vec3.T) {
	u.position = p
	u.Emit(UnitPositionUpdateEvent{Unit: u, Position: p})
}

type UnitPositionUpdateEvent struct {
	Unit     *Unit
	Position vec3.T
}

func (e UnitPositionUpdateEvent) Source() any {
	return e.Unit
}

type TaskEvent struct {
	Task Task
}

func (e TaskEvent) Source() any {
	return e.Task
}

type MoveToTask struct {
	Target vec3.T
	Speed  float32
}

func (t *MoveToTask) Start(u Actor) {
	// start moving towards the target
	log.Println("unit", u.Name(), "move to", t.Target)
}

func (t *MoveToTask) Step(u Actor, dt float32) bool {
	// get distance and direction to target
	dir := t.Target.Sub(u.Position())
	distSqr := dir.LengthSqr()
	dir.Normalize()
	step := t.Speed * dt

	log.Println("unit", u.Name(), "move step", dir, "remaining", distSqr, "speed", t.Speed)

	// check if we reached the target
	if distSqr < step*step {
		// reached target
		log.Println("unit", u.Name(), "arrived at", t.Target)
		u.SetPosition(t.Target)
		return true
	}

	// move towards the target
	offset := dir.Scaled(step)
	u.SetPosition(u.Position().Add(offset))
	return false
}
