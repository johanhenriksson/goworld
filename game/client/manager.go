package client

import (
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
)

type Manager struct {
	object.Object
	Client     *Client
	Controller *LocalController
	Entities   map[server.Identity]Entity

	hostname    string
	lock        sync.Mutex
	events      []Event
	nextAttempt time.Time
}

func NewManager(hostname string) *Manager {
	return object.New("GameManager", &Manager{
		Controller: NewLocalController(),
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

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, event := range m.events {
		t := reflect.TypeOf(event)
		log.Printf("client event: %s %+v\n", t.Name(), event)
		event.Apply(m)
	}
	m.events = m.events[:0]

	// update world
	m.Object.Update(scene, dt)
}
