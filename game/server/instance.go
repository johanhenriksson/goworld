package server

import (
	"log"
	"time"

	"github.com/johanhenriksson/goworld/game/net"
)

const TickRate = 50 * time.Millisecond

type Event struct {
	Sender  *Client
	Payload net.Packet
}

type Instance struct {
	Clients []*Client

	onEvent <-chan Event
}

func (m *Instance) loop() {
	tick := time.After(TickRate)
	events := make([]Event, 0, 1024)
	for {
		select {
		case <-tick:
			// process updates
			for _, e := range events {
				m.process(e)
			}
			events = events[:0]
			tick = time.After(TickRate)

		case u := <-m.onEvent:
			events = append(events, u)
		}
	}
}

func (m *Instance) process(ev Event) {
	// process instance events
	switch data := ev.Payload.(type) {
	case *net.MovePacket:
		// move entity
		m.processEntityMove(ev.Sender, data)
	default:
		// unknown event
		log.Printf("unknown event: %v", ev)
	}

	// release packet memory
	ev.Payload.Message().Release()
}

func (m *Instance) processEntityMove(sender Entity, move *net.MovePacket) {
	moved := Identity(move.Uid())
	if sender.ID() != moved {
		// attempt to move unobserved unit
		log.Printf("client %v attempted to move unobserved unit %v", sender.ID(), moved)
		return
	}

	for _, c := range m.Clients {
		if c.ID() == moved {
			// dont return packet to sender
			continue
		}
		c.Send(move)
	}
}
