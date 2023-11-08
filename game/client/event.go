package client

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Event interface {
	Apply(*Manager) error
}

//
//
//

type DisconnectEvent struct{}

func (e DisconnectEvent) Apply(m *Manager) error {
	m.Client = nil
	m.events = m.events[:0]
	return nil
}

//
// enter world
//

type EnterWorldEvent struct {
	Map string
}

func (e EnterWorldEvent) Apply(c *Manager) error {
	// unload current world
	if c.World != nil {
		c.World.Destroy()
		object.Detach(c.World)
		c.World = nil
	}

	// load map
	m := terrain.NewMap("default", 32)
	world := terrain.NewWorld(m, 200)

	object.Attach(c, world)
	c.World = world

	return nil
}

//
// entity spawn
//

type EntitySpawnEvent struct {
	EntityID server.Identity
	Position vec3.T
	Rotation float32
}

func (e EntitySpawnEvent) Apply(c *Manager) error {
	if _, exists := c.Entities[e.EntityID]; exists {
		return fmt.Errorf("entity %x already exists", e.EntityID)
	}

	// entity object
	entity := object.Builder(NewEntity(e.EntityID, e.Position, e.Rotation)).
		Position(e.Position).
		Rotation(quat.Euler(0, e.Rotation, 0)).
		Parent(c).
		Create()

	c.Entities[e.EntityID] = entity
	return nil
}

//
// entity despawn
//

type EntityDespawnEvent struct {
	EntityID server.Identity
}

func (e EntityDespawnEvent) Apply(c *Manager) error {
	if entity, exists := c.Entities[e.EntityID]; exists {
		delete(c.Entities, e.EntityID)
		object.Detach(entity)
		entity.Destroy()
	} else {
		return fmt.Errorf("cant despawn entity %x, it does not exist", e.EntityID)
	}
	return nil
}

//
// entity observe
//

type EntityObserveEvent struct {
	EntityID server.Identity
}

func (e EntityObserveEvent) Apply(c *Manager) error {
	if c.World == nil {
		return fmt.Errorf("cant observe entity %x, not in a world yet", e.EntityID)
	}

	entity, exists := c.Entities[e.EntityID]
	if !exists {
		return fmt.Errorf("cant observe entity %x, it does not exist", e.EntityID)
	}

	c.Controller.Observe(entity)

	return nil
}

type EntityMoveEvent struct {
	EntityID server.Identity
	Position vec3.T
	Rotation float32
	Stopped  bool
	Delta    float32
}

func (e EntityMoveEvent) Apply(c *Manager) error {
	entity, exists := c.Entities[e.EntityID]
	if !exists {
		return fmt.Errorf("cant move entity %x, it does not exist", e.EntityID)
	}

	// move entity
	entity.Move(e)

	return nil
}
