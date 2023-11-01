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

func (c *Client) decode() (net.Message, error) {
	msg, err := c.decoder.Decode()
	if err != nil {
		log.Println("accept failed:", err)
	}

	pkt, err := net.ReadRootPacket(msg)
	if err != nil {
		return nil, err
	}

	switch pkt.Which() {
	case net.Packet_Which_auth:
		return nil, fmt.Errorf("%w: auth packet is client->server only", net.ErrInvalidPacket)

	case net.Packet_Which_move:
		return pkt.Move()
	}

	return nil, fmt.Errorf("%w: received packet with type %v", net.ErrUnknownPacket, pkt.Which())
}

func (c *Client) readLoop() {
	for {
		msg, err := c.decode()
		if err != nil {
			log.Fatalln("client error:", err)
		}

		// todo: submit for processing somewhere
		log.Println("<-client:", msg)

		// todo: release should happen somewhere else
		msg.Message().Release()
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
	// defer msg.Release()

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
