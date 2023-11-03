package server

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Event interface {
	Apply(*Instance) error
}

//
// spawn entity
//

type EntitySpawnEvent struct {
	Entity Entity
}

func (e EntitySpawnEvent) Apply(instance *Instance) error {
	log.Printf("spawn entity %v in world %v\n", e.Entity, instance)
	instance.Spawn(e.Entity)
	return nil
}

//
// enter world event
//

type EnterWorldEvent struct {
	Client *Client
	Player Entity
}

func (e *EnterWorldEvent) Apply(instance *Instance) error {
	log.Printf("client %v enter world %v with entity %v\n", e.Client, instance, e.Player)
	// assign client to instance
	e.Client.Instance = instance
	e.Client.Entity = e.Player

	// spawn entity
	instance.Spawn(e.Client)

	// observe entity
	return e.Client.Observe(e.Player)
}

//
// move
//

type EntityMoveEvent struct {
	Sender   *Client
	Entity   Identity
	Position vec3.T
	Stop     bool
}

func (e EntityMoveEvent) Apply(instance *Instance) error {
	if e.Sender.ID() != e.Entity {
		// attempt to move unobserved unit
		return fmt.Errorf("client %v attempted to move unobserved unit %v", e.Sender.ID(), e.Entity)
	}
	return nil
}
