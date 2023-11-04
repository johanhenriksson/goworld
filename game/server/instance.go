package server

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

const TickRate = 50 * time.Millisecond

var instanceId = 1

type Instance struct {
	Entities map[Identity]Entity

	id      int
	onEvent chan Event
}

func NewInstance() *Instance {
	instance := &Instance{
		Entities: make(map[Identity]Entity, 1024),
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
	m.Entities[entity.ID()] = entity

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
	if _, exists := m.Entities[entity.ID()]; !exists {
		return
	}

	delete(m.Entities, entity.ID())

	// send entity despawn update to other clients
	// for _, other := range m.Entities {
	// 	if client, isClient := other.(*Client); isClient {
	// 	}
	// }
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
				t := reflect.TypeOf(e)
				log.Printf("%s: %s %+v\n", m, t.Name(), e)
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
