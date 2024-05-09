package srv

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Unit struct {
	Emitter

	id       Identity
	name     string
	position vec3.T
	target   Identity
	area     Area
}

var _ Actor = (*Unit)(nil)

func NewUnit(name string) *Unit {
	return &Unit{
		id:       None,
		name:     name,
		position: vec3.Zero,
	}
}

func (u *Unit) Name() string {
	return u.name
}

func (u *Unit) String() string {
	return u.name
}

func (u *Unit) Position() vec3.T {
	return u.position
}

// Updates the units world position.
// Emits a position update
func (u *Unit) SetPosition(p vec3.T) {
	u.position = p
	u.Emit(PositionUpdateEvent{Unit: u, Position: p})
}

func (u *Unit) Target() Identity {
	return u.target
}

func (u *Unit) SetTarget(t Identity) {
	u.target = t
	u.Emit(TargetUpdateEvent{Unit: u, Target: t})
}

func (u *Unit) ID() Identity { return u.id }

func (u *Unit) Area() Area { return u.area }

func (u *Unit) Spawn(area Area, id Identity, position vec3.T) {
	// reset the unit
	u.id = id
	u.area = area
	u.position = position
	// u.state = Alive
	// u.health = u.maxHealth
	// u.power = u.maxPower

	u.Emit(AreaEnterEvent{
		Unit:     id,
		Area:     area,
		Position: position,
	})
}

func (u *Unit) Despawn() {
	u.id = None
	u.area = nil
}
