package client

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Event interface {
	Apply(*Manager) error
}

type DisconnectEvent struct{}

func (e DisconnectEvent) Apply(m *Manager) error {
	m.Client = nil
	m.events = m.events[:0]
	return nil
}

type EntitySpawnEvent struct {
	EntityID server.Identity
	Position vec3.T
}

func (e EntitySpawnEvent) Apply(c *Manager) error {
	if _, exists := c.Entities[e.EntityID]; exists {
		return fmt.Errorf("entity %x already exists", e.EntityID)
	}

	// entity object
	entity := object.Builder(NewEntity(e.EntityID, e.Position)).
		Position(e.Position).
		Parent(c).
		Create()

	c.Entities[e.EntityID] = entity
	return nil
}

type EntityObserveEvent struct {
	EntityID server.Identity
}

func (e EntityObserveEvent) Apply(c *Manager) error {
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
}

func (e EntityMoveEvent) Apply(c *Manager) error {
	entity, exists := c.Entities[e.EntityID]
	if !exists {
		return fmt.Errorf("cant move entity %x, it does not exist", e.EntityID)
	}

	// move entity
	entity.Transform().SetPosition(e.Position)
	entity.Transform().SetRotation(quat.Euler(0, e.Rotation, 0))

	return nil
}
