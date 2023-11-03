package client

import (
	"fmt"
	"log"
	osnet "net"

	"capnproto.org/go/capnp/v3"
	"github.com/johanhenriksson/goworld/game/net"
)

type Client struct {
	conn    osnet.Conn
	arena   *capnp.SingleSegmentArena
	encoder *capnp.Encoder
	decoder *capnp.Decoder
}

func NewClient() (*Client, error) {
	conn, err := osnet.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", net.GamePort))
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:    conn,
		arena:   capnp.SingleSegment(nil),
		encoder: capnp.NewEncoder(conn),
		decoder: capnp.NewDecoder(conn),
	}
	go client.readLoop()
	return client, err
}

func (c *Client) decode() (*net.Packet, error) {
	msg, err := c.decoder.Decode()
	if err != nil {
		return nil, err
	}

	pkt, err := net.ReadRootPacket(msg)
	if err != nil {
		return nil, err
	}

	return &pkt, nil
}

func (c *Client) handlePacket(msg *net.Packet) error {
	switch msg.Which() {
	case net.Packet_Which_auth:
		return fmt.Errorf("%w: auth packet is client->server only", net.ErrInvalidPacket)

	case net.Packet_Which_entityMove:
		log.Println("client: entity move")

	case net.Packet_Which_entityStop:
		log.Println("client: entity stop")

	case net.Packet_Which_entitySpawn:
		spawn, err := msg.EntitySpawn()
		if err != nil {
			return err
		}
		log.Println("spawn entity", spawn.Id())

	case net.Packet_Which_entityObserve:
		observe, err := msg.EntityObserve()
		if err != nil {
			return err
		}
		log.Println("observe entity", observe.Entity())

	default:
		return fmt.Errorf("%w: received packet with type %s", net.ErrUnknownPacket, msg.Which())
	}

	return nil
}

func (c *Client) readLoop() {
	for {
		msg, err := c.decode()
		if err != nil {
			log.Println("client read error:", err)
			return
		}
		defer msg.Message().Release()

		if err := c.handlePacket(msg); err != nil {
			log.Println("client packet error:", err)
			return
		}
	}
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
	log.Println("client->:", msg)

	return c.encoder.Encode(msg)
}

func (c *Client) SendAuthToken(token uint64) error {
	return c.Send(func(p *net.Packet) error {
		auth, err := net.NewAuthPacket(p.Segment())
		if err != nil {
			return err
		}

		// encode authentication token
		auth.SetToken(token)

		return p.SetAuth(auth)
	})
}
