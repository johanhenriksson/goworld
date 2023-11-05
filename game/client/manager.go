package client

import (
	"log"
	"sync"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/game/terrain"
)

type Manager struct {
	object.Object
	Client     *Client
	Controller *LocalController
	World      *terrain.World
	Entities   map[server.Identity]Entity

	hostname    string
	lock        sync.Mutex
	events      []Event
	nextAttempt time.Time
}

func NewManager(hostname string) *Manager {
	return object.New("GameManager", &Manager{
		Controller: object.Builder(NewLocalController()).Active(false).Create(),
		Entities:   make(map[server.Identity]Entity, 64),

		hostname:    hostname,
		events:      make([]Event, 0, 1024),
		nextAttempt: time.Now(),
	})
}

func (m *Manager) queueEvent(event Event) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.events = append(m.events, event)
}

func (m *Manager) Update(scene object.Component, dt float32) {
	if m.Client == nil && m.nextAttempt.Before(time.Now()) {
		log.Println("connecting to", m.hostname)
		cli := NewClient(m.queueEvent)
		if err := cli.Connect(m.hostname); err != nil {
			log.Println("failed to connect:", err)
			m.nextAttempt = time.Now().Add(5 * time.Second)
			return
		}

		if err := cli.SendAuthToken(uint64(m.ID())); err != nil {
			log.Println("failed to authenticate:", err)
			m.nextAttempt = time.Now().Add(5 * time.Second)
			return
		}

		m.Client = cli
	}

	m.handleEvents()

	// update world (if ingame)
	if m.World != nil {
		m.Object.Update(scene, dt)
	}
}

func (m *Manager) handleEvents() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, event := range m.events {
		// t := reflect.TypeOf(event)
		// log.Printf("client event: %s %+v\n", t.Name(), event)
		if err := event.Apply(m); err != nil {
			log.Println("failed to apply event:", err)
			m.Client.Disconnect()
			m.Client = nil
			break
		}
	}
	m.events = m.events[:0]
}
