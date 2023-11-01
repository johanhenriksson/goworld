package server

import (
	"fmt"
	"log"
	osnet "net"

	"capnproto.org/go/capnp/v3"
	"github.com/johanhenriksson/goworld/game/net"
)

type DropInfo struct{}

type Client struct {
	Entity

	conn    osnet.Conn
	arena   *capnp.SingleSegmentArena
	encoder *capnp.Encoder
	decoder *capnp.Decoder
}

func NewClient(conn osnet.Conn) *Client {
	return &Client{
		conn:    conn,
		arena:   capnp.SingleSegment(nil),
		encoder: capnp.NewEncoder(conn),
		decoder: capnp.NewDecoder(conn),
	}
}

func (c *Client) Drop(reason string) error {
	log.Println("dropping client:", reason)
	if c.Entity != nil {
		// todo: remove entity from world etc?
	}
	if c.conn == nil {
		// wtf
		return fmt.Errorf("client already dropped")
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	c.conn = nil
	return nil
}

func (c *Client) Observe(entity Entity) error {
	if c.Entity != nil {
		// is this bad?
	}
	c.Entity = entity
	return nil
}

func (c *Client) decode() (net.Message, error) {
	msg, err := c.decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("packet decode failed: %w", err)
	}

	pkt, err := net.ReadRootPacket(msg)
	if err != nil {
		return nil, fmt.Errorf("packet read failed: %w", err)
	}
	log.Println("<-server:", msg)

	switch pkt.Which() {
	case net.Packet_Which_auth:
		return pkt.Auth()
	case net.Packet_Which_move:
		return pkt.Move()
	}

	return nil, fmt.Errorf("%w: received packet with type %v", net.ErrUnknownPacket, pkt.Which())
}

// readLoop is a goroutine that continuously reads packets from the client
// and submits them to the current instances message queue.
// todo: in the event of a read error, the client is dropped and a drop event is posted to the instance event queue.
func (c *Client) readLoop() {
	for {
		msg, err := c.decode()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("<-server:", msg)

		// todo: submit packet event to instance
		// todo: dont release here
		msg.Message().Release()
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
	log.Println("server->:", msg)

	return c.encoder.Encode(msg)
}

func (c *Client) SendMove() error {
	return c.Send(func(wrap *net.Packet) error {
		pkt, err := net.NewMovePacket(wrap.Segment())
		if err != nil {
			return err
		}

		// encode entity id
		pkt.SetUid(uint64(c.ID()))

		// encode position
		epos, err := net.FromVec3(wrap.Segment(), c.Position())
		if err != nil {
			return err
		}
		pkt.SetPosition(epos)

		// todo: encode facing

		return wrap.SetMove(pkt)
	})
}
