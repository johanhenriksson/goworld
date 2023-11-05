package server

import (
	"fmt"
	"log"
	osnet "net"

	"capnproto.org/go/capnp/v3"
	"github.com/johanhenriksson/goworld/game/net"
	"github.com/johanhenriksson/goworld/math/vec3"
)

var clientId = 1

type DropInfo struct{}

type Client struct {
	Entity
	Instance *Instance

	id      int
	conn    osnet.Conn
	arena   *capnp.SingleSegmentArena
	encoder *capnp.Encoder
	decoder *capnp.Decoder
}

func NewClient(conn osnet.Conn) *Client {
	id := clientId
	clientId++
	return &Client{
		id:      id,
		conn:    conn,
		arena:   capnp.SingleSegment(nil),
		encoder: capnp.NewEncoder(conn),
		decoder: capnp.NewDecoder(conn),
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("Client[%d]", c.id)
}

func (c *Client) Drop(reason string) error {
	log.Println("dropping client:", reason)
	if c.conn == nil {
		return fmt.Errorf("client already dropped")
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	c.conn = nil
	log.Println("dropped client:", reason)

	if c.Entity != nil {
		// todo: remove entity from world etc?
	}

	return nil
}

func (c *Client) Observe(entity Entity) error {
	if c.Entity != nil && c.Entity != entity {
		return fmt.Errorf("client is already observing entity %x", c.Entity.ID())
	}

	log.Println("client", c, "observes", entity)
	c.Entity = entity

	// send observe entity packet
	return c.Send(func(wrap *net.Packet) error {
		observe, err := net.NewEntityObservePacket(wrap.Segment())
		if err != nil {
			return err
		}

		observe.SetEntity(uint64(entity.ID()))

		return wrap.SetEntityObserve(observe)
	})
}

func (c *Client) decode() (*net.Packet, error) {
	msg, err := c.decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("packet decode failed: %w", err)
	}

	pkt, err := net.ReadRootPacket(msg)
	if err != nil {
		return nil, fmt.Errorf("packet read failed: %w", err)
	}

	return &pkt, nil
}

func (c *Client) handlePacket(msg *net.Packet) error {
	defer msg.Message().Release()

	switch msg.Which() {
	case net.Packet_Which_auth:
		// auth packets should be handled by the server during the accept phase.
		// if we receive one here, consider it an error
		return fmt.Errorf("%w: unexpected auth packet", net.ErrInvalidPacket)

	case net.Packet_Which_entityMove:
		move, err := msg.EntityMove()
		if err != nil {
			return err
		}
		pos, err := move.Position()
		if err != nil {
			return err
		}
		c.Instance.SubmitEvent(EntityMoveEvent{
			Sender:   c,
			Entity:   Identity(move.Entity()),
			Position: net.ToVec3(pos),
			Rotation: move.Rotation(),
			Stopped:  move.Stopped(),
		})

	default:
		return fmt.Errorf("%w: received packet with type %v", net.ErrUnknownPacket, msg.Which())
	}

	return nil
}

// readLoop is a goroutine that continuously reads packets from the client
// and submits them to the current instances message queue.
// todo: in the event of a read error, the client is dropped and a drop event is posted to the instance event queue.
func (c *Client) readLoop() {
	defer func() {
		// disconnected
		c.Instance.SubmitEvent(DisconnectEvent{
			Entity: c,
		})
		log.Println(c, "disconnected")
	}()

	for {
		msg, err := c.decode()
		if err != nil {
			log.Println(err)
			return
		}

		if err := c.handlePacket(msg); err != nil {
			// something went to shit.
			log.Println("server error:", err)
			// todo: handle errors somehow. drop client?
			return
		}
	}
}

// Resend a message received from another client.
// Does not release the message buffer after transmission.
func (c *Client) Resend(msg net.Message) error {
	return c.encoder.Encode(msg.Message())
}

type PacketBuilderFn func(*net.Packet) error

func (c *Client) Send(fn PacketBuilderFn) error {
	// allocate message
	msg, seg, err := capnp.NewMessage(c.arena)
	if err != nil {
		return err
	}

	// remember to clean up after transmission
	defer msg.Release()

	// create packet wrapper
	wrap, err := net.NewRootPacket(seg)
	if err != nil {
		return err
	}

	// delegate to packet builder func
	if err := fn(&wrap); err != nil {
		return err
	}

	return c.encoder.Encode(msg)
}

func (c *Client) SendMove(id Identity, pos vec3.T, rot float32, stopped bool) error {
	return c.Send(func(wrap *net.Packet) error {
		move, err := net.NewEntityMovePacket(wrap.Segment())
		if err != nil {
			return err
		}

		// encode entity id
		move.SetEntity(uint64(id))

		// encode position
		epos, err := net.FromVec3(wrap.Segment(), pos)
		if err != nil {
			return err
		}
		move.SetPosition(epos)

		move.SetRotation(rot)

		move.SetStopped(stopped)

		return wrap.SetEntityMove(move)
	})
}

func (c *Client) SendSpawn(entity Entity) error {
	return c.Send(func(wrap *net.Packet) error {
		spawn, err := net.NewEntitySpawnPacket(wrap.Segment())
		if err != nil {
			return err
		}

		// entity id
		spawn.SetEntity(uint64(entity.ID()))

		// entity position
		pos, err := net.FromVec3(wrap.Segment(), entity.Position())
		if err != nil {
			return err
		}
		spawn.SetPosition(pos)

		spawn.SetRotation(entity.Rotation())

		return wrap.SetEntitySpawn(spawn)
	})
}
