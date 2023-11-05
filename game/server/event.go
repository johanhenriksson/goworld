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

	// send enter world
	if err := e.Client.SendEnterWorld("the_world"); err != nil {
		return err
	}

	// send other entities to client
	for _, other := range instance.Entities {
		if err := e.Client.SendSpawn(other); err != nil {
			return err
		}
	}

	// spawn entity
	instance.Spawn(e.Client)

	// observe entity
	return e.Client.Observe(e.Player)
}

type DisconnectEvent struct {
	Entity Entity
}

func (e DisconnectEvent) Apply(instance *Instance) error {
	log.Printf("disconnect entity %v\n", e.Entity)
	instance.Despawn(e.Entity)
	return nil
}

//
// move
//

type EntityMoveEvent struct {
	Sender   *Client
	Entity   Identity
	Position vec3.T
	Rotation float32
	Delta    float32
	Stopped  bool
}

func (e EntityMoveEvent) Apply(instance *Instance) error {
	if e.Sender.ID() != e.Entity {
		// attempt to move unobserved unit
		return fmt.Errorf("client %v attempted to move unobserved unit %v", e.Sender.ID(), e.Entity)
	}

	// update entity position
	instance.Entities[e.Entity].SetPosition(e.Position)
	instance.Entities[e.Entity].SetRotation(e.Rotation)

	// send move to other clients
	for _, other := range instance.Entities {
		if client, isClient := other.(*Client); isClient {
			if client != e.Sender {
				client.SendMove(e.Entity, e.Position, e.Rotation, e.Stopped, e.Delta)
			}
		}
	}

	return nil
}
