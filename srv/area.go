package srv

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Area is a region of the world
// It connects units, world geometry etc
type Area interface {
	Geometry
	Observable

	Action(Action)
	Join(Actor) Identity
	Leave(Identity) error

	Actor(Identity) (Actor, error)
	CastSphere(origin vec3.T, radius float32) []Actor
	CastCone(origin vec3.T, direction vec3.T, angle float32, radius float32) []Actor
}

type Geometry interface {
	// HeightAt returns the height at a given point
	HeightAt(vec3.T) float32

	// Path returns a path between two points
	Path(vec3.T, vec3.T) []vec3.T

	// Visible returns true if the two points are visible to each other
	Visible(vec3.T, vec3.T) bool
}

var nextAreaID = 1

type SimpleArea struct {
	Geometry
	Emitter

	id      int
	actions chan Action
	events  chan Event
	actors  *Pool[Actor]
	clients []Client
}

var _ Area = (*SimpleArea)(nil)

func NewSimpleArea() *SimpleArea {
	id := nextAreaID
	nextAreaID++

	a := &SimpleArea{
		id:      id,
		actions: make(chan Action, 1024),
		events:  make(chan Event, 1024),
		actors:  NewPool[Actor](1, 8),
		clients: make([]Client, 0, 8),
	}
	go a.loop()
	return a
}

func (a *SimpleArea) String() string {
	return fmt.Sprintf("Area%d", a.id)
}

func (a *SimpleArea) Action(action Action) {
	a.actions <- action
}

func (a *SimpleArea) Join(actor Actor) Identity {
	id := a.actors.Add(actor)

	actor.Subscribe(a, func(ev Event) {
		a.events <- ev
	})

	actor.Spawn(a, id, vec3.Zero)
	return id
}

func (a *SimpleArea) Leave(id Identity) error {
	actor, err := a.actors.Get(id)
	if err != nil {
		return err
	}

	actor.Despawn()
	actor.Unsubscribe(a)

	a.actors.Remove(id)
	// actor.Emit(LeaveEvent{Actor: actor, Area: a})

	return nil
}

func (a *SimpleArea) Actor(id Identity) (Actor, error) {
	return a.actors.Get(id)
}

func (a *SimpleArea) CastSphere(origin vec3.T, radius float32) []Actor {
	if radius <= 0 {
		return nil
	}

	radiusSquared := radius * radius
	return a.actors.Filter(func(actor Actor) bool {
		toActor := actor.Position().Sub(origin)
		distSquared := toActor.LengthSqr()
		return distSquared < radiusSquared
	})
}

func (a *SimpleArea) CastCone(origin vec3.T, direction vec3.T, angle float32, radius float32) []Actor {
	if radius <= 0 || angle <= 0 || angle >= math.Pi {
		return nil
	}

	radiusSquared := radius * radius
	thresh := math.Cos(angle)

	return a.actors.Filter(func(actor Actor) bool {
		toActor := actor.Position().Sub(origin)

		distSquared := toActor.LengthSqr()
		if distSquared > radiusSquared {
			return false
		}

		inside := vec3.Dot(toActor, direction) > thresh
		return inside
	})
}

func (a *SimpleArea) loop() {
	for {
		// handle area events
		select {
		case action := <-a.actions:
			if err := action.Apply(a); err != nil {
				log.Println("failed to apply action", err)
			}

		case ev := <-a.events:
			// all the events of all actors propagate up here
			// pass them on to all clients
			a.Emit(ev)
		}
	}
}

type AreaEnterEvent struct {
	Area     Area
	Unit     Identity
	Position vec3.T
}

func (e AreaEnterEvent) Source() any { return e.Unit }
