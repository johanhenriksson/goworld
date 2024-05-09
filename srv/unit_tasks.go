package srv

import (
	"log"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type MoveToTask struct {
	Target vec3.T
	Speed  float32
}

var _ Task = (*MoveToTask)(nil)

func (t *MoveToTask) Start(u Actor) {
	// start moving towards the target
	log.Println("unit", u.Name(), "move to", t.Target)
}

func (t *MoveToTask) Stop(u Actor) {
	// stop moving
	log.Println("unit", u.Name(), "stopped moving at", u.Position())
}

func (t *MoveToTask) Step(u Actor, dt float32) bool {
	// get distance and direction to target
	dir := t.Target.Sub(u.Position())
	distSqr := dir.LengthSqr()
	dir.Normalize()
	step := t.Speed * dt

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

type AggroTask struct {
	Range float32
}

var _ Task = (*AggroTask)(nil)

func (a *AggroTask) Step(actor Actor, dt float32) bool {
	// find all units within range
	actor.Area().CastSphere(actor.Position(), a.Range)
	return false
}

func (a *AggroTask) Start(Actor) {}
func (a *AggroTask) Stop(Actor)  {}
