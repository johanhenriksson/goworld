package client

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/player"
	"github.com/johanhenriksson/goworld/game/server"
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
	// todo
	return nil
}

type EntityObserveEvent struct {
	EntityID uint64
}

func (e EntityObserveEvent) Apply(c *Manager) error {
	// character
	char := player.New()
	char.Transform().SetPosition(vec3.New(5, 32, 5))
	object.Attach(c, char)

	return nil
}
