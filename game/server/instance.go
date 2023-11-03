package server

import (
	"fmt"
	"log"
	"time"
)

const TickRate = 50 * time.Millisecond

var instanceId = 1

type Instance struct {
	Entities []Entity

	id      int
	onEvent chan Event
}

func NewInstance() *Instance {
	instance := &Instance{
		Entities: make([]Entity, 0, 1024),
		id:       instanceId,
		onEvent:  make(chan Event),
	}
	instanceId++
	go instance.loop()
	return instance
}

func (m *Instance) String() string {
	return fmt.Sprintf("Instance[%d]", m.id)
}

// Spawn an entity in the instance.
// (should not be called from the instance loop)
func (m *Instance) Spawn(entity Entity) {
	log.Println("spawn entity", entity)
	m.Entities = append(m.Entities, entity)

	// send entity spawn update to clients
	for _, other := range m.Entities {
		if client, isClient := other.(*Client); isClient {
			log.Println("send spawn to", client)
			if err := client.SendSpawn(entity); err != nil {
				log.Println("failed to send entity spawn")
			}
		}
	}
}

// Despawns an entity in the instance.
// (should be called from the instance loop)
func (m *Instance) Despawn(entity Entity) {
	existed := false
	for i, other := range m.Entities {
		if other == entity {
			m.Entities = append(m.Entities[:i], m.Entities[i+1:]...)
			existed = true
		}
	}

	if !existed {
		return
	}

	// send entity despawn update to other clients
	for _, other := range m.Entities {
		if client, isClient := other.(*Client); isClient {
			client.SendMove()
		}
	}
}

func (i *Instance) SubmitEvent(event Event) {
	i.onEvent <- event
}

func (m *Instance) loop() {
	tick := time.After(TickRate)
	events := make([]Event, 0, 1024)
	for {
		select {
		case <-tick:
			// reset timer
			tick = time.After(TickRate)

			// process updates
			for _, e := range events {
				log.Printf("instance event: %+v\n", e)
				if err := e.Apply(m); err != nil {
					log.Println("instance error:", err)
				}
			}
			events = events[:0]

		case u := <-m.onEvent:
			events = append(events, u)
		}
	}
}
